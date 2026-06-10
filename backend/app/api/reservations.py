"""Reservation endpoints (Module 5)."""
from __future__ import annotations

from flask import Blueprint, jsonify, request

from app.api.helpers import current_user, get_pagination, json_body
from app.models import Reservation
from app.models.enums import ReservationStatus
from app.repositories import reservation_repo
from app.services import reservation_service
from app.utils.security import permission_required

bp = Blueprint("reservations", __name__, url_prefix="/api/reservations")


@bp.get("")
@permission_required("reservations")
def list_reservations():
    page, per_page, offset = get_pagination()
    status = request.args.get("status", "").strip()
    filters = []
    if status:
        filters.append(Reservation.status == ReservationStatus(status))
    total = reservation_repo.count(filters or None)
    items = reservation_repo.list(filters=filters or None, limit=per_page, offset=offset)
    return jsonify(
        {"items": [r.to_dict() for r in items], "total": total, "page": page, "per_page": per_page}
    )


@bp.get("/<int:reservation_id>")
@permission_required("reservations")
def get_reservation(reservation_id: int):
    return jsonify(reservation_repo.get_or_404(reservation_id).to_dict())


@bp.post("")
@permission_required("reservations")
def create_reservation():
    reservation = reservation_service.create_reservation(json_body(), actor=current_user())
    return jsonify(reservation.to_dict()), 201


@bp.post("/<int:reservation_id>/status")
@permission_required("reservations")
def update_status(reservation_id: int):
    body = json_body()
    reservation = reservation_service.update_status(
        reservation_id, body.get("status", ""), actor=current_user()
    )
    return jsonify(reservation.to_dict())


@bp.post("/<int:reservation_id>/cancel")
@permission_required("reservations")
def cancel_reservation(reservation_id: int):
    reservation = reservation_service.cancel(reservation_id, actor=current_user())
    return jsonify(reservation.to_dict())
