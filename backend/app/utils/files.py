"""Helpers for storing uploaded files (document scans, photos, logos)."""
from __future__ import annotations

import uuid
from pathlib import Path

from werkzeug.datastructures import FileStorage
from werkzeug.utils import secure_filename

from flask import current_app

ALLOWED_IMAGE_EXTENSIONS = {".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".pdf"}


def _upload_dir() -> Path:
    path: Path = current_app.config["UPLOAD_DIR"]
    path.mkdir(parents=True, exist_ok=True)
    return path


def save_upload(file: FileStorage, subfolder: str = "") -> str:
    """Persist an uploaded file and return its path relative to UPLOAD_DIR.

    A random prefix is added to avoid collisions while keeping the original
    extension for previewing.
    """
    original = secure_filename(file.filename or "file")
    ext = Path(original).suffix.lower()
    if ext and ext not in ALLOWED_IMAGE_EXTENSIONS:
        from app.utils.errors import ValidationError

        raise ValidationError(f"Недопустимый тип файла: {ext}")

    target_dir = _upload_dir()
    if subfolder:
        target_dir = target_dir / secure_filename(subfolder)
        target_dir.mkdir(parents=True, exist_ok=True)

    filename = f"{uuid.uuid4().hex}{ext}"
    relative = filename if not subfolder else f"{subfolder}/{filename}"
    file.save(target_dir / filename)
    return relative


def absolute_path(relative: str | None) -> Path | None:
    """Resolve a stored relative path to an absolute filesystem path."""
    if not relative:
        return None
    return _upload_dir() / relative
