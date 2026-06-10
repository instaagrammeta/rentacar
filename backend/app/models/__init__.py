"""Domain models package.

Importing this package makes every model available to SQLAlchemy's metadata
(required for ``create_all`` and Alembic autogenerate).
"""
from __future__ import annotations

from app.models.accident import Accident, AccidentPhoto
from app.models.audit import AuditLog
from app.models.blacklist import BlacklistEntry
from app.models.car import Car, CarPhoto
from app.models.client import Client
from app.models.enums import (
    BlacklistReason,
    CarStatus,
    ClientStatus,
    PaymentMethod,
    PaymentType,
    RentalStatus,
    ReservationStatus,
    UserRole,
)
from app.models.payment import Payment
from app.models.rental import Rental, VehicleReturn
from app.models.reservation import Reservation
from app.models.settings import CompanySettings
from app.models.user import User

__all__ = [
    "Accident",
    "AccidentPhoto",
    "AuditLog",
    "BlacklistEntry",
    "Car",
    "CarPhoto",
    "Client",
    "CompanySettings",
    "Payment",
    "Rental",
    "Reservation",
    "User",
    "VehicleReturn",
    "BlacklistReason",
    "CarStatus",
    "ClientStatus",
    "PaymentMethod",
    "PaymentType",
    "RentalStatus",
    "ReservationStatus",
    "UserRole",
]
