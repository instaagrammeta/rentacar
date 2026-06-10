"""Automatic SQLite backup service.

Creates timestamped copies of the SQLite database inside the backups folder and
prunes copies older than the configured retention period. Designed to be run on
a daily schedule by APScheduler (see :mod:`app.scheduler`).
"""
from __future__ import annotations

import logging
import shutil
from datetime import datetime, timedelta
from pathlib import Path

from flask import current_app

logger = logging.getLogger("rentacar.backup")


def _database_file() -> Path | None:
    uri: str = current_app.config["SQLALCHEMY_DATABASE_URI"]
    if not uri.startswith("sqlite:///"):
        return None
    return Path(uri.replace("sqlite:///", "", 1))


def create_backup() -> str | None:
    """Create a timestamped backup of the SQLite database."""
    db_file = _database_file()
    if db_file is None or not db_file.exists():
        logger.warning("Резервное копирование пропущено: база данных не найдена")
        return None

    backup_dir: Path = current_app.config["BACKUP_DIR"]
    backup_dir.mkdir(parents=True, exist_ok=True)
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    target = backup_dir / f"rentacar_backup_{timestamp}.db"
    shutil.copy2(db_file, target)
    logger.info("Создана резервная копия: %s", target)

    prune_old_backups()
    return str(target)


def prune_old_backups() -> int:
    """Remove backups older than the retention window. Returns count removed."""
    backup_dir: Path = current_app.config["BACKUP_DIR"]
    retention_days: int = current_app.config["BACKUP_RETENTION_DAYS"]
    cutoff = datetime.now() - timedelta(days=retention_days)
    removed = 0
    for file in backup_dir.glob("rentacar_backup_*.db"):
        if datetime.fromtimestamp(file.stat().st_mtime) < cutoff:
            file.unlink(missing_ok=True)
            removed += 1
    if removed:
        logger.info("Удалено старых резервных копий: %s", removed)
    return removed


def list_backups() -> list[dict]:
    backup_dir: Path = current_app.config["BACKUP_DIR"]
    items = []
    for file in sorted(backup_dir.glob("rentacar_backup_*.db"), reverse=True):
        stat = file.stat()
        items.append(
            {
                "name": file.name,
                "size": stat.st_size,
                "created_at": datetime.fromtimestamp(stat.st_mtime).isoformat(),
            }
        )
    return items
