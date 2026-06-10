"""Client management endpoints (Modules 2 & 3)."""
from __future__ import annotations

from flask import Blueprint, jsonify, request

from app.api.helpers import current_user, get_pagination, json_body
from app.models import Client
from app.models.enums import ClientStatus
from app.repositories import client_repo
from app.services import client_service
from app.utils.security import permission_required

bp = Blueprint("clients", __name__, url_prefix="/api/clients")


@bp.get("")
@permission_required("clients")
def list_clients():
    page, per_page, offset = get_pagination()
    search = request.args.get("search", "").strip()
    status = request.args.get("status", "").strip()

    filters = []
    if status:
        filters.append(Client.status == ClientStatus(status))

    if search:
        items = client_repo.search(search)
        total = len(items)
        items = items[offset : offset + per_page]
    else:
        total = client_repo.count(filters or None)
        items = client_repo.list(
            filters=filters or None, limit=per_page, offset=offset, order_by=Client.last_name.asc()
        )
    return jsonify(
        {
            "items": [c.to_dict() for c in items],
            "total": total,
            "page": page,
            "per_page": per_page,
        }
    )


@bp.get("/<int:client_id>")
@permission_required("clients")
def get_client(client_id: int):
    client = client_repo.get_or_404(client_id)
    return jsonify(client.to_dict())


@bp.get("/<int:client_id>/history")
@permission_required("clients")
def client_history(client_id: int):
    return jsonify(client_service.get_history(client_id))


@bp.get("/search")
@permission_required("clients")
def search_client():
    """VIP / returning client lookup by phone, QR/client code, or free text."""
    phone = request.args.get("phone")
    code = request.args.get("code")
    term = request.args.get("term")
    result = client_service.find_client(phone=phone, code=code, term=term)
    if isinstance(result, list):
        return jsonify([c.to_dict() for c in result])
    return jsonify(result.to_dict())


@bp.post("")
@permission_required("clients")
def create_client():
    client = client_service.create_client(json_body(), actor=current_user())
    return jsonify(client.to_dict()), 201


@bp.put("/<int:client_id>")
@permission_required("clients")
def update_client(client_id: int):
    client = client_service.update_client(client_id, json_body(), actor=current_user())
    return jsonify(client.to_dict())


@bp.delete("/<int:client_id>")
@permission_required("clients")
def delete_client(client_id: int):
    client_service.delete_client(client_id, actor=current_user())
    return jsonify({"message": "Клиент удалён"})
