"""Pricing helpers used by rentals and returns.

The pricing logic chooses the most favourable tariff for the customer: monthly
blocks are billed at the monthly price, remaining whole weeks at the weekly
price and the rest at the daily price. When weekly / monthly prices are not
configured the calculation gracefully falls back to the daily rate.
"""
from __future__ import annotations

from dataclasses import dataclass
from datetime import date


@dataclass
class PriceBreakdown:
    days: int
    months: int
    weeks: int
    extra_days: int
    total: float

    def to_dict(self) -> dict:
        return {
            "days": self.days,
            "months": self.months,
            "weeks": self.weeks,
            "extra_days": self.extra_days,
            "total": round(self.total, 2),
        }


def rental_days(start: date, end: date) -> int:
    """Return the number of billable days (inclusive of the first day)."""
    delta = (end - start).days
    return max(delta, 1)


def calculate_price(
    start: date,
    end: date,
    *,
    daily_price: float,
    weekly_price: float = 0.0,
    monthly_price: float = 0.0,
) -> PriceBreakdown:
    """Calculate the total rental price for a date range."""
    days = rental_days(start, end)
    remaining = days
    total = 0.0

    months = 0
    if monthly_price and monthly_price > 0:
        months = remaining // 30
        total += months * monthly_price
        remaining -= months * 30

    weeks = 0
    if weekly_price and weekly_price > 0:
        weeks = remaining // 7
        total += weeks * weekly_price
        remaining -= weeks * 7

    total += remaining * daily_price
    return PriceBreakdown(days=days, months=months, weeks=weeks, extra_days=remaining, total=total)
