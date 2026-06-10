"""Dashboard aggregation service (Module 1)."""
from __future__ import annotations

from datetime import date, datetime, timedelta

from sqlalchemy import func, select

from app.extensions import db
from app.models import Car, Payment, Rental
from app.models.enums import CarStatus, RentalStatus


def _car_counts() -> dict[str, int]:
    stmt = select(Car.status, func.count(Car.id)).group_by(Car.status)
    counts = {status.value: 0 for status in CarStatus}
    for status, count in db.session.execute(stmt).all():
        key = status.value if isinstance(status, CarStatus) else status
        counts[key] = count
    return counts


def _revenue(start: datetime, end: datetime) -> float:
    stmt = select(func.coalesce(func.sum(Payment.amount), 0.0)).where(
        Payment.paid_at >= start, Payment.paid_at < end
    )
    return float(db.session.execute(stmt).scalar_one())


def get_summary() -> dict:
    """Return the headline dashboard metrics."""
    counts = _car_counts()

    today = date.today()
    day_start = datetime.combine(today, datetime.min.time())
    month_start = datetime.combine(today.replace(day=1), datetime.min.time())

    today_revenue = _revenue(day_start, day_start + timedelta(days=1))
    next_month = (month_start.replace(day=28) + timedelta(days=4)).replace(day=1)
    monthly_revenue = _revenue(month_start, next_month)

    active_rentals = db.session.execute(
        select(func.count(Rental.id)).where(Rental.status == RentalStatus.ACTIVE)
    ).scalar_one()

    upcoming_returns = db.session.execute(
        select(func.count(Rental.id)).where(
            Rental.status == RentalStatus.ACTIVE,
            Rental.rental_end <= today + timedelta(days=3),
        )
    ).scalar_one()

    return {
        "cars_available": counts.get(CarStatus.AVAILABLE.value, 0),
        "cars_rented": counts.get(CarStatus.RENTED.value, 0),
        "cars_reserved": counts.get(CarStatus.RESERVED.value, 0),
        "cars_maintenance": counts.get(CarStatus.MAINTENANCE.value, 0),
        "today_revenue": round(today_revenue, 2),
        "monthly_revenue": round(monthly_revenue, 2),
        "active_rentals": active_rentals,
        "upcoming_returns": upcoming_returns,
    }


def revenue_by_month(months: int = 12) -> list[dict]:
    """Return revenue totals grouped by calendar month."""
    today = date.today().replace(day=1)
    results: list[dict] = []
    cursor = today
    # Walk backwards `months` months.
    series: list[date] = []
    for _ in range(months):
        series.append(cursor)
        cursor = (cursor - timedelta(days=1)).replace(day=1)
    for month_start in reversed(series):
        next_month = (month_start.replace(day=28) + timedelta(days=4)).replace(day=1)
        amount = _revenue(
            datetime.combine(month_start, datetime.min.time()),
            datetime.combine(next_month, datetime.min.time()),
        )
        results.append({"month": month_start.strftime("%Y-%m"), "revenue": round(amount, 2)})
    return results


def top_rented_cars(limit: int = 5) -> list[dict]:
    """Return the most frequently rented cars."""
    stmt = (
        select(Car, func.count(Rental.id).label("rentals"))
        .join(Rental, Rental.car_id == Car.id)
        .group_by(Car.id)
        .order_by(func.count(Rental.id).desc())
        .limit(limit)
    )
    rows = db.session.execute(stmt).all()
    return [{"car": car.display_name, "rentals": count} for car, count in rows]


def rental_statistics() -> dict:
    """Return rental counts grouped by status."""
    stmt = select(Rental.status, func.count(Rental.id)).group_by(Rental.status)
    stats = {status.value: 0 for status in RentalStatus}
    for status, count in db.session.execute(stmt).all():
        key = status.value if isinstance(status, RentalStatus) else status
        stats[key] = count
    return stats
