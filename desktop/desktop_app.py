"""Alternative single-process desktop launcher using PyWebView.

This is a lightweight Python-only alternative to the Electron build. It starts
the Flask backend (served by Waitress in a background thread) and opens a native
window pointing at it. The Vue frontend must be built (``npm run build``) and
its ``dist`` placed next to the backend, or served by Flask.

Run::

    python desktop_app.py

Package to an .exe with PyInstaller::

    pyinstaller --noconsole --name "Rentacar CRM" desktop_app.py
"""
from __future__ import annotations

import sys
import threading
import time
from pathlib import Path

# Ensure the backend package is importable when run from the project root.
BACKEND_DIR = Path(__file__).resolve().parent.parent / "backend"
sys.path.insert(0, str(BACKEND_DIR))

import webview  # noqa: E402  (import after sys.path tweak)
from waitress import serve  # noqa: E402

from app import create_app  # noqa: E402

HOST = "127.0.0.1"
PORT = 5000


def _run_server(app) -> None:
    serve(app, host=HOST, port=PORT, threads=8)


def main() -> None:
    app = create_app("production")

    server_thread = threading.Thread(target=_run_server, args=(app,), daemon=True)
    server_thread.start()

    # Give the server a moment to bind the port.
    time.sleep(1.5)

    webview.create_window(
        "Rentacar CRM",
        f"http://{HOST}:{PORT}",
        width=1440,
        height=900,
        min_size=(1024, 700),
    )
    webview.start()


if __name__ == "__main__":
    main()
