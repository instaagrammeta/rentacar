"""Company settings model (Module 12), used in PDFs and exports."""
from __future__ import annotations

from sqlalchemy import String, Text
from sqlalchemy.orm import Mapped, mapped_column

from app.models.base import BaseModel


class CompanySettings(BaseModel):
    """Single-row table that stores company information."""

    __tablename__ = "company_settings"

    company_name: Mapped[str] = mapped_column(String(160), default="Rentacar CRM", nullable=False)
    address: Mapped[str | None] = mapped_column(Text)
    phone: Mapped[str | None] = mapped_column(String(64))
    email: Mapped[str | None] = mapped_column(String(120))
    logo_path: Mapped[str | None] = mapped_column(String(255))
    currency: Mapped[str] = mapped_column(String(8), default="RUB", nullable=False)
    # Free-form extra text printed on contracts (terms & conditions).
    contract_terms: Mapped[str | None] = mapped_column(Text)
