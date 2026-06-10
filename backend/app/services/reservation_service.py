"""Reservation service (Module 5)."""
from __future__ import annotations

from datetime import date
from typing import Any

from app.models import Reservation
from app.models.enums import CarStatus, ReservationStatus
from app.repositories import car_repo, client_repo, reservation_repo
from app.services import audit_service, car_service, client_service
from app.utils.errors import ValidationError


def _parse_date(value: Any) -> date:
    if isinstance(value, date):
        return value
    if not value:
        raise ValidationError("Дата обязательна")
    return date.fromisoformat(str(value)[:10])


def create_reservation(data: dict[str, Any], actor=None) -> Reservation:
    client = client_repo.get_or_404(int(data["client_id"]))
    car = car_repo.get_or_404(int(data["car_id"]))

    # Blacklist / eligibility check up front.
    client_service.assert_rental_eligibility(client)

    start = _parse_date(data.get("start_date"))
    end = _parse_date(data.get("end_date"))
    if end < start:
        raise ValidationError("Дата окончания не может быть раньше даты начала")
    if car.status in (CarStatus.RENTED, CarStatus.MAINTENANCE):
        raise ValidationError("Автомобиль недоступен для бронирования")

    reservation = Reservation(
        client_id=client.id,
        car_id=car.id,
        start_date=start,
        end_date=end,
        deposit=float(data.get("deposit") or car.deposit_amount),
        notes=data.get("notes"),
        status=ReservationStatus.RESERVED,
    )
    reservation_repo.add(reservation)
    car_service.set_status(car, CarStatus.RESERVED)

    audit_service.record(
        "create_reservation",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="reservation",
        entity_id=reservation.id,
    )
    return reservation


def update_status(reservation_id: int, status: str, actor=None) -> Reservation:
    reservation = reservation_repo.get_or_404(reservation_id)
    new_status = ReservationStatus(status)
    reservation.status = new_status
    reservation_repo.update(reservation, {})

    car = car_repo.get(reservation.car_id)
    if car:
        if new_status == ReservationStatus.CANCELLED and car.status == CarStatus.RESERVED:
            car_service.set_status(car, CarStatus.AVAILABLE)
        elif new_status in (ReservationStatus.RESERVED, ReservationStatus.CONFIRMED):
            car_service.set_status(car, CarStatus.RESERVED)

    audit_service.record(
        "reservation_status",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="reservation",
        entity_id=reservation.id,
        details=f"Статус: {new_status.label}",
    )
    return reservation


def cancel(reservation_id: int, actor=None) -> Reservation:
    return update_status(reservation_id, ReservationStatus.CANCELLED.value, actor=actor)
