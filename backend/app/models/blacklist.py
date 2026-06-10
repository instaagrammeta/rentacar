"""Blacklist entry model (Module 9)."""
from __future__ import annotations

from sqlalchemy import Enum, ForeignKey, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel
from app.models.enums import BlacklistReason


class BlacklistEntry(BaseModel):
    """Marks a client as blacklisted, preventing new rentals."""

    __tablename__ = "blacklist_entries"

    client_id: Mapped[int] = mapped_column(
        ForeignKey("clients.id"), unique=True, nullable=False
    )
    reason: Mapped[BlacklistReason] = mapped_column(
        Enum(BlacklistReason, native_enum=False, length=32), nullable=False
    )
    comment: Mapped[str | None] = mapped_column(Text)
    created_by: Mapped[int | None] = mapped_column(ForeignKey("users.id"))

    client: Mapped["Client"] = relationship(back_populates="blacklist_entry")  # noqa: F821

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        data = super().to_dict(exclude=exclude)
        data["reason"] = (
            self.reason.value if isinstance(self.reason, BlacklistReason) else self.reason
        )
        data["reason_label"] = BlacklistReason(self.reason).label
        if self.client:
            data["client_name"] = self.client.full_name
        return data
