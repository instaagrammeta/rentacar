"""Rental contract service (Module 6)."""
from __future__ import annotations

from datetime import date
from typing import Any

from app.models import Rental
from app.models.enums import CarStatus, RentalStatus, ReservationStatus
from app.repositories import (
    car_repo,
    client_repo,
    rental_repo,
    reservation_repo,
    settings_repo,
)
from app.services import audit_service, car_service, client_service, pdf_service
from app.utils.errors import ValidationError
from app.utils.pricing import calculate_price


def _parse_date(value: Any) -> date:
    if isinstance(value, date):
        return value
    if not value:
        raise ValidationError("Дата обязательна")
    return date.fromisoformat(str(value)[:10])


def generate_contract_number() -> str:
    last = rental_repo.list(limit=1)
    next_id = (last[0].id + 1) if last else 1
    return f"DOG-{date.today().year}-{next_id:05d}"


def create_rental(data: dict[str, Any], actor=None) -> Rental:
    """Create a rental contract, calculate price and generate the PDF."""
    client = client_repo.get_or_404(int(data["client_id"]))
    car = car_repo.get_or_404(int(data["car_id"]))

    # Enforce age / experience / blacklist rules before signing a contract.
    client_service.assert_rental_eligibility(client)

    if car.status == CarStatus.RENTED:
        raise ValidationError("Автомобиль уже находится в аренде")
    if car.status == CarStatus.MAINTENANCE:
        raise ValidationError("Автомобиль находится на обслуживании")

    start = _parse_date(data.get("rental_start"))
    end = _parse_date(data.get("rental_end"))
    if end < start:
        raise ValidationError("Дата окончания не может быть раньше даты начала")

    breakdown = calculate_price(
        start,
        end,
        daily_price=car.daily_price,
        weekly_price=car.weekly_price,
        monthly_price=car.monthly_price,
    )
    total = float(data.get("total_price") or breakdown.total)

    rental = Rental(
        contract_number=generate_contract_number(),
        client_id=client.id,
        car_id=car.id,
        reservation_id=data.get("reservation_id"),
        employee_id=getattr(actor, "id", None) or data.get("employee_id"),
        rental_start=start,
        rental_end=end,
        deposit=float(data.get("deposit") or car.deposit_amount),
        daily_price=car.daily_price,
        total_price=total,
        start_mileage=int(data["start_mileage"]) if data.get("start_mileage") else car.mileage,
        notes=data.get("notes"),
        status=RentalStatus.ACTIVE,
    )
    rental_repo.add(rental)

    # Update car & related reservation state.
    car_service.set_status(car, CarStatus.RENTED)
    if rental.reservation_id:
        reservation = reservation_repo.get(rental.reservation_id)
        if reservation:
            reservation.status = ReservationStatus.CONFIRMED
            reservation_repo.update(reservation, {})

    # Generate the contract PDF.
    settings = settings_repo.get_settings()
    rental.pdf_path = pdf_service.generate_contract_pdf(rental, settings)
    rental_repo.update(rental, {})

    audit_service.record(
        "create_rental",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="rental",
        entity_id=rental.id,
        details=f"Договор {rental.contract_number}",
    )
    return rental


def cancel_rental(rental_id: int, actor=None) -> Rental:
    rental = rental_repo.get_or_404(rental_id)
    if rental.status == RentalStatus.COMPLETED:
        raise ValidationError("Невозможно отменить завершённую аренду")
    rental.status = RentalStatus.CANCELLED
    rental_repo.update(rental, {})
    car = car_repo.get(rental.car_id)
    if car:
        car_service.set_status(car, CarStatus.AVAILABLE)
    audit_service.record(
        "cancel_rental",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="rental",
        entity_id=rental.id,
    )
    return rental


def regenerate_pdf(rental_id: int) -> Rental:
    rental = rental_repo.get_or_404(rental_id)
    settings = settings_repo.get_settings()
    rental.pdf_path = pdf_service.generate_contract_pdf(rental, settings)
    rental_repo.update(rental, {})
    return rental
