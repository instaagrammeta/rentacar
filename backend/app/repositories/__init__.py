"""Repository layer.

Each repository encapsulates persistence logic for one aggregate. Service
classes depend on repositories rather than on the ORM session directly.
"""
from __future__ import annotations

from sqlalchemy import or_, select

from app.extensions import db
from app.models import (
    Accident,
    AuditLog,
    BlacklistEntry,
    Car,
    Client,
    CompanySettings,
    Payment,
    Rental,
    Reservation,
    User,
)
from app.repositories.base_repository import BaseRepository


class UserRepository(BaseRepository[User]):
    model = User

    def get_by_username(self, username: str) -> User | None:
        stmt = select(User).where(User.username == username)
        return db.session.execute(stmt).scalar_one_or_none()


class ClientRepository(BaseRepository[Client]):
    model = Client

    def get_by_code(self, code: str) -> Client | None:
        stmt = select(Client).where(Client.client_code == code)
        return db.session.execute(stmt).scalar_one_or_none()

    def get_by_phone(self, phone: str) -> Client | None:
        stmt = select(Client).where(Client.phone == phone)
        return db.session.execute(stmt).scalars().first()

    def search(self, term: str) -> list[Client]:
        like = f"%{term}%"
        stmt = select(Client).where(
            or_(
                Client.first_name.ilike(like),
                Client.last_name.ilike(like),
                Client.phone.ilike(like),
                Client.client_code.ilike(like),
                Client.email.ilike(like),
            )
        )
        return list(db.session.execute(stmt).scalars().all())


class CarRepository(BaseRepository[Car]):
    model = Car


class ReservationRepository(BaseRepository[Reservation]):
    model = Reservation


class RentalRepository(BaseRepository[Rental]):
    model = Rental


class PaymentRepository(BaseRepository[Payment]):
    model = Payment


class AccidentRepository(BaseRepository[Accident]):
    model = Accident


class BlacklistRepository(BaseRepository[BlacklistEntry]):
    model = BlacklistEntry

    def get_by_client(self, client_id: int) -> BlacklistEntry | None:
        stmt = select(BlacklistEntry).where(BlacklistEntry.client_id == client_id)
        return db.session.execute(stmt).scalar_one_or_none()


class AuditRepository(BaseRepository[AuditLog]):
    model = AuditLog


class SettingsRepository(BaseRepository[CompanySettings]):
    model = CompanySettings

    def get_settings(self) -> CompanySettings:
        """Return the singleton settings row, creating it if missing."""
        settings = db.session.execute(select(CompanySettings)).scalars().first()
        if settings is None:
            settings = CompanySettings()
            db.session.add(settings)
            db.session.commit()
        return settings


# Convenient singletons
user_repo = UserRepository()
client_repo = ClientRepository()
car_repo = CarRepository()
reservation_repo = ReservationRepository()
rental_repo = RentalRepository()
payment_repo = PaymentRepository()
accident_repo = AccidentRepository()
blacklist_repo = BlacklistRepository()
audit_repo = AuditRepository()
settings_repo = SettingsRepository()

__all__ = [
    "UserRepository",
    "ClientRepository",
    "CarRepository",
    "ReservationRepository",
    "RentalRepository",
    "PaymentRepository",
    "AccidentRepository",
    "BlacklistRepository",
    "AuditRepository",
    "SettingsRepository",
    "user_repo",
    "client_repo",
    "car_repo",
    "reservation_repo",
    "rental_repo",
    "payment_repo",
    "accident_repo",
    "blacklist_repo",
    "audit_repo",
    "settings_repo",
]
