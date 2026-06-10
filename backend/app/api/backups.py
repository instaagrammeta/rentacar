"""Backup management endpoints."""
from __future__ import annotations

from flask import Blueprint, jsonify

from app.api.helpers import current_user
from app.services import backup_service
from app.utils.security import permission_required

bp = Blueprint("backups", __name__, url_prefix="/api/backups")


@bp.get("")
@permission_required("backups")
def list_backups():
    return jsonify(backup_service.list_backups())


@bp.post("")
@permission_required("backups")
def create_backup():
    path = backup_service.create_backup()
    if path is None:
        return jsonify({"error": "База данных не найдена"}), 400
    return jsonify({"message": "Резервная копия создана", "path": path}), 201
