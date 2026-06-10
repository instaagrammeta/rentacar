"""Authentication & user-management endpoints."""
from __future__ import annotations

from flask import Blueprint, jsonify
from flask_jwt_extended import jwt_required

from app.api.helpers import current_user, json_body
from app.models.enums import UserRole
from app.services import auth_service
from app.utils.errors import ValidationError
from app.utils.security import role_required

bp = Blueprint("auth", __name__, url_prefix="/api/auth")


@bp.post("/login")
def login():
    body = json_body()
    username = body.get("username", "").strip()
    password = body.get("password", "")
    if not username or not password:
        raise ValidationError("Введите имя пользователя и пароль")
    result = auth_service.authenticate(username, password)
    return jsonify(result)


@bp.get("/me")
@jwt_required()
def me():
    user = current_user()
    if user is None:
        return jsonify({"error": "Не авторизован"}), 401
    return jsonify(user.to_dict())


@bp.get("/roles")
@jwt_required()
def roles():
    return jsonify([{"value": r.value, "label": r.label} for r in UserRole])


@bp.get("/users")
@role_required(UserRole.ADMINISTRATOR)
def list_users():
    return jsonify([u.to_dict() for u in auth_service.list_users()])


@bp.post("/users")
@role_required(UserRole.ADMINISTRATOR)
def create_user():
    body = json_body()
    user = auth_service.create_user(
        username=body.get("username", ""),
        password=body.get("password", ""),
        full_name=body.get("full_name", ""),
        email=body.get("email"),
        role=body.get("role", UserRole.OPERATOR.value),
        actor=current_user(),
    )
    return jsonify(user.to_dict()), 201


@bp.post("/users/<int:user_id>/password")
@role_required(UserRole.ADMINISTRATOR)
def change_password(user_id: int):
    from app.repositories import user_repo

    user = user_repo.get_or_404(user_id)
    body = json_body()
    auth_service.change_password(user, body.get("password", ""))
    return jsonify({"message": "Пароль обновлён"})
