"""Accident endpoints (Module 10)."""
from __future__ import annotations

from flask import Blueprint, jsonify, request

from app.api.helpers import current_user, get_pagination, json_body
from app.models import Accident
from app.repositories import accident_repo
from app.services import accident_service
from app.utils.security import permission_required

bp = Blueprint("accidents", __name__, url_prefix="/api/accidents")


@bp.get("")
@permission_required("accidents")
def list_accidents():
    page, per_page, offset = get_pagination()
    car_id = request.args.get("car_id", type=int)
    filters = []
    if car_id:
        filters.append(Accident.car_id == car_id)
    total = accident_repo.count(filters or None)
    items = accident_repo.list(filters=filters or None, limit=per_page, offset=offset)
    return jsonify(
        {"items": [a.to_dict() for a in items], "total": total, "page": page, "per_page": per_page}
    )


@bp.get("/<int:accident_id>")
@permission_required("accidents")
def get_accident(accident_id: int):
    return jsonify(accident_repo.get_or_404(accident_id).to_dict())


@bp.post("")
@permission_required("accidents")
def create_accident():
    accident = accident_service.create_accident(json_body(), actor=current_user())
    return jsonify(accident.to_dict()), 201


@bp.put("/<int:accident_id>")
@permission_required("accidents")
def update_accident(accident_id: int):
    accident = accident_service.update_accident(accident_id, json_body(), actor=current_user())
    return jsonify(accident.to_dict())


@bp.delete("/<int:accident_id>")
@permission_required("accidents")
def delete_accident(accident_id: int):
    accident_service.delete_accident(accident_id, actor=current_user())
    return jsonify({"message": "Запись о ДТП удалена"})
