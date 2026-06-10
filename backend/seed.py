"""Seed the database with default users, settings and demo data.

Usage::

    python seed.py            # create admin + demo data (idempotent)
    python seed.py --minimal  # only the admin account and settings

The default administrator credentials are ``admin`` / ``admin123``. Change the
password after the first login.
"""
from __future__ import annotations

import argparse
from datetime import date, timedelta

from app import create_app
from app.extensions import db
from app.models import Car, Client, CompanySettings, User
from app.models.enums import CarStatus, ClientStatus, UserRole
from app.repositories import user_repo
from app.services import auth_service
from app.utils.security import hash_password


def ensure_admin() -> User:
    admin = user_repo.get_by_username("admin")
    if admin is None:
        admin = User(
            username="admin",
            full_name="Администратор системы",
            email="admin@rentacar.local",
            role=UserRole.ADMINISTRATOR,
            password_hash=hash_password("admin123"),
        )
        db.session.add(admin)
        db.session.commit()
        print("Создан администратор: admin / admin123")
    return admin


def ensure_demo_users() -> None:
    demo = [
        ("manager", "Менеджер Иванов", UserRole.RENTAL_MANAGER),
        ("cashier", "Кассир Петрова", UserRole.CASHIER),
        ("operator", "Оператор Сидоров", UserRole.OPERATOR),
    ]
    for username, full_name, role in demo:
        if user_repo.get_by_username(username) is None:
            auth_service.create_user(
                username=username,
                password=f"{username}123",
                full_name=full_name,
                role=role.value,
            )
            print(f"Создан пользователь: {username} / {username}123 ({role.label})")


def ensure_settings() -> None:
    if db.session.query(CompanySettings).first() is None:
        db.session.add(
            CompanySettings(
                company_name="ООО «Рентакар»",
                address="г. Москва, ул. Примерная, д. 1",
                phone="+7 (495) 000-00-00",
                email="info@rentacar.local",
                currency="RUB",
                contract_terms=(
                    "Арендатор обязуется бережно использовать транспортное средство и "
                    "вернуть его в надлежащем состоянии в указанный срок."
                ),
            )
        )
        db.session.commit()
        print("Созданы настройки компании")


def ensure_demo_data() -> None:
    if db.session.query(Car).count() == 0:
        cars = [
            Car(brand="Toyota", model="Camry", year=2022, color="Чёрный", plate_number="А001АА777",
                vin="JT2BF22K1W0123456", mileage=25000, daily_price=3500, weekly_price=21000,
                monthly_price=78000, deposit_amount=15000, status=CarStatus.AVAILABLE),
            Car(brand="Kia", model="Rio", year=2021, color="Белый", plate_number="В002ВВ777",
                vin="KNADN512AB6123456", mileage=40000, daily_price=2000, weekly_price=12000,
                monthly_price=45000, deposit_amount=8000, status=CarStatus.AVAILABLE),
            Car(brand="BMW", model="X5", year=2023, color="Синий", plate_number="С003СС777",
                vin="5UXKR0C50J0123456", mileage=12000, daily_price=8000, weekly_price=49000,
                monthly_price=185000, deposit_amount=40000, status=CarStatus.AVAILABLE),
        ]
        db.session.add_all(cars)
        db.session.commit()
        print(f"Создано автомобилей: {len(cars)}")

    if db.session.query(Client).count() == 0:
        clients = [
            Client(client_code="CL-000001", first_name="Алексей", last_name="Смирнов",
                   phone="+79161112233", email="smirnov@example.com",
                   date_of_birth=date.today() - timedelta(days=365 * 30),
                   driver_license_number="7700123456", driver_experience_years=8,
                   status=ClientStatus.ACTIVE, is_vip=True),
            Client(client_code="CL-000002", first_name="Мария", last_name="Кузнецова",
                   phone="+79162223344", email="kuznetsova@example.com",
                   date_of_birth=date.today() - timedelta(days=365 * 25),
                   driver_license_number="7700654321", driver_experience_years=3,
                   status=ClientStatus.ACTIVE),
        ]
        db.session.add_all(clients)
        db.session.commit()
        print(f"Создано клиентов: {len(clients)}")


def main() -> None:
    parser = argparse.ArgumentParser(description="Seed the Rentacar database")
    parser.add_argument("--minimal", action="store_true", help="Only admin + settings")
    args = parser.parse_args()

    app = create_app()
    with app.app_context():
        db.create_all()
        ensure_admin()
        ensure_settings()
        if not args.minimal:
            ensure_demo_users()
            ensure_demo_data()
    print("Готово. База данных инициализирована.")


if __name__ == "__main__":
    main()
