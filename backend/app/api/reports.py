"""Reporting endpoints (Module 11)."""
from __future__ import annotations

from pathlib import Path

from flask import Blueprint, jsonify, request, send_file

from app.services import report_service
from app.utils.security import permission_required

bp = Blueprint("reports", __name__, url_prefix="/api/reports")


@bp.get("/daily-revenue")
@permission_required("reports")
def daily_revenue():
    day = request.args.get("date")
    from datetime import date as date_cls

    parsed = date_cls.fromisoformat(day) if day else None
    return jsonify(report_service.daily_revenue(parsed))


@bp.get("/monthly-revenue")
@permission_required("reports")
def monthly_revenue():
    from datetime import date as date_cls

    year = request.args.get("year", default=date_cls.today().year, type=int)
    month = request.args.get("month", default=date_cls.today().month, type=int)
    return jsonify(report_service.monthly_revenue(year, month))


@bp.get("/yearly-revenue")
@permission_required("reports")
def yearly_revenue():
    from datetime import date as date_cls

    year = request.args.get("year", default=date_cls.today().year, type=int)
    return jsonify(report_service.yearly_revenue(year))


@bp.get("/profitable-cars")
@permission_required("reports")
def profitable_cars():
    limit = request.args.get("limit", default=10, type=int)
    return jsonify(report_service.most_profitable_cars(limit))


@bp.get("/active-rentals")
@permission_required("reports")
def active_rentals():
    return jsonify(report_service.active_rentals())


@bp.get("/debtors")
@permission_required("reports")
def debtors():
    return jsonify(report_service.debtors())


@bp.get("/client-statistics")
@permission_required("reports")
def client_statistics():
    return jsonify(report_service.client_statistics())


@bp.get("/export/<report_type>")
@permission_required("reports")
def export_report(report_type: str):
    """Generate a styled Excel export for the given report type and download it."""
    params = request.args.to_dict()
    path = report_service.export_report_excel(report_type, params)
    file = Path(path)
    return send_file(file, as_attachment=True, download_name=file.name)
