"""Accident service (Module 10)."""
from __future__ import annotations

from datetime import date
from typing import Any

from app.models import Accident, AccidentPhoto
from app.repositories import accident_repo, car_repo
from app.services import audit_service
from app.utils.errors import ValidationError


def _parse_date(value: Any) -> date:
    if isinstance(value, date):
        return value
    if not value:
        return date.today()
    return date.fromisoformat(str(value)[:10])


def create_accident(data: dict[str, Any], actor=None) -> Accident:
    car = car_repo.get_or_404(int(data["car_id"]))

    accident = Accident(
        car_id=car.id,
        client_id=int(data["client_id"]) if data.get("client_id") else None,
        rental_id=int(data["rental_id"]) if data.get("rental_id") else None,
        accident_date=_parse_date(data.get("accident_date")),
        description=data.get("description"),
        repair_cost=float(data.get("repair_cost") or 0),
    )
    for path in data.get("photos", []) or []:
        accident.photos.append(AccidentPhoto(file_path=path))

    accident_repo.add(accident)
    audit_service.record(
        "create_accident",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="accident",
        entity_id=accident.id,
    )
    return accident


def update_accident(accident_id: int, data: dict[str, Any], actor=None) -> Accident:
    accident = accident_repo.get_or_404(accident_id)
    if "description" in data:
        accident.description = data["description"]
    if "repair_cost" in data and data["repair_cost"] is not None:
        accident.repair_cost = float(data["repair_cost"])
    if "accident_date" in data and data["accident_date"]:
        accident.accident_date = _parse_date(data["accident_date"])
    if "photos" in data and data["photos"] is not None:
        accident.photos.clear()
        for path in data["photos"]:
            accident.photos.append(AccidentPhoto(file_path=path))
    accident_repo.update(accident, {})
    return accident


def delete_accident(accident_id: int, actor=None) -> None:
    accident = accident_repo.get_or_404(accident_id)
    accident_repo.delete(accident)
    audit_service.record(
        "delete_accident",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="accident",
        entity_id=accident_id,
    )
