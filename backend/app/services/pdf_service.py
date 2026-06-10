"""PDF generation for rental contracts and payment receipts.

Uses ReportLab. Cyrillic text is supported by registering a Unicode TrueType
font when available, otherwise ReportLab's built-in Helvetica is used (which
still renders Cyrillic when the platform provides the glyphs).
"""
from __future__ import annotations

import logging
from pathlib import Path

from reportlab.lib import colors
from reportlab.lib.pagesizes import A4
from reportlab.lib.styles import ParagraphStyle, getSampleStyleSheet
from reportlab.lib.units import mm
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.platypus import (
    Image,
    Paragraph,
    SimpleDocTemplate,
    Spacer,
    Table,
    TableStyle,
)

from flask import current_app

logger = logging.getLogger("rentacar.pdf")

GREEN = colors.HexColor("#1B5E20")
LIGHT_GREEN = colors.HexColor("#E8F5E9")

_FONT_NAME = "Helvetica"
_FONT_BOLD = "Helvetica-Bold"
_FONT_REGISTERED = False

# Common Linux/Windows locations for a Unicode font that includes Cyrillic.
_FONT_CANDIDATES = [
    "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
    "/usr/share/fonts/dejavu/DejaVuSans.ttf",
    "C:/Windows/Fonts/arial.ttf",
    "C:/Windows/Fonts/DejaVuSans.ttf",
]
_FONT_BOLD_CANDIDATES = [
    "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
    "/usr/share/fonts/dejavu/DejaVuSans-Bold.ttf",
    "C:/Windows/Fonts/arialbd.ttf",
]


def _ensure_font() -> None:
    """Register a Cyrillic-capable TrueType font if one is available."""
    global _FONT_NAME, _FONT_BOLD, _FONT_REGISTERED
    if _FONT_REGISTERED:
        return
    for path in _FONT_CANDIDATES:
        if Path(path).exists():
            try:
                pdfmetrics.registerFont(TTFont("AppFont", path))
                _FONT_NAME = "AppFont"
                break
            except Exception:  # pragma: no cover - defensive
                continue
    for path in _FONT_BOLD_CANDIDATES:
        if Path(path).exists():
            try:
                pdfmetrics.registerFont(TTFont("AppFont-Bold", path))
                _FONT_BOLD = "AppFont-Bold"
                break
            except Exception:  # pragma: no cover - defensive
                continue
    _FONT_REGISTERED = True


def _styles() -> dict:
    _ensure_font()
    base = getSampleStyleSheet()
    title = ParagraphStyle(
        "AppTitle",
        parent=base["Title"],
        fontName=_FONT_BOLD,
        textColor=GREEN,
        fontSize=18,
    )
    heading = ParagraphStyle(
        "AppHeading", parent=base["Heading2"], fontName=_FONT_BOLD, textColor=GREEN
    )
    normal = ParagraphStyle("AppNormal", parent=base["Normal"], fontName=_FONT_NAME, fontSize=10)
    return {"title": title, "heading": heading, "normal": normal}


def _company_header(elements: list, styles: dict, settings) -> None:
    if settings and settings.logo_path:
        logo = current_app.config["UPLOAD_DIR"] / settings.logo_path
        if logo.exists():
            try:
                elements.append(Image(str(logo), width=40 * mm, height=20 * mm))
            except Exception:  # pragma: no cover
                pass
    company = settings.company_name if settings else "Rentacar CRM"
    elements.append(Paragraph(company, styles["title"]))
    if settings and settings.address:
        elements.append(Paragraph(settings.address, styles["normal"]))
    if settings and settings.phone:
        elements.append(Paragraph(f"Телефон: {settings.phone}", styles["normal"]))
    elements.append(Spacer(1, 8 * mm))


def _table(data: list[list[str]]) -> Table:
    table = Table(data, hAlign="LEFT", colWidths=[60 * mm, 110 * mm])
    table.setStyle(
        TableStyle(
            [
                ("FONTNAME", (0, 0), (-1, -1), _FONT_NAME),
                ("FONTSIZE", (0, 0), (-1, -1), 10),
                ("BACKGROUND", (0, 0), (0, -1), LIGHT_GREEN),
                ("TEXTCOLOR", (0, 0), (0, -1), GREEN),
                ("GRID", (0, 0), (-1, -1), 0.5, colors.HexColor("#C8E6C9")),
                ("VALIGN", (0, 0), (-1, -1), "TOP"),
                ("PADDING", (0, 0), (-1, -1), 6),
            ]
        )
    )
    return table


def generate_contract_pdf(rental, settings) -> str:
    """Generate a rental contract PDF and return its relative path."""
    styles = _styles()
    export_dir: Path = current_app.config["EXPORT_DIR"] / "contracts"
    export_dir.mkdir(parents=True, exist_ok=True)
    relative = f"contracts/{rental.contract_number}.pdf"
    full_path = current_app.config["EXPORT_DIR"] / relative

    doc = SimpleDocTemplate(str(full_path), pagesize=A4, title=f"Договор {rental.contract_number}")
    elements: list = []
    _company_header(elements, styles, settings)

    elements.append(Paragraph(f"Договор аренды № {rental.contract_number}", styles["heading"]))
    elements.append(Spacer(1, 4 * mm))

    client = rental.client
    car = rental.car
    rows = [
        ["Клиент", client.full_name if client else "—"],
        ["Телефон", client.phone if client else "—"],
        ["Вод. удостоверение", (client.driver_license_number or "—") if client else "—"],
        ["Автомобиль", car.display_name if car else "—"],
        ["VIN", (car.vin or "—") if car else "—"],
        ["Начало аренды", rental.rental_start.isoformat()],
        ["Окончание аренды", rental.rental_end.isoformat()],
        ["Цена за сутки", f"{rental.daily_price:.2f}"],
        ["Депозит", f"{rental.deposit:.2f}"],
        ["Итоговая стоимость", f"{rental.total_price:.2f}"],
        ["Ответственный сотрудник", rental.employee.full_name if rental.employee else "—"],
    ]
    elements.append(_table(rows))
    elements.append(Spacer(1, 8 * mm))

    if settings and settings.contract_terms:
        elements.append(Paragraph("Условия договора", styles["heading"]))
        elements.append(Paragraph(settings.contract_terms, styles["normal"]))
        elements.append(Spacer(1, 8 * mm))

    elements.append(Spacer(1, 16 * mm))
    sign = _table([["Подпись клиента", "______________________"], ["Подпись сотрудника", "______________________"]])
    elements.append(sign)

    doc.build(elements)
    logger.info("Создан PDF договора: %s", full_path)
    return relative


def generate_receipt_pdf(payment, settings) -> str:
    """Generate a payment receipt PDF and return its relative path."""
    styles = _styles()
    export_dir: Path = current_app.config["EXPORT_DIR"] / "receipts"
    export_dir.mkdir(parents=True, exist_ok=True)
    relative = f"receipts/{payment.receipt_number}.pdf"
    full_path = current_app.config["EXPORT_DIR"] / relative

    doc = SimpleDocTemplate(str(full_path), pagesize=A4, title=f"Квитанция {payment.receipt_number}")
    elements: list = []
    _company_header(elements, styles, settings)
    elements.append(Paragraph(f"Квитанция № {payment.receipt_number}", styles["heading"]))
    elements.append(Spacer(1, 4 * mm))

    rows = [
        ["Клиент", payment.client.full_name if payment.client else "—"],
        ["Дата оплаты", payment.paid_at.strftime("%d.%m.%Y %H:%M")],
        ["Тип платежа", payment.payment_type.label if hasattr(payment.payment_type, "label") else str(payment.payment_type)],
        ["Способ оплаты", payment.method.label if hasattr(payment.method, "label") else str(payment.method)],
        ["Сумма", f"{payment.amount:.2f}"],
    ]
    if payment.notes:
        rows.append(["Примечание", payment.notes])
    elements.append(_table(rows))
    doc.build(elements)
    logger.info("Создан PDF квитанции: %s", full_path)
    return relative
