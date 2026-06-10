"""Application specific exceptions and a Flask error handler registrar.

Errors raised by the service layer are translated into consistent JSON
responses with Russian messages so the frontend can display them directly.
"""
from __future__ import annotations

from flask import Flask, jsonify


class AppError(Exception):
    """Base class for all expected application errors."""

    status_code = 400

    def __init__(self, message: str, status_code: int | None = None, payload: dict | None = None):
        super().__init__(message)
        self.message = message
        if status_code is not None:
            self.status_code = status_code
        self.payload = payload or {}

    def to_dict(self) -> dict:
        data = {"error": self.message, **self.payload}
        return data


class NotFoundError(AppError):
    status_code = 404


class ValidationError(AppError):
    status_code = 422


class AuthError(AppError):
    status_code = 401


class PermissionError(AppError):  # noqa: A001 - intentional domain name
    status_code = 403


class ConflictError(AppError):
    status_code = 409


class BusinessRuleError(AppError):
    """Raised when a domain rule (e.g. age limit) is violated."""

    status_code = 422


def register_error_handlers(app: Flask) -> None:
    """Attach JSON error handlers to the Flask application."""

    @app.errorhandler(AppError)
    def handle_app_error(error: AppError):  # type: ignore[no-redef]
        response = jsonify(error.to_dict())
        response.status_code = error.status_code
        return response

    @app.errorhandler(404)
    def handle_404(_error):  # type: ignore[no-redef]
        return jsonify({"error": "Ресурс не найден"}), 404

    @app.errorhandler(500)
    def handle_500(error):  # type: ignore[no-redef]
        app.logger.exception("Внутренняя ошибка сервера: %s", error)
        return jsonify({"error": "Внутренняя ошибка сервера"}), 500
