"""Password hashing and role based access control helpers."""
from __future__ import annotations

from functools import wraps
from typing import Callable, Iterable

from flask_jwt_extended import get_jwt, jwt_required, verify_jwt_in_request
from passlib.context import CryptContext

from app.models.enums import UserRole
from app.utils.errors import PermissionError as AppPermissionError

# bcrypt based password hashing context.
_pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def hash_password(password: str) -> str:
    """Return a bcrypt hash for the given plain-text password."""
    return _pwd_context.hash(password)


def verify_password(password: str, password_hash: str) -> bool:
    """Verify a plain-text password against a stored hash."""
    try:
        return _pwd_context.verify(password, password_hash)
    except ValueError:
        return False


# Module access matrix. Maps a logical permission to the roles allowed to use it.
ROLE_PERMISSIONS: dict[str, set[UserRole]] = {
    "dashboard": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER, UserRole.CASHIER, UserRole.OPERATOR},
    "clients": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER},
    "cars": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER},
    "reservations": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER, UserRole.OPERATOR},
    "rentals": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER},
    "returns": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER},
    "payments": {UserRole.ADMINISTRATOR, UserRole.CASHIER},
    "blacklist": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER},
    "accidents": {UserRole.ADMINISTRATOR, UserRole.RENTAL_MANAGER},
    "reports": {UserRole.ADMINISTRATOR},
    "settings": {UserRole.ADMINISTRATOR},
    "users": {UserRole.ADMINISTRATOR},
    "backups": {UserRole.ADMINISTRATOR},
}


def current_role() -> UserRole:
    """Return the :class:`UserRole` from the current JWT claims."""
    claims = get_jwt()
    return UserRole(claims.get("role", UserRole.OPERATOR.value))


def role_required(*roles: UserRole) -> Callable:
    """Decorator enforcing that the caller has one of ``roles``."""

    def decorator(fn: Callable) -> Callable:
        @wraps(fn)
        @jwt_required()
        def wrapper(*args, **kwargs):
            if current_role() not in roles:
                raise AppPermissionError("Недостаточно прав для выполнения операции")
            return fn(*args, **kwargs)

        return wrapper

    return decorator


def permission_required(permission: str) -> Callable:
    """Decorator enforcing access based on the :data:`ROLE_PERMISSIONS` matrix."""
    allowed: Iterable[UserRole] = ROLE_PERMISSIONS.get(permission, set())

    def decorator(fn: Callable) -> Callable:
        @wraps(fn)
        def wrapper(*args, **kwargs):
            verify_jwt_in_request()
            if current_role() not in allowed:
                raise AppPermissionError("Недостаточно прав для выполнения операции")
            return fn(*args, **kwargs)

        return wrapper

    return decorator
