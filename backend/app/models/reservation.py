"""Reservation model (Module 5)."""
from __future__ import annotations

from datetime import date

from sqlalchemy import Date, Enum, Float, ForeignKey, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel
from app.models.enums import ReservationStatus


class Reservation(BaseModel):
    """A booking that reserves a car for a client for a date range."""

    __tablename__ = "reservations"

    client_id: Mapped[int] = mapped_column(ForeignKey("clients.id"), nullable=False)
    car_id: Mapped[int] = mapped_column(ForeignKey("cars.id"), nullable=False)

    start_date: Mapped[date] = mapped_column(Date, nullable=False)
    end_date: Mapped[date] = mapped_column(Date, nullable=False)
    deposit: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    notes: Mapped[str | None] = mapped_column(Text)

    status: Mapped[ReservationStatus] = mapped_column(
        Enum(ReservationStatus, native_enum=False, length=32),
        default=ReservationStatus.RESERVED,
        nullable=False,
    )

    client: Mapped["Client"] = relationship(back_populates="reservations")  # noqa: F821
    car: Mapped["Car"] = relationship(back_populates="reservations")  # noqa: F821

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        data = super().to_dict(exclude=exclude)
        data["status_label"] = ReservationStatus(self.status).label
        data["status"] = (
            self.status.value if isinstance(self.status, ReservationStatus) else self.status
        )
        if self.client:
            data["client_name"] = self.client.full_name
        if self.car:
            data["car_name"] = self.car.display_name
        return data
