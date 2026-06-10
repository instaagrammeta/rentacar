"""API package: registers every REST blueprint on the Flask app."""
from __future__ import annotations

from flask import Flask


def register_blueprints(app: Flask) -> None:
    """Import and register all API blueprints."""
    from app.api import (
        accidents,
        audit,
        auth,
        backups,
        blacklist,
        cars,
        clients,
        dashboard,
        payments,
        project,
        rentals,
        reports,
        reservations,
        settings,
        uploads,
    )

    for module in (
        auth,
        clients,
        cars,
        reservations,
        rentals,
        payments,
        blacklist,
        accidents,
        dashboard,
        reports,
        settings,
        backups,
        project,
        uploads,
        audit,
    ):
        app.register_blueprint(module.bp)
