"""Client management service (Modules 2 & 3).

Implements validation rules (minimum age & driving experience), client codes,
QR codes for VIP / returning customers and aggregated client history.
"""
from __future__ import annotations

from datetime import date
from typing import Any

from flask import current_app

from app.models import Client
from app.models.enums import ClientStatus
from app.repositories import (
    accident_repo,
    blacklist_repo,
    client_repo,
    payment_repo,
    rental_repo,
)
from app.services import audit_service
from app.utils.errors import BusinessRuleError, ConflictError, NotFoundError, ValidationError
from app.utils.qr import generate_qr


def _calculate_age(born: date | None, on: date | None = None) -> int | None:
    if born is None:
        return None
    on = on or date.today()
    return on.year - born.year - ((on.month, on.day) < (born.month, born.day))


def generate_client_code() -> str:
    """Generate the next sequential client code, e.g. ``CL-000042``."""
    last = client_repo.list(limit=1)
    next_id = (last[0].id + 1) if last else 1
    return f"CL-{next_id:06d}"


def assert_rental_eligibility(client: Client) -> None:
    """Raise :class:`BusinessRuleError` if the client may not rent a car.

    Enforces the minimum age, minimum driving experience and blacklist rules.
    """
    if client.status == ClientStatus.BLACKLISTED or client.blacklist_entry is not None:
        raise BusinessRuleError("Клиент находится в чёрном списке и не может арендовать автомобиль")

    min_age = current_app.config["MIN_CLIENT_AGE"]
    min_exp = current_app.config["MIN_DRIVING_EXPERIENCE_YEARS"]

    age = _calculate_age(client.date_of_birth)
    if age is None or age < min_age:
        raise BusinessRuleError(
            f"Минимальный возраст клиента — {min_age} лет (текущий: {age if age is not None else 'не указан'})"
        )
    if (client.driver_experience_years or 0) < min_exp:
        raise BusinessRuleError(
            f"Минимальный стаж вождения — {min_exp} год(а) (текущий: {client.driver_experience_years or 0})"
        )


def create_client(data: dict[str, Any], actor=None) -> Client:
    """Create a client, generating a unique code and QR identifier."""
    if not data.get("first_name") or not data.get("last_name"):
        raise ValidationError("Имя и фамилия обязательны")
    if not data.get("phone"):
        raise ValidationError("Телефон обязателен")

    if client_repo.get_by_phone(data["phone"]):
        raise ConflictError("Клиент с таким телефоном уже существует")

    client = Client(
        client_code=generate_client_code(),
        first_name=data["first_name"].strip(),
        last_name=data["last_name"].strip(),
        phone=data["phone"].strip(),
        email=(data.get("email") or None),
        date_of_birth=_parse_date(data.get("date_of_birth")),
        passport_number=data.get("passport_number"),
        driver_license_number=data.get("driver_license_number"),
        driver_license_issue_date=_parse_date(data.get("driver_license_issue_date")),
        driver_experience_years=int(data.get("driver_experience_years") or 0),
        passport_scan=data.get("passport_scan"),
        driver_license_scan=data.get("driver_license_scan"),
        is_vip=bool(data.get("is_vip", False)),
        notes=data.get("notes"),
        status=ClientStatus(data.get("status", ClientStatus.PENDING_VERIFICATION.value)),
    )
    client_repo.add(client)

    # Generate the QR code now that we have an id and code.
    client.qr_code_path = generate_qr(client.client_code, f"{client.client_code}.png")
    client_repo.update(client, {})

    audit_service.record(
        "create_client",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="client",
        entity_id=client.id,
        details=f"Создан клиент {client.full_name}",
    )
    return client


def update_client(client_id: int, data: dict[str, Any], actor=None) -> Client:
    client = client_repo.get_or_404(client_id)
    for field in (
        "first_name",
        "last_name",
        "phone",
        "email",
        "passport_number",
        "driver_license_number",
        "passport_scan",
        "driver_license_scan",
        "notes",
        "is_vip",
    ):
        if field in data and data[field] is not None:
            setattr(client, field, data[field])

    if "driver_experience_years" in data and data["driver_experience_years"] is not None:
        client.driver_experience_years = int(data["driver_experience_years"])
    for date_field in ("date_of_birth", "driver_license_issue_date"):
        if date_field in data:
            setattr(client, date_field, _parse_date(data[date_field]))
    if "status" in data and data["status"]:
        client.status = ClientStatus(data["status"])

    client_repo.update(client, {})
    audit_service.record(
        "update_client",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="client",
        entity_id=client.id,
    )
    return client


def delete_client(client_id: int, actor=None) -> None:
    client = client_repo.get_or_404(client_id)
    if client.rentals:
        raise ConflictError("Невозможно удалить клиента с историей аренды")
    client_repo.delete(client)
    audit_service.record(
        "delete_client",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="client",
        entity_id=client_id,
    )


def find_client(*, phone: str | None = None, code: str | None = None, term: str | None = None):
    """Look up a returning client by phone, QR/client code, or free text."""
    if code:
        client = client_repo.get_by_code(code)
        if client:
            return client
    if phone:
        client = client_repo.get_by_phone(phone)
        if client:
            return client
    if term:
        return client_repo.search(term)
    raise NotFoundError("Клиент не найден")


def get_history(client_id: int) -> dict:
    """Aggregate the full history of a client (rentals, payments, etc.)."""
    client = client_repo.get_or_404(client_id)
    rentals = [r.to_dict() for r in client.rentals]
    payments = [p.to_dict() for p in client.payments]
    accidents = [a.to_dict() for a in client.accidents]
    penalties = [
        {
            "rental_id": r.id,
            "contract_number": r.contract_number,
            "amount": r.vehicle_return.penalties + r.vehicle_return.late_fee,
        }
        for r in client.rentals
        if r.vehicle_return
        and (r.vehicle_return.penalties or r.vehicle_return.late_fee)
    ]
    return {
        "client": client.to_dict(),
        "rentals": rentals,
        "payments": payments,
        "penalties": penalties,
        "accidents": accidents,
    }


def _parse_date(value: Any) -> date | None:
    if not value:
        return None
    if isinstance(value, date):
        return value
    return date.fromisoformat(str(value)[:10])
