"""Audit log endpoints (security activity trail)."""
from __future__ import annotations

from flask import Blueprint, jsonify, request

from app.models.enums import UserRole
from app.services import audit_service
from app.utils.security import role_required

bp = Blueprint("audit", __name__, url_prefix="/api/audit")


@bp.get("")
@role_required(UserRole.ADMINISTRATOR)
def list_logs():
    limit = request.args.get("limit", default=200, type=int)
    offset = request.args.get("offset", default=0, type=int)
    logs = audit_service.list_logs(limit=limit, offset=offset)
    return jsonify([log.to_dict() for log in logs])
