"""User / employee model used for authentication and authorisation."""
from __future__ import annotations

from sqlalchemy import Boolean, Enum, String
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel
from app.models.enums import UserRole


class User(BaseModel):
    """An employee account that can log into the CRM."""

    __tablename__ = "users"

    username: Mapped[str] = mapped_column(String(64), unique=True, nullable=False, index=True)
    email: Mapped[str | None] = mapped_column(String(120), unique=True)
    full_name: Mapped[str] = mapped_column(String(160), nullable=False, default="")
    password_hash: Mapped[str] = mapped_column(String(255), nullable=False)
    role: Mapped[UserRole] = mapped_column(
        Enum(UserRole, native_enum=False, length=32), nullable=False, default=UserRole.OPERATOR
    )
    is_active: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)

    # Relationships
    audit_logs: Mapped[list["AuditLog"]] = relationship(  # noqa: F821
        back_populates="user", cascade="all, delete-orphan"
    )

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        exclude = (exclude or set()) | {"password_hash"}
        data = super().to_dict(exclude=exclude)
        data["role"] = self.role.value if isinstance(self.role, UserRole) else self.role
        data["role_label"] = UserRole(self.role).label
        return data
