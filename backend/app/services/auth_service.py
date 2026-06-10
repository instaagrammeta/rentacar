"""Authentication and user management service."""
from __future__ import annotations

from flask_jwt_extended import create_access_token, create_refresh_token

from app.models import User
from app.models.enums import UserRole
from app.repositories import user_repo
from app.services import audit_service
from app.utils.errors import AuthError, ConflictError, ValidationError
from app.utils.security import hash_password, verify_password


def _tokens_for(user: User) -> dict:
    """Build access / refresh tokens carrying the user's role claim."""
    identity = str(user.id)
    additional_claims = {"role": user.role.value, "username": user.username}
    access = create_access_token(identity=identity, additional_claims=additional_claims)
    refresh = create_refresh_token(identity=identity, additional_claims=additional_claims)
    return {"access_token": access, "refresh_token": refresh}


def authenticate(username: str, password: str) -> dict:
    """Validate credentials and return tokens + user payload."""
    user = user_repo.get_by_username(username)
    if user is None or not verify_password(password, user.password_hash):
        raise AuthError("Неверное имя пользователя или пароль")
    if not user.is_active:
        raise AuthError("Учётная запись отключена")

    audit_service.record("login", user_id=user.id, username=user.username, entity="user")
    return {**_tokens_for(user), "user": user.to_dict()}


def create_user(
    *,
    username: str,
    password: str,
    full_name: str = "",
    email: str | None = None,
    role: str = UserRole.OPERATOR.value,
    actor: User | None = None,
) -> User:
    """Create a new employee account."""
    if not username or not password:
        raise ValidationError("Имя пользователя и пароль обязательны")
    if len(password) < 6:
        raise ValidationError("Пароль должен содержать не менее 6 символов")
    if user_repo.get_by_username(username):
        raise ConflictError("Пользователь с таким именем уже существует")

    try:
        role_enum = UserRole(role)
    except ValueError as exc:  # pragma: no cover - guarded by schema
        raise ValidationError(f"Неизвестная роль: {role}") from exc

    user = User(
        username=username,
        full_name=full_name,
        email=email,
        role=role_enum,
        password_hash=hash_password(password),
    )
    user_repo.add(user)
    audit_service.record(
        "create_user",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="user",
        entity_id=user.id,
        details=f"Создан пользователь {username} ({role})",
    )
    return user


def change_password(user: User, new_password: str) -> None:
    if len(new_password) < 6:
        raise ValidationError("Пароль должен содержать не менее 6 символов")
    user.password_hash = hash_password(new_password)
    user_repo.update(user, {})


def list_users() -> list[User]:
    return list(user_repo.list(order_by=User.username.asc()))
