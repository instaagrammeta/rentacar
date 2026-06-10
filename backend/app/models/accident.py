"""Accident record model (Module 10) and its photos."""
from __future__ import annotations

from datetime import date

from sqlalchemy import Date, Float, ForeignKey, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel


class Accident(BaseModel):
    """An accident or damage event involving a car and (optionally) a client."""

    __tablename__ = "accidents"

    car_id: Mapped[int] = mapped_column(ForeignKey("cars.id"), nullable=False)
    client_id: Mapped[int | None] = mapped_column(ForeignKey("clients.id"))
    rental_id: Mapped[int | None] = mapped_column(ForeignKey("rentals.id"))

    accident_date: Mapped[date] = mapped_column(Date, nullable=False)
    description: Mapped[str | None] = mapped_column(Text)
    repair_cost: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)

    car: Mapped["Car"] = relationship(back_populates="accidents")  # noqa: F821
    client: Mapped["Client | None"] = relationship(back_populates="accidents")  # noqa: F821
    photos: Mapped[list["AccidentPhoto"]] = relationship(
        back_populates="accident", cascade="all, delete-orphan"
    )

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        data = super().to_dict(exclude=exclude)
        data["photos"] = [p.file_path for p in self.photos]
        if self.car:
            data["car_name"] = self.car.display_name
        if self.client:
            data["client_name"] = self.client.full_name
        return data


class AccidentPhoto(BaseModel):
    """A photo attached to an accident record."""

    __tablename__ = "accident_photos"

    accident_id: Mapped[int] = mapped_column(
        ForeignKey("accidents.id", ondelete="CASCADE"), nullable=False
    )
    file_path: Mapped[str] = mapped_column(String(255), nullable=False)

    accident: Mapped["Accident"] = relationship(back_populates="photos")
