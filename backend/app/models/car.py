"""Car / vehicle model and its photo gallery."""
from __future__ import annotations

from sqlalchemy import Enum, Float, ForeignKey, Integer, String
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel
from app.models.enums import CarStatus


class Car(BaseModel):
    """A vehicle available for rental."""

    __tablename__ = "cars"

    brand: Mapped[str] = mapped_column(String(80), nullable=False, index=True)
    model: Mapped[str] = mapped_column(String(80), nullable=False, index=True)
    year: Mapped[int] = mapped_column(Integer, nullable=False)
    color: Mapped[str | None] = mapped_column(String(40))
    vin: Mapped[str | None] = mapped_column(String(32), unique=True, index=True)
    plate_number: Mapped[str] = mapped_column(String(20), unique=True, index=True, nullable=False)
    mileage: Mapped[int] = mapped_column(Integer, default=0, nullable=False)

    daily_price: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    weekly_price: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    monthly_price: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)
    deposit_amount: Mapped[float] = mapped_column(Float, default=0.0, nullable=False)

    status: Mapped[CarStatus] = mapped_column(
        Enum(CarStatus, native_enum=False, length=32),
        default=CarStatus.AVAILABLE,
        nullable=False,
    )

    # Relationships
    photos: Mapped[list["CarPhoto"]] = relationship(
        back_populates="car", cascade="all, delete-orphan"
    )
    reservations: Mapped[list["Reservation"]] = relationship(back_populates="car")  # noqa: F821
    rentals: Mapped[list["Rental"]] = relationship(back_populates="car")  # noqa: F821
    accidents: Mapped[list["Accident"]] = relationship(back_populates="car")  # noqa: F821

    @property
    def display_name(self) -> str:
        return f"{self.brand} {self.model} ({self.plate_number})"

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        data = super().to_dict(exclude=exclude)
        data["display_name"] = self.display_name
        data["status"] = self.status.value if isinstance(self.status, CarStatus) else self.status
        data["status_label"] = CarStatus(self.status).label
        data["photos"] = [p.file_path for p in self.photos]
        return data


class CarPhoto(BaseModel):
    """A single photo belonging to a car."""

    __tablename__ = "car_photos"

    car_id: Mapped[int] = mapped_column(ForeignKey("cars.id", ondelete="CASCADE"), nullable=False)
    file_path: Mapped[str] = mapped_column(String(255), nullable=False)

    car: Mapped["Car"] = relationship(back_populates="photos")
