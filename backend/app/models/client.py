"""Client (customer) model including VIP identification helpers."""
from __future__ import annotations

from datetime import date

from sqlalchemy import Date, Enum, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel
from app.models.enums import ClientStatus


class Client(BaseModel):
    """A rental customer with documents and verification state."""

    __tablename__ = "clients"

    # Unique business identifier used for VIP lookup / QR codes (e.g. CL-000123)
    client_code: Mapped[str] = mapped_column(String(32), unique=True, index=True, nullable=False)

    first_name: Mapped[str] = mapped_column(String(80), nullable=False)
    last_name: Mapped[str] = mapped_column(String(80), nullable=False)
    phone: Mapped[str] = mapped_column(String(32), index=True, nullable=False)
    email: Mapped[str | None] = mapped_column(String(120), index=True)
    date_of_birth: Mapped[date | None] = mapped_column(Date)

    passport_number: Mapped[str | None] = mapped_column(String(64), index=True)
    driver_license_number: Mapped[str | None] = mapped_column(String(64), index=True)
    driver_license_issue_date: Mapped[date | None] = mapped_column(Date)
    # Driving experience expressed in whole years.
    driver_experience_years: Mapped[int] = mapped_column(default=0, nullable=False)

    # File paths (relative to the upload directory) for scanned documents.
    passport_scan: Mapped[str | None] = mapped_column(String(255))
    driver_license_scan: Mapped[str | None] = mapped_column(String(255))

    # QR code image path generated for VIP / returning clients.
    qr_code_path: Mapped[str | None] = mapped_column(String(255))
    is_vip: Mapped[bool] = mapped_column(default=False, nullable=False)

    notes: Mapped[str | None] = mapped_column(Text)
    status: Mapped[ClientStatus] = mapped_column(
        Enum(ClientStatus, native_enum=False, length=32),
        default=ClientStatus.PENDING_VERIFICATION,
        nullable=False,
    )

    # Relationships
    reservations: Mapped[list["Reservation"]] = relationship(back_populates="client")  # noqa: F821
    rentals: Mapped[list["Rental"]] = relationship(back_populates="client")  # noqa: F821
    payments: Mapped[list["Payment"]] = relationship(back_populates="client")  # noqa: F821
    accidents: Mapped[list["Accident"]] = relationship(back_populates="client")  # noqa: F821
    blacklist_entry: Mapped["BlacklistEntry | None"] = relationship(  # noqa: F821
        back_populates="client", uselist=False
    )

    @property
    def full_name(self) -> str:
        return f"{self.last_name} {self.first_name}".strip()

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        data = super().to_dict(exclude=exclude)
        data["full_name"] = self.full_name
        data["status"] = self.status.value if isinstance(self.status, ClientStatus) else self.status
        data["status_label"] = ClientStatus(self.status).label
        return data
