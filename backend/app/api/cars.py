"""Car management endpoints (Module 4)."""
from __future__ import annotations

from flask import Blueprint, jsonify, request

from app.api.helpers import current_user, get_pagination, json_body
from app.models import Car
from app.models.enums import CarStatus
from app.repositories import car_repo
from app.services import car_service
from app.utils.security import permission_required

bp = Blueprint("cars", __name__, url_prefix="/api/cars")


@bp.get("")
@permission_required("cars")
def list_cars():
    page, per_page, offset = get_pagination()
    status = request.args.get("status", "").strip()
    search = request.args.get("search", "").strip()

    filters = []
    if status:
        filters.append(Car.status == CarStatus(status))
    if search:
        like = f"%{search}%"
        filters.append(
            Car.brand.ilike(like) | Car.model.ilike(like) | Car.plate_number.ilike(like)
        )

    total = car_repo.count(filters or None)
    items = car_repo.list(filters=filters or None, limit=per_page, offset=offset)
    return jsonify(
        {"items": [c.to_dict() for c in items], "total": total, "page": page, "per_page": per_page}
    )


@bp.get("/available")
@permission_required("cars")
def available_cars():
    items = car_repo.list(filters=[Car.status == CarStatus.AVAILABLE], limit=500)
    return jsonify([c.to_dict() for c in items])


@bp.get("/<int:car_id>")
@permission_required("cars")
def get_car(car_id: int):
    return jsonify(car_repo.get_or_404(car_id).to_dict())


@bp.post("")
@permission_required("cars")
def create_car():
    car = car_service.create_car(json_body(), actor=current_user())
    return jsonify(car.to_dict()), 201


@bp.put("/<int:car_id>")
@permission_required("cars")
def update_car(car_id: int):
    car = car_service.update_car(car_id, json_body(), actor=current_user())
    return jsonify(car.to_dict())


@bp.delete("/<int:car_id>")
@permission_required("cars")
def delete_car(car_id: int):
    car_service.delete_car(car_id, actor=current_user())
    return jsonify({"message": "Автомобиль удалён"})
