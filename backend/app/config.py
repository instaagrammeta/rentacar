"""Application configuration.

Centralised configuration objects for the Rentacar CRM backend. Values can be
overridden through environment variables which makes the same code base usable
in development, testing and the packaged desktop build.
"""
from __future__ import annotations

import os
from datetime import timedelta
from pathlib import Path

# Base directory of the backend package (…/backend)
BASE_DIR = Path(__file__).resolve().parent.parent

# Directory that stores runtime data: SQLite database, uploads, backups, exports.
# When running as a packaged desktop application this can be redirected to the
# user's AppData directory through the ``RENTACAR_DATA_DIR`` environment variable.
DATA_DIR = Path(os.environ.get("RENTACAR_DATA_DIR", BASE_DIR / "data"))
UPLOAD_DIR = DATA_DIR / "uploads"
BACKUP_DIR = Path(os.environ.get("RENTACAR_BACKUP_DIR", BASE_DIR.parent / "backups"))
EXPORT_DIR = DATA_DIR / "exports"

for _directory in (DATA_DIR, UPLOAD_DIR, BACKUP_DIR, EXPORT_DIR):
    _directory.mkdir(parents=True, exist_ok=True)


class BaseConfig:
    """Shared configuration values."""

    SECRET_KEY: str = os.environ.get("RENTACAR_SECRET_KEY", "change-me-in-production")
    JWT_SECRET_KEY: str = os.environ.get("RENTACAR_JWT_SECRET", SECRET_KEY)
    JWT_ACCESS_TOKEN_EXPIRES: timedelta = timedelta(
        hours=int(os.environ.get("RENTACAR_JWT_HOURS", "12"))
    )
    JWT_REFRESH_TOKEN_EXPIRES: timedelta = timedelta(days=30)

    SQLALCHEMY_DATABASE_URI: str = os.environ.get(
        "RENTACAR_DATABASE_URI", f"sqlite:///{(DATA_DIR / 'rentacar.db').as_posix()}"
    )
    SQLALCHEMY_TRACK_MODIFICATIONS: bool = False
    SQLALCHEMY_ENGINE_OPTIONS: dict = {"pool_pre_ping": True}

    # File handling
    MAX_CONTENT_LENGTH: int = 32 * 1024 * 1024  # 32 MB upload limit
    UPLOAD_DIR: Path = UPLOAD_DIR
    BACKUP_DIR: Path = BACKUP_DIR
    EXPORT_DIR: Path = EXPORT_DIR
    DATA_DIR: Path = DATA_DIR

    # Business rules
    MIN_CLIENT_AGE: int = 21
    MIN_DRIVING_EXPERIENCE_YEARS: int = 1

    # Automatic backups
    ENABLE_SCHEDULED_BACKUPS: bool = (
        os.environ.get("RENTACAR_ENABLE_BACKUPS", "true").lower() == "true"
    )
    BACKUP_RETENTION_DAYS: int = int(os.environ.get("RENTACAR_BACKUP_RETENTION", "30"))


class DevelopmentConfig(BaseConfig):
    DEBUG = True


class TestingConfig(BaseConfig):
    TESTING = True
    SQLALCHEMY_DATABASE_URI = "sqlite:///:memory:"
    ENABLE_SCHEDULED_BACKUPS = False


class ProductionConfig(BaseConfig):
    DEBUG = False


_CONFIG_MAP = {
    "development": DevelopmentConfig,
    "testing": TestingConfig,
    "production": ProductionConfig,
}


def get_config(name: str | None = None) -> type[BaseConfig]:
    """Return the configuration class for the given environment name."""
    name = name or os.environ.get("RENTACAR_ENV", "development")
    return _CONFIG_MAP.get(name, DevelopmentConfig)
