"""Development / production entry point for the Rentacar CRM backend.

Run with the Flask development server::

    python run.py

For production (used by the desktop build) the application is served by
Waitress, a pure-Python WSGI server that works well on Windows::

    python run.py --serve
"""
from __future__ import annotations

import argparse
import os

from app import create_app

app = create_app(os.environ.get("RENTACAR_ENV"))


def main() -> None:
    parser = argparse.ArgumentParser(description="Rentacar CRM backend server")
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=5000)
    parser.add_argument(
        "--serve", action="store_true", help="Serve with Waitress (production)"
    )
    args = parser.parse_args()

    if args.serve:
        from waitress import serve

        serve(app, host=args.host, port=args.port)
    else:
        app.run(host=args.host, port=args.port, debug=app.config.get("DEBUG", False))


if __name__ == "__main__":
    main()
