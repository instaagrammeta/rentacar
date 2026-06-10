"""Enumerations used across the domain models.

The ``value`` of each member is the canonical machine value stored in the
database. Human readable Russian labels are provided through :meth:`label` and
are primarily used for exports (Excel / PDF) where the backend renders text.
The Vue frontend additionally localises these values through Vue I18n.
"""
from __future__ import annotations

import enum


class UserRole(str, enum.Enum):
    """System roles controlling access to modules."""

    ADMINISTRATOR = "administrator"
    RENTAL_MANAGER = "rental_manager"
    CASHIER = "cashier"
    OPERATOR = "operator"

    @property
    def label(self) -> str:
        return {
            "administrator": "Администратор",
            "rental_manager": "Менеджер аренды",
            "cashier": "Кассир",
            "operator": "Оператор",
        }[self.value]


class ClientStatus(str, enum.Enum):
    ACTIVE = "active"
    PENDING_VERIFICATION = "pending_verification"
    BLACKLISTED = "blacklisted"

    @property
    def label(self) -> str:
        return {
            "active": "Активен",
            "pending_verification": "Ожидает проверки",
            "blacklisted": "В чёрном списке",
        }[self.value]


class CarStatus(str, enum.Enum):
    AVAILABLE = "available"
    RESERVED = "reserved"
    RENTED = "rented"
    MAINTENANCE = "maintenance"

    @property
    def label(self) -> str:
        return {
            "available": "Доступен",
            "reserved": "Забронирован",
            "rented": "В аренде",
            "maintenance": "На обслуживании",
        }[self.value]


class ReservationStatus(str, enum.Enum):
    RESERVED = "reserved"
    CONFIRMED = "confirmed"
    CANCELLED = "cancelled"

    @property
    def label(self) -> str:
        return {
            "reserved": "Забронировано",
            "confirmed": "Подтверждено",
            "cancelled": "Отменено",
        }[self.value]


class RentalStatus(str, enum.Enum):
    ACTIVE = "active"
    COMPLETED = "completed"
    CANCELLED = "cancelled"

    @property
    def label(self) -> str:
        return {
            "active": "Активна",
            "completed": "Завершена",
            "cancelled": "Отменена",
        }[self.value]


class PaymentMethod(str, enum.Enum):
    CASH = "cash"
    BANK_TRANSFER = "bank_transfer"
    CARD = "card"

    @property
    def label(self) -> str:
        return {
            "cash": "Наличные",
            "bank_transfer": "Банковский перевод",
            "card": "Карта",
        }[self.value]


class PaymentType(str, enum.Enum):
    DEPOSIT = "deposit"
    RENTAL = "rental"
    PENALTY = "penalty"
    DAMAGE = "damage"
    REFUND = "refund"

    @property
    def label(self) -> str:
        return {
            "deposit": "Депозит",
            "rental": "Оплата аренды",
            "penalty": "Штраф",
            "damage": "Возмещение ущерба",
            "refund": "Возврат средств",
        }[self.value]


class BlacklistReason(str, enum.Enum):
    FRAUD = "fraud"
    VEHICLE_DAMAGE = "vehicle_damage"
    NON_PAYMENT = "non_payment"
    SERIOUS_VIOLATION = "serious_violation"

    @property
    def label(self) -> str:
        return {
            "fraud": "Мошенничество",
            "vehicle_damage": "Повреждение автомобиля",
            "non_payment": "Неоплата",
            "serious_violation": "Серьёзное нарушение",
        }[self.value]
