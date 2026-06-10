"""Vehicle return service (Module 7).

Computes extra days, late fees, damage costs and the final payment owed (or to
be refunded) when a car is returned.
"""
from __future__ import annotations

from datetime import date, datetime
from typing import Any

from app.models import VehicleReturn
from app.models.enums import CarStatus, RentalStatus
from app.repositories import car_repo, rental_repo
from app.services import audit_service, car_service
from app.utils.errors import ConflictError, ValidationError


def _parse_datetime(value: Any) -> datetime:
    if isinstance(value, datetime):
        return value
    if not value:
        return datetime.now()
    text = str(value)
    try:
        return datetime.fromisoformat(text)
    except ValueError:
        return datetime.combine(date.fromisoformat(text[:10]), datetime.min.time())


def calculate_return(
    rental,
    *,
    return_date: datetime,
    damage_cost: float = 0.0,
    penalties: float = 0.0,
    late_fee_multiplier: float = 1.5,
) -> dict:
    """Calculate the financial outcome of a return.

    * Extra days are days returned after the agreed end date.
    * Late fee = extra_days * daily_price * multiplier.
    * Final payment = late fee + damage cost + penalties - deposit.
      A negative value means the company must refund the client.
    """
    agreed_end = rental.rental_end
    actual_end = return_date.date() if isinstance(return_date, datetime) else return_date

    extra_days = max((actual_end - agreed_end).days, 0)
    late_fee = round(extra_days * rental.daily_price * late_fee_multiplier, 2)
    final_payment = round(late_fee + damage_cost + penalties - rental.deposit, 2)

    return {
        "extra_days": extra_days,
        "late_fee": late_fee,
        "damage_cost": round(damage_cost, 2),
        "penalties": round(penalties, 2),
        "final_payment": final_payment,
    }


def create_return(rental_id: int, data: dict[str, Any], actor=None) -> VehicleReturn:
    rental = rental_repo.get_or_404(rental_id)
    if rental.vehicle_return is not None:
        raise ConflictError("Возврат для этой аренды уже оформлен")
    if rental.status == RentalStatus.CANCELLED:
        raise ValidationError("Невозможно оформить возврат для отменённой аренды")

    return_date = _parse_datetime(data.get("return_date"))
    damage_cost = float(data.get("damage_cost") or 0)
    penalties = float(data.get("penalties") or 0)

    calc = calculate_return(
        rental,
        return_date=return_date,
        damage_cost=damage_cost,
        penalties=penalties,
    )

    vehicle_return = VehicleReturn(
        rental_id=rental.id,
        return_date=return_date,
        mileage=int(data.get("mileage") or rental.car.mileage),
        fuel_level=int(data.get("fuel_level") or 100),
        damages=data.get("damages"),
        extra_days=calc["extra_days"],
        late_fee=calc["late_fee"],
        damage_cost=calc["damage_cost"],
        penalties=calc["penalties"],
        final_payment=calc["final_payment"],
    )
    rental.vehicle_return = vehicle_return
    rental.status = RentalStatus.COMPLETED

    # Update odometer and free the car.
    car = car_repo.get(rental.car_id)
    if car:
        car.mileage = vehicle_return.mileage
        car_service.set_status(car, CarStatus.AVAILABLE)

    rental_repo.update(rental, {})

    audit_service.record(
        "create_return",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="rental",
        entity_id=rental.id,
        details=f"Возврат: итог {calc['final_payment']:.2f}",
    )
    return vehicle_return
