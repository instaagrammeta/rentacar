"""Background scheduler for automatic daily backups."""
from __future__ import annotations

import atexit
import logging

from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.cron import CronTrigger

from flask import Flask

logger = logging.getLogger("rentacar.scheduler")

_scheduler: BackgroundScheduler | None = None


def init_scheduler(app: Flask) -> None:
    """Start a background scheduler that creates a database backup daily."""
    global _scheduler
    if not app.config.get("ENABLE_SCHEDULED_BACKUPS"):
        return
    if _scheduler is not None:  # already initialised (e.g. reloader)
        return

    scheduler = BackgroundScheduler(daemon=True)

    def _job() -> None:
        from app.services import backup_service

        with app.app_context():
            try:
                backup_service.create_backup()
            except Exception:  # pragma: no cover - defensive
                logger.exception("Не удалось создать автоматическую резервную копию")

    # Every day at 02:00 local time.
    scheduler.add_job(_job, CronTrigger(hour=2, minute=0), id="daily_backup", replace_existing=True)
    scheduler.start()
    logger.info("Планировщик резервного копирования запущен")

    atexit.register(lambda: scheduler.shutdown(wait=False))
    _scheduler = scheduler
