"""Rental contract model (Module 6) and the vehicle return (Module 7)."""
from __future__ import annotations

from datetime import date, datetime

from sqlalchemy import Date, DateTime, Enum, Float, ForeignKey, Integer, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel
from app.models.enums import RentalStatus


class Rental(BaseModel):
    """A signed rental contract between a client and the company."""

    __tablename__ = "rentals"

    contract_number: Mapped[str] = mapped_column(String(32), unique=True, index=True, nullable=False)

    client_id: Mapped[int] = mapped_column(ForeignKey("clients.id"), nullable=False)
    car_id: Mapped[int] = mapped_column(ForeignKey("cars.id"), nullable=False)
    reservation_id: Mapped[int | None] = mapped_column(ForeignKey("reservations.id"))
    employee_id: Mapped[int | None] = mapped_column(ForeignKey("users.id"))

    rental_start: Mapped[date] = mapped_column(Date, nullable=False)
    rental_end: Mapped[date] = mapped_column(Date, nullable=False)

    deposit: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    daily_price: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    total_price: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)

    start_mileage: Mapped[int | None] = mapped_column(Integer)
    pdf_path: Mapped[str | None] = mapped_column(String(255))
    notes: Mapped[str | None] = mapped_column(Text)

    status: Mapped[RentalStatus] = mapped_column(
        Enum(RentalStatus, native_enum=False, length=32),
        default=RentalStatus.ACTIVE,
        nullable=False,
    )

    client: Mapped["Client"] = relationship(back_populates="rentals")  # noqa: F821
    car: Mapped["Car"] = relationship(back_populates="rentals")  # noqa: F821
    employee: Mapped["User | None"] = relationship()  # noqa: F821
    payments: Mapped[list["Payment"]] = relationship(back_populates="rental")  # noqa: F821
    vehicle_return: Mapped["VehicleReturn | None"] = relationship(
        back_populates="rental", uselist=False, cascade="all, delete-orphan"
    )

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        data = super().to_dict(exclude=exclude)
        data["status_label"] = RentalStatus(self.status).label
        data["status"] = self.status.value if isinstance(self.status, RentalStatus) else self.status
        if self.client:
            data["client_name"] = self.client.full_name
        if self.car:
            data["car_name"] = self.car.display_name
        if self.employee:
            data["employee_name"] = self.employee.full_name
        data["has_return"] = self.vehicle_return is not None
        return data


class VehicleReturn(BaseModel):
    """Records the state of a car when it is returned and final charges."""

    __tablename__ = "vehicle_returns"

    rental_id: Mapped[int] = mapped_column(
        ForeignKey("rentals.id", ondelete="CASCADE"), unique=True, nullable=False
    )

    return_date: Mapped[datetime] = mapped_column(DateTime, nullable=False)
    mileage: Mapped[int] = mapped_column(Integer, default=0, nullable=False)
    fuel_level: Mapped[int] = mapped_column(Integer, default=100, nullable=False)  # percent
    damages: Mapped[str | None] = mapped_column(Text)

    extra_days: Mapped[int] = mapped_column(Integer, default=0, nullable=False)
    late_fee: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    damage_cost: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    penalties: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    final_payment: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)

    rental: Mapped["Rental"] = relationship(back_populates="vehicle_return")
