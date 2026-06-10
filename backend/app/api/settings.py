"""Company settings endpoints (Module 12)."""
from __future__ import annotations

from flask import Blueprint, jsonify

from app.api.helpers import current_user, json_body
from app.models.enums import UserRole
from app.services import settings_service
from app.utils.security import role_required
from flask_jwt_extended import jwt_required

bp = Blueprint("settings", __name__, url_prefix="/api/settings")


@bp.get("")
@jwt_required()
def get_settings():
    # All authenticated users may read company info (used for headers / PDFs).
    return jsonify(settings_service.get_settings().to_dict())


@bp.put("")
@role_required(UserRole.ADMINISTRATOR)
def update_settings():
    settings = settings_service.update_settings(json_body(), actor=current_user())
    return jsonify(settings.to_dict())
