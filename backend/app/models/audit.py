"""Audit log model used for the security / activity trail."""
from __future__ import annotations

from sqlalchemy import ForeignKey, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel


class AuditLog(BaseModel):
    """Records logins, actions and data changes performed by users."""

    __tablename__ = "audit_logs"

    user_id: Mapped[int | None] = mapped_column(ForeignKey("users.id"))
    username: Mapped[str | None] = mapped_column(String(64))
    action: Mapped[str] = mapped_column(String(64), nullable=False, index=True)
    entity: Mapped[str | None] = mapped_column(String(64), index=True)
    entity_id: Mapped[int | None] = mapped_column()
    ip_address: Mapped[str | None] = mapped_column(String(64))
    details: Mapped[str | None] = mapped_column(Text)

    user: Mapped["User | None"] = relationship(back_populates="audit_logs")  # noqa: F821
