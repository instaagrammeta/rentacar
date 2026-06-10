"""Shared helpers for API blueprints."""
from __future__ import annotations

from flask import request
from flask_jwt_extended import get_jwt_identity

from app.models import User
from app.repositories import user_repo


def current_user() -> User | None:
    """Return the :class:`User` referenced by the current JWT identity."""
    identity = get_jwt_identity()
    if identity is None:
        return None
    try:
        return user_repo.get(int(identity))
    except (TypeError, ValueError):
        return None


def get_pagination() -> tuple[int, int, int]:
    """Parse ``page`` / ``per_page`` query params. Returns (page, per_page, offset)."""
    try:
        page = max(int(request.args.get("page", 1)), 1)
    except ValueError:
        page = 1
    try:
        per_page = min(max(int(request.args.get("per_page", 20)), 1), 200)
    except ValueError:
        per_page = 20
    return page, per_page, (page - 1) * per_page


def json_body() -> dict:
    """Return the request JSON body as a dict (empty dict if none)."""
    return request.get_json(silent=True) or {}
