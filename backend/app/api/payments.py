"""Payment endpoints (Module 8)."""
from __future__ import annotations

from flask import Blueprint, current_app, jsonify, request, send_file

from app.api.helpers import current_user, get_pagination, json_body
from app.models import Payment
from app.repositories import payment_repo
from app.services import payment_service
from app.utils.security import permission_required

bp = Blueprint("payments", __name__, url_prefix="/api/payments")


@bp.get("")
@permission_required("payments")
def list_payments():
    page, per_page, offset = get_pagination()
    client_id = request.args.get("client_id", type=int)
    rental_id = request.args.get("rental_id", type=int)
    filters = []
    if client_id:
        filters.append(Payment.client_id == client_id)
    if rental_id:
        filters.append(Payment.rental_id == rental_id)
    total = payment_repo.count(filters or None)
    items = payment_repo.list(filters=filters or None, limit=per_page, offset=offset)
    return jsonify(
        {"items": [p.to_dict() for p in items], "total": total, "page": page, "per_page": per_page}
    )


@bp.get("/<int:payment_id>")
@permission_required("payments")
def get_payment(payment_id: int):
    return jsonify(payment_repo.get_or_404(payment_id).to_dict())


@bp.post("")
@permission_required("payments")
def create_payment():
    payment = payment_service.create_payment(json_body(), actor=current_user())
    return jsonify(payment.to_dict()), 201


@bp.get("/<int:payment_id>/receipt")
@permission_required("payments")
def download_receipt(payment_id: int):
    payment = payment_repo.get_or_404(payment_id)
    if not payment.receipt_path:
        return jsonify({"error": "Квитанция не найдена"}), 404
    full_path = current_app.config["EXPORT_DIR"] / payment.receipt_path
    return send_file(
        full_path, as_attachment=True, download_name=f"{payment.receipt_number}.pdf"
    )
