"""Car management service (Module 4)."""
from __future__ import annotations

from typing import Any

from app.models import Car, CarPhoto
from app.models.enums import CarStatus
from app.repositories import car_repo
from app.services import audit_service
from app.utils.errors import ConflictError, ValidationError

_NUMERIC_FIELDS = (
    "year",
    "mileage",
    "daily_price",
    "weekly_price",
    "monthly_price",
    "deposit_amount",
)


def create_car(data: dict[str, Any], actor=None) -> Car:
    if not data.get("brand") or not data.get("model"):
        raise ValidationError("Марка и модель обязательны")
    if not data.get("plate_number"):
        raise ValidationError("Гос. номер обязателен")

    existing = car_repo.list(filters=[Car.plate_number == data["plate_number"]])
    if existing:
        raise ConflictError("Автомобиль с таким гос. номером уже существует")

    car = Car(
        brand=data["brand"].strip(),
        model=data["model"].strip(),
        year=int(data.get("year") or 0),
        color=data.get("color"),
        vin=data.get("vin"),
        plate_number=data["plate_number"].strip(),
        mileage=int(data.get("mileage") or 0),
        daily_price=float(data.get("daily_price") or 0),
        weekly_price=float(data.get("weekly_price") or 0),
        monthly_price=float(data.get("monthly_price") or 0),
        deposit_amount=float(data.get("deposit_amount") or 0),
        status=CarStatus(data.get("status", CarStatus.AVAILABLE.value)),
    )
    car_repo.add(car)

    for path in data.get("photos", []) or []:
        car.photos.append(CarPhoto(file_path=path))
    if car.photos:
        car_repo.update(car, {})

    audit_service.record(
        "create_car",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="car",
        entity_id=car.id,
        details=f"Добавлен автомобиль {car.display_name}",
    )
    return car


def update_car(car_id: int, data: dict[str, Any], actor=None) -> Car:
    car = car_repo.get_or_404(car_id)
    for field in ("brand", "model", "color", "vin", "plate_number"):
        if field in data and data[field] is not None:
            setattr(car, field, data[field])
    for field in _NUMERIC_FIELDS:
        if field in data and data[field] is not None:
            setattr(car, field, type(getattr(car, field))(data[field]))
    if "status" in data and data["status"]:
        car.status = CarStatus(data["status"])
    if "photos" in data and data["photos"] is not None:
        car.photos.clear()
        for path in data["photos"]:
            car.photos.append(CarPhoto(file_path=path))

    car_repo.update(car, {})
    audit_service.record(
        "update_car",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="car",
        entity_id=car.id,
    )
    return car


def set_status(car: Car, status: CarStatus) -> None:
    """Update a car's availability status."""
    car.status = status
    car_repo.update(car, {})


def delete_car(car_id: int, actor=None) -> None:
    car = car_repo.get_or_404(car_id)
    if car.rentals:
        raise ConflictError("Невозможно удалить автомобиль с историей аренды")
    car_repo.delete(car)
    audit_service.record(
        "delete_car",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="car",
        entity_id=car_id,
    )
