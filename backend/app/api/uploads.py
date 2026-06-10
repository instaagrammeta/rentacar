"""File upload & static media serving endpoints."""
from __future__ import annotations

from flask import Blueprint, current_app, jsonify, request, send_from_directory
from flask_jwt_extended import jwt_required

from app.utils.errors import ValidationError
from app.utils.files import save_upload

bp = Blueprint("uploads", __name__, url_prefix="/api/uploads")


@bp.post("")
@jwt_required()
def upload_file():
    """Upload one or more files; returns the stored relative paths."""
    if "files" not in request.files and "file" not in request.files:
        raise ValidationError("Файл не передан")
    subfolder = request.form.get("subfolder", "")

    stored: list[str] = []
    files = request.files.getlist("files") or [request.files["file"]]
    for file in files:
        if file and file.filename:
            stored.append(save_upload(file, subfolder=subfolder))
    return jsonify({"paths": stored}), 201


@bp.get("/<path:relative_path>")
def serve_file(relative_path: str):
    """Serve an uploaded media file (document scan, photo, QR code, logo)."""
    upload_dir = current_app.config["UPLOAD_DIR"]
    return send_from_directory(upload_dir, relative_path)
