"""Blacklist service (Module 9)."""
from __future__ import annotations

from typing import Any

from app.models import BlacklistEntry
from app.models.enums import BlacklistReason, ClientStatus
from app.repositories import blacklist_repo, client_repo
from app.services import audit_service
from app.utils.errors import ConflictError, ValidationError


def add_to_blacklist(data: dict[str, Any], actor=None) -> BlacklistEntry:
    client = client_repo.get_or_404(int(data["client_id"]))
    if blacklist_repo.get_by_client(client.id):
        raise ConflictError("Клиент уже в чёрном списке")

    try:
        reason = BlacklistReason(data["reason"])
    except (KeyError, ValueError) as exc:
        raise ValidationError("Укажите корректную причину") from exc

    entry = BlacklistEntry(
        client_id=client.id,
        reason=reason,
        comment=data.get("comment"),
        created_by=getattr(actor, "id", None),
    )
    blacklist_repo.add(entry)

    client.status = ClientStatus.BLACKLISTED
    client_repo.update(client, {})

    audit_service.record(
        "blacklist_add",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="client",
        entity_id=client.id,
        details=f"Причина: {reason.label}",
    )
    return entry


def remove_from_blacklist(client_id: int, actor=None) -> None:
    entry = blacklist_repo.get_by_client(client_id)
    if entry is None:
        raise ConflictError("Клиент не находится в чёрном списке")
    blacklist_repo.delete(entry)

    client = client_repo.get(client_id)
    if client:
        client.status = ClientStatus.ACTIVE
        client_repo.update(client, {})

    audit_service.record(
        "blacklist_remove",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="client",
        entity_id=client_id,
    )


def list_blacklist() -> list[BlacklistEntry]:
    return list(blacklist_repo.list())
