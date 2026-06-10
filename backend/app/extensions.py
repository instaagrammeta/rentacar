"""Flask extension singletons.

Extensions are instantiated here without an application so they can be imported
across the project without creating circular imports. They are bound to the
Flask app inside the application factory (:func:`app.create_app`).
"""
from __future__ import annotations

from flask_cors import CORS
from flask_jwt_extended import JWTManager
from flask_migrate import Migrate
from flask_sqlalchemy import SQLAlchemy

# ORM / database
db = SQLAlchemy()

# Database migrations (Alembic wrapper)
migrate = Migrate()

# JWT authentication
jwt = JWTManager()

# Cross-Origin Resource Sharing (needed for the Vite dev server / Electron)
cors = CORS()
