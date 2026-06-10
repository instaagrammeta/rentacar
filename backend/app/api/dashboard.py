"""Dashboard endpoints (Module 1)."""
from __future__ import annotations

from flask import Blueprint, jsonify, request

from app.services import dashboard_service
from app.utils.security import permission_required

bp = Blueprint("dashboard", __name__, url_prefix="/api/dashboard")


@bp.get("/summary")
@permission_required("dashboard")
def summary():
    return jsonify(dashboard_service.get_summary())


@bp.get("/revenue-by-month")
@permission_required("dashboard")
def revenue_by_month():
    months = request.args.get("months", default=12, type=int)
    return jsonify(dashboard_service.revenue_by_month(months))


@bp.get("/top-cars")
@permission_required("dashboard")
def top_cars():
    limit = request.args.get("limit", default=5, type=int)
    return jsonify(dashboard_service.top_rented_cars(limit))


@bp.get("/rental-statistics")
@permission_required("dashboard")
def rental_statistics():
    return jsonify(dashboard_service.rental_statistics())
