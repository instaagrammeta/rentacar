"""Custom ``.rentacar`` project file format (Export / Import Project).

A ``.rentacar`` file is a ZIP archive containing:

* ``manifest.json`` — format version and metadata.
* ``data.json``     — a full dump of every domain table (clients, cars,
  reservations, rentals, returns, payments, accidents, blacklist, settings,
  users).
* ``uploads/``      — referenced document scans, photos, QR codes and logo.

Importing restores all data into the current database. The import is
transactional at the table level: existing rows are cleared and replaced so the
opened project fully represents the file's contents.
"""
from __future__ import annotations

import json
import logging
import zipfile
from datetime import date, datetime
from pathlib import Path

from flask import current_app

from app.extensions import db
from app.models import (
    Accident,
    AccidentPhoto,
    BlacklistEntry,
    Car,
    CarPhoto,
    Client,
    CompanySettings,
    Payment,
    Rental,
    Reservation,
    User,
    VehicleReturn,
)
from app.services import audit_service
from app.utils.errors import ValidationError

logger = logging.getLogger("rentacar.project")

FORMAT_VERSION = "1.0"
FILE_EXTENSION = ".rentacar"

# Order matters for import because of foreign keys.
_EXPORT_ORDER: list[tuple[str, type]] = [
    ("users", User),
    ("company_settings", CompanySettings),
    ("clients", Client),
    ("cars", Car),
    ("car_photos", CarPhoto),
    ("reservations", Reservation),
    ("rentals", Rental),
    ("vehicle_returns", VehicleReturn),
    ("payments", Payment),
    ("accidents", Accident),
    ("accident_photos", AccidentPhoto),
    ("blacklist_entries", BlacklistEntry),
]


def _serialise_row(obj) -> dict:
    """Serialise a model row into JSON-friendly primitives."""
    row: dict = {}
    for column in obj.__table__.columns:
        value = getattr(obj, column.name)
        if isinstance(value, (datetime, date)):
            value = value.isoformat()
        elif hasattr(value, "value"):  # Enum
            value = value.value
        row[column.name] = value
    return row


def _collect_upload_paths(data: dict) -> set[str]:
    """Gather all relative upload paths referenced by the exported data."""
    paths: set[str] = set()
    for table, rows in data.items():
        for row in rows:
            for key in ("passport_scan", "driver_license_scan", "qr_code_path", "file_path", "logo_path"):
                value = row.get(key)
                if value:
                    paths.add(value)
    return paths


def export_project(actor=None) -> str:
    """Export the entire database into a ``.rentacar`` file. Returns its path."""
    data: dict[str, list[dict]] = {}
    for table_name, model in _EXPORT_ORDER:
        rows = db.session.query(model).all()
        data[table_name] = [_serialise_row(r) for r in rows]

    manifest = {
        "format": "rentacar",
        "version": FORMAT_VERSION,
        "exported_at": datetime.now().isoformat(),
        "tables": {name: len(rows) for name, rows in data.items()},
    }

    export_dir: Path = current_app.config["EXPORT_DIR"] / "projects"
    export_dir.mkdir(parents=True, exist_ok=True)
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    target = export_dir / f"project_{timestamp}{FILE_EXTENSION}"

    upload_root: Path = current_app.config["UPLOAD_DIR"]
    with zipfile.ZipFile(target, "w", zipfile.ZIP_DEFLATED) as archive:
        archive.writestr("manifest.json", json.dumps(manifest, ensure_ascii=False, indent=2))
        archive.writestr("data.json", json.dumps(data, ensure_ascii=False, indent=2))
        for rel_path in _collect_upload_paths(data):
            source = upload_root / rel_path
            if source.exists():
                archive.write(source, f"uploads/{rel_path}")

    audit_service.record(
        "export_project",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="project",
        details=str(target),
    )
    logger.info("Проект экспортирован: %s", target)
    return str(target)


def _coerce_value(model, column_name: str, value):
    """Convert serialised primitives back into Python types for the column."""
    if value is None:
        return None
    column = model.__table__.columns.get(column_name)
    if column is None:
        return value
    python_type = getattr(column.type, "python_type", None)
    try:
        if python_type is datetime and isinstance(value, str):
            return datetime.fromisoformat(value)
        if python_type is date and isinstance(value, str):
            return date.fromisoformat(value[:10])
    except (ValueError, NotImplementedError):
        return value
    return value


def import_project(file_path: str, actor=None) -> dict:
    """Import a ``.rentacar`` archive, replacing all current data."""
    path = Path(file_path)
    if not path.exists():
        raise ValidationError("Файл проекта не найден")
    if not zipfile.is_zipfile(path):
        raise ValidationError("Некорректный формат файла .rentacar")

    with zipfile.ZipFile(path, "r") as archive:
        names = archive.namelist()
        if "data.json" not in names:
            raise ValidationError("В файле отсутствует data.json")
        data = json.loads(archive.read("data.json").decode("utf-8"))

        # Restore upload files.
        upload_root: Path = current_app.config["UPLOAD_DIR"]
        for name in names:
            if name.startswith("uploads/") and not name.endswith("/"):
                relative = name[len("uploads/") :]
                target = upload_root / relative
                target.parent.mkdir(parents=True, exist_ok=True)
                target.write_bytes(archive.read(name))

    # Clear existing data in reverse FK order, then insert.
    for _, model in reversed(_EXPORT_ORDER):
        db.session.query(model).delete()
    db.session.commit()

    counts: dict[str, int] = {}
    for table_name, model in _EXPORT_ORDER:
        rows = data.get(table_name, [])
        for row in rows:
            cleaned = {k: _coerce_value(model, k, v) for k, v in row.items()}
            db.session.add(model(**cleaned))
        counts[table_name] = len(rows)
    db.session.commit()

    audit_service.record(
        "import_project",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="project",
        details=f"Импортировано: {counts}",
    )
    logger.info("Проект импортирован: %s", counts)
    return {"imported": counts}
