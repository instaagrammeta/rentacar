# Database Migrations (Alembic / Flask-Migrate)

This directory contains the Alembic migration environment.

The application can run without migrations because `create_app()` calls
`db.create_all()` on start-up (convenient for the packaged desktop build).
For controlled schema evolution in team/production environments, use the
migration commands below.

## Commands

From the `backend/` directory with the virtual environment active:

```bash
export FLASK_APP=run.py            # Windows: set FLASK_APP=run.py

# Generate the initial migration from the models
flask db migrate -m "initial schema"

# Apply migrations
flask db upgrade

# Roll back the last migration
flask db downgrade -1
```

`versions/` will hold the generated revision files. `render_as_batch=True` is
enabled in `env.py` so SQLite `ALTER TABLE` operations work correctly.
