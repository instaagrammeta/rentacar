"""Payment service (Module 8)."""
from __future__ import annotations

from datetime import date, datetime
from typing import Any

from app.models import Payment
from app.models.enums import PaymentMethod, PaymentType
from app.repositories import client_repo, payment_repo, rental_repo, settings_repo
from app.services import audit_service, pdf_service
from app.utils.errors import ValidationError


def generate_receipt_number() -> str:
    last = payment_repo.list(limit=1)
    next_id = (last[0].id + 1) if last else 1
    return f"KV-{date.today().year}-{next_id:05d}"


def create_payment(data: dict[str, Any], actor=None) -> Payment:
    client = client_repo.get_or_404(int(data["client_id"]))
    amount = float(data.get("amount") or 0)
    if amount <= 0:
        raise ValidationError("Сумма платежа должна быть больше нуля")

    rental_id = data.get("rental_id")
    if rental_id:
        rental_repo.get_or_404(int(rental_id))

    try:
        method = PaymentMethod(data.get("method", PaymentMethod.CASH.value))
        payment_type = PaymentType(data.get("payment_type", PaymentType.RENTAL.value))
    except ValueError as exc:
        raise ValidationError("Некорректный способ или тип оплаты") from exc

    payment = Payment(
        receipt_number=generate_receipt_number(),
        client_id=client.id,
        rental_id=int(rental_id) if rental_id else None,
        cashier_id=getattr(actor, "id", None),
        amount=amount,
        method=method,
        payment_type=payment_type,
        paid_at=datetime.now(),
        notes=data.get("notes"),
    )
    payment_repo.add(payment)

    settings = settings_repo.get_settings()
    payment.receipt_path = pdf_service.generate_receipt_pdf(payment, settings)
    payment_repo.update(payment, {})

    audit_service.record(
        "create_payment",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="payment",
        entity_id=payment.id,
        details=f"Платёж {payment.receipt_number} на сумму {amount:.2f}",
    )
    return payment


def list_payments(*, client_id: int | None = None, rental_id: int | None = None) -> list[Payment]:
    filters = []
    if client_id:
        filters.append(Payment.client_id == client_id)
    if rental_id:
        filters.append(Payment.rental_id == rental_id)
    return list(payment_repo.list(filters=filters or None))
