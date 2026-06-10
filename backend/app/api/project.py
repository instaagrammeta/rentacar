"""Project import / export endpoints (.rentacar custom file format)."""
from __future__ import annotations

import tempfile
from pathlib import Path

from flask import Blueprint, jsonify, request, send_file

from app.api.helpers import current_user
from app.models.enums import UserRole
from app.services import project_service
from app.utils.errors import ValidationError
from app.utils.security import role_required

bp = Blueprint("project", __name__, url_prefix="/api/project")


@bp.post("/export")
@role_required(UserRole.ADMINISTRATOR)
def export_project():
    """Файл → Сохранить проект — export the whole database to a .rentacar file."""
    path = project_service.export_project(actor=current_user())
    file = Path(path)
    return send_file(
        file,
        as_attachment=True,
        download_name=file.name,
        mimetype="application/octet-stream",
    )


@bp.post("/import")
@role_required(UserRole.ADMINISTRATOR)
def import_project():
    """Файл → Открыть проект — import a .rentacar file and restore all data."""
    if "file" not in request.files:
        raise ValidationError("Файл проекта не передан")
    upload = request.files["file"]
    if not upload.filename or not upload.filename.endswith(".rentacar"):
        raise ValidationError("Ожидается файл с расширением .rentacar")

    # Persist to a temporary file, then import.
    with tempfile.NamedTemporaryFile(delete=False, suffix=".rentacar") as tmp:
        upload.save(tmp.name)
        temp_path = tmp.name
    try:
        result = project_service.import_project(temp_path, actor=current_user())
    finally:
        Path(temp_path).unlink(missing_ok=True)
    return jsonify({"message": "Проект импортирован", **result})
