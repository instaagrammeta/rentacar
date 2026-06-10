"""Generic repository implementing the Repository Pattern.

The repository isolates SQLAlchemy query logic from the service layer so the
business code never talks to the ORM session directly. Concrete repositories
subclass :class:`BaseRepository` and set the ``model`` attribute.
"""
from __future__ import annotations

from typing import Generic, Sequence, TypeVar

from sqlalchemy import func, select

from app.extensions import db

T = TypeVar("T")


class BaseRepository(Generic[T]):
    """CRUD helper bound to a single model class."""

    model: type[T]

    def __init__(self, model: type[T] | None = None) -> None:
        if model is not None:
            self.model = model

    # ------------------------------------------------------------------ reads
    def get(self, entity_id: int) -> T | None:
        """Return a single entity by primary key or ``None``."""
        return db.session.get(self.model, entity_id)

    def get_or_404(self, entity_id: int) -> T:
        """Return an entity or raise a 404 error."""
        from app.utils.errors import NotFoundError

        entity = self.get(entity_id)
        if entity is None:
            raise NotFoundError(f"{self.model.__name__} #{entity_id} не найден")
        return entity

    def list(
        self,
        *,
        filters: list | None = None,
        order_by=None,
        limit: int | None = None,
        offset: int | None = None,
    ) -> Sequence[T]:
        """Return entities matching optional filters / ordering / pagination."""
        stmt = select(self.model)
        if filters:
            stmt = stmt.where(*filters)
        if order_by is not None:
            stmt = stmt.order_by(order_by)
        else:
            stmt = stmt.order_by(self.model.id.desc())
        if limit is not None:
            stmt = stmt.limit(limit)
        if offset is not None:
            stmt = stmt.offset(offset)
        return db.session.execute(stmt).scalars().all()

    def count(self, filters: list | None = None) -> int:
        """Count rows matching optional filters."""
        stmt = select(func.count()).select_from(self.model)
        if filters:
            stmt = stmt.where(*filters)
        return db.session.execute(stmt).scalar_one()

    # ----------------------------------------------------------------- writes
    def add(self, entity: T, *, commit: bool = True) -> T:
        db.session.add(entity)
        if commit:
            db.session.commit()
        else:
            db.session.flush()
        return entity

    def update(self, entity: T, data: dict, *, commit: bool = True) -> T:
        for key, value in data.items():
            if hasattr(entity, key):
                setattr(entity, key, value)
        if commit:
            db.session.commit()
        return entity

    def delete(self, entity: T, *, commit: bool = True) -> None:
        db.session.delete(entity)
        if commit:
            db.session.commit()
