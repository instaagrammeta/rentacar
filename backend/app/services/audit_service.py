"""Audit logging service (security activity trail)."""
from __future__ import annotations

import logging

from flask import has_request_context, request

from app.models import AuditLog
from app.repositories import audit_repo

logger = logging.getLogger("rentacar.audit")


def record(
    action: str,
    *,
    user_id: int | None = None,
    username: str | None = None,
    entity: str | None = None,
    entity_id: int | None = None,
    details: str | None = None,
) -> AuditLog:
    """Create an audit log entry.

    The function is defensive: failures to write the audit trail must never
    break the main business operation, therefore exceptions are swallowed and
    logged instead of propagated.
    """
    ip_address = None
    if has_request_context():
        ip_address = request.headers.get("X-Forwarded-For", request.remote_addr)

    entry = AuditLog(
        user_id=user_id,
        username=username,
        action=action,
        entity=entity,
        entity_id=entity_id,
        ip_address=ip_address,
        details=details,
    )
    try:
        audit_repo.add(entry)
        logger.info("AUDIT %s by %s on %s#%s", action, username, entity, entity_id)
    except Exception:  # pragma: no cover - defensive
        logger.exception("Не удалось записать журнал аудита")
    return entry


def list_logs(limit: int = 200, offset: int = 0) -> list[AuditLog]:
    return list(audit_repo.list(limit=limit, offset=offset))
