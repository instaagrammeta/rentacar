"""Blacklist endpoints (Module 9)."""
from __future__ import annotations

from flask import Blueprint, jsonify

from app.api.helpers import current_user, json_body
from app.models.enums import BlacklistReason
from app.services import blacklist_service
from app.utils.security import permission_required

bp = Blueprint("blacklist", __name__, url_prefix="/api/blacklist")


@bp.get("")
@permission_required("blacklist")
def list_blacklist():
    return jsonify([e.to_dict() for e in blacklist_service.list_blacklist()])


@bp.get("/reasons")
@permission_required("blacklist")
def reasons():
    return jsonify([{"value": r.value, "label": r.label} for r in BlacklistReason])


@bp.post("")
@permission_required("blacklist")
def add_to_blacklist():
    entry = blacklist_service.add_to_blacklist(json_body(), actor=current_user())
    return jsonify(entry.to_dict()), 201


@bp.delete("/<int:client_id>")
@permission_required("blacklist")
def remove_from_blacklist(client_id: int):
    blacklist_service.remove_from_blacklist(client_id, actor=current_user())
    return jsonify({"message": "Клиент удалён из чёрного списка"})
