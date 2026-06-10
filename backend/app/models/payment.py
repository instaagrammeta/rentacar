"""Payment / receipt model (Module 8)."""
from __future__ import annotations

from datetime import datetime

from sqlalchemy import DateTime, Enum, Float, ForeignKey, String, Text
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.models.base import BaseModel, utcnow
from app.models.enums import PaymentMethod, PaymentType


class Payment(BaseModel):
    """A money transaction linked to a client and optionally a rental."""

    __tablename__ = "payments"

    receipt_number: Mapped[str] = mapped_column(String(32), unique=True, index=True, nullable=False)

    client_id: Mapped[int] = mapped_column(ForeignKey("clients.id"), nullable=False)
    rental_id: Mapped[int | None] = mapped_column(ForeignKey("rentals.id"))
    cashier_id: Mapped[int | None] = mapped_column(ForeignKey("users.id"))

    amount: Mapped[float] = mapped_column(Float, nullable=False)
    method: Mapped[PaymentMethod] = mapped_column(
        Enum(PaymentMethod, native_enum=False, length=32), nullable=False
    )
    payment_type: Mapped[PaymentType] = mapped_column(
        Enum(PaymentType, native_enum=False, length=32),
        default=PaymentType.RENTAL,
        nullable=False,
    )
    paid_at: Mapped[datetime] = mapped_column(DateTime, default=utcnow, nullable=False)
    receipt_path: Mapped[str | None] = mapped_column(String(255))
    notes: Mapped[str | None] = mapped_column(Text)

    client: Mapped["Client"] = relationship(back_populates="payments")  # noqa: F821
    rental: Mapped["Rental | None"] = relationship(back_populates="payments")  # noqa: F821
    cashier: Mapped["User | None"] = relationship()  # noqa: F821

    def to_dict(self, exclude: set[str] | None = None) -> dict:
        data = super().to_dict(exclude=exclude)
        data["method"] = self.method.value if isinstance(self.method, PaymentMethod) else self.method
        data["method_label"] = PaymentMethod(self.method).label
        data["payment_type"] = (
            self.payment_type.value
            if isinstance(self.payment_type, PaymentType)
            else self.payment_type
        )
        data["payment_type_label"] = PaymentType(self.payment_type).label
        if self.client:
            data["client_name"] = self.client.full_name
        return data
