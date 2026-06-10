"""Rental contract & vehicle return endpoints (Modules 6 & 7)."""
from __future__ import annotations

from flask import Blueprint, jsonify, request, send_file

from app.api.helpers import current_user, get_pagination, json_body
from app.models import Rental
from app.models.enums import RentalStatus
from app.repositories import rental_repo
from app.services import rental_service, return_service
from app.utils.security import permission_required

from flask import current_app

bp = Blueprint("rentals", __name__, url_prefix="/api/rentals")


@bp.get("")
@permission_required("rentals")
def list_rentals():
    page, per_page, offset = get_pagination()
    status = request.args.get("status", "").strip()
    filters = []
    if status:
        filters.append(Rental.status == RentalStatus(status))
    total = rental_repo.count(filters or None)
    items = rental_repo.list(filters=filters or None, limit=per_page, offset=offset)
    return jsonify(
        {"items": [r.to_dict() for r in items], "total": total, "page": page, "per_page": per_page}
    )


@bp.get("/<int:rental_id>")
@permission_required("rentals")
def get_rental(rental_id: int):
    rental = rental_repo.get_or_404(rental_id)
    data = rental.to_dict()
    if rental.vehicle_return:
        data["vehicle_return"] = rental.vehicle_return.to_dict()
    return jsonify(data)


@bp.post("")
@permission_required("rentals")
def create_rental():
    rental = rental_service.create_rental(json_body(), actor=current_user())
    return jsonify(rental.to_dict()), 201


@bp.post("/<int:rental_id>/cancel")
@permission_required("rentals")
def cancel_rental(rental_id: int):
    rental = rental_service.cancel_rental(rental_id, actor=current_user())
    return jsonify(rental.to_dict())


@bp.get("/<int:rental_id>/contract")
@permission_required("rentals")
def download_contract(rental_id: int):
    rental = rental_repo.get_or_404(rental_id)
    if not rental.pdf_path:
        rental = rental_service.regenerate_pdf(rental_id)
    full_path = current_app.config["EXPORT_DIR"] / rental.pdf_path
    return send_file(full_path, as_attachment=True, download_name=f"{rental.contract_number}.pdf")


# ----------------------------------------------------------- vehicle returns
@bp.post("/<int:rental_id>/return/preview")
@permission_required("returns")
def preview_return(rental_id: int):
    """Calculate return charges without persisting (live preview in the UI)."""
    from datetime import datetime

    rental = rental_repo.get_or_404(rental_id)
    body = json_body()
    return_date = body.get("return_date")
    parsed = (
        datetime.fromisoformat(return_date) if return_date else datetime.now()
    )
    calc = return_service.calculate_return(
        rental,
        return_date=parsed,
        damage_cost=float(body.get("damage_cost") or 0),
        penalties=float(body.get("penalties") or 0),
    )
    return jsonify(calc)


@bp.post("/<int:rental_id>/return")
@permission_required("returns")
def create_return(rental_id: int):
    vehicle_return = return_service.create_return(rental_id, json_body(), actor=current_user())
    return jsonify(vehicle_return.to_dict()), 201
