"""Application factory for the Rentacar CRM backend.

Creates and configures the Flask application, binds extensions, registers API
blueprints, configures logging, JWT error handlers and the backup scheduler.
"""
from __future__ import annotations

import logging
from logging.handlers import RotatingFileHandler

from flask import Flask, jsonify

from app.api import register_blueprints
from app.config import BaseConfig, get_config
from app.extensions import cors, db, jwt, migrate
from app.utils.errors import register_error_handlers


def _configure_logging(app: Flask) -> None:
    """Configure rotating file + console logging."""
    log_dir = app.config["DATA_DIR"] / "logs"
    log_dir.mkdir(parents=True, exist_ok=True)
    handler = RotatingFileHandler(
        log_dir / "rentacar.log", maxBytes=2 * 1024 * 1024, backupCount=5, encoding="utf-8"
    )
    formatter = logging.Formatter("%(asctime)s [%(levelname)s] %(name)s: %(message)s")
    handler.setFormatter(formatter)
    level = logging.DEBUG if app.config.get("DEBUG") else logging.INFO

    root = logging.getLogger("rentacar")
    root.setLevel(level)
    if not root.handlers:
        root.addHandler(handler)
        console = logging.StreamHandler()
        console.setFormatter(formatter)
        root.addHandler(console)


def _register_jwt_handlers() -> None:
    """Return consistent JSON errors for authentication failures."""

    @jwt.unauthorized_loader
    def _missing_token(reason: str):  # type: ignore[no-redef]
        return jsonify({"error": "Требуется авторизация"}), 401

    @jwt.invalid_token_loader
    def _invalid_token(reason: str):  # type: ignore[no-redef]
        return jsonify({"error": "Недействительный токен"}), 401

    @jwt.expired_token_loader
    def _expired_token(header, payload):  # type: ignore[no-redef]
        return jsonify({"error": "Срок действия токена истёк"}), 401


def create_app(config_name: str | None = None) -> Flask:
    """Application factory."""
    app = Flask(__name__)
    config: type[BaseConfig] = get_config(config_name)
    app.config.from_object(config)

    # Bind extensions.
    db.init_app(app)
    migrate.init_app(app, db)
    jwt.init_app(app)
    cors.init_app(app, resources={r"/api/*": {"origins": "*"}})

    _configure_logging(app)
    _register_jwt_handlers()
    register_error_handlers(app)
    register_blueprints(app)

    # Import models so they are registered with SQLAlchemy metadata.
    from app import models  # noqa: F401

    @app.get("/api/health")
    def health():
        return jsonify({"status": "ok", "service": "Rentacar CRM"})

    # Create tables automatically when not using migrations (desktop build).
    with app.app_context():
        db.create_all()

    # Start the daily backup scheduler.
    from app.scheduler import init_scheduler

    init_scheduler(app)

    app.logger.info("Приложение Rentacar CRM инициализировано")
    return app
