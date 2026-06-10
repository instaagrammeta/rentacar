"""Company settings service (Module 12)."""
from __future__ import annotations

from typing import Any

from app.models import CompanySettings
from app.repositories import settings_repo
from app.services import audit_service


def get_settings() -> CompanySettings:
    return settings_repo.get_settings()


def update_settings(data: dict[str, Any], actor=None) -> CompanySettings:
    settings = settings_repo.get_settings()
    for field in (
        "company_name",
        "address",
        "phone",
        "email",
        "logo_path",
        "currency",
        "contract_terms",
    ):
        if field in data and data[field] is not None:
            setattr(settings, field, data[field])
    settings_repo.update(settings, {})
    audit_service.record(
        "update_settings",
        user_id=getattr(actor, "id", None),
        username=getattr(actor, "username", None),
        entity="settings",
        entity_id=settings.id,
    )
    return settings
