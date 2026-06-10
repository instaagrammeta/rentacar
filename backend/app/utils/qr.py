"""QR code generation helper used for VIP / returning client identification."""
from __future__ import annotations

from pathlib import Path

import qrcode

from flask import current_app


def generate_qr(data: str, filename: str) -> str:
    """Generate a QR code PNG for ``data`` and return its relative path."""
    upload_dir: Path = current_app.config["UPLOAD_DIR"] / "qr"
    upload_dir.mkdir(parents=True, exist_ok=True)

    qr = qrcode.QRCode(
        version=1,
        error_correction=qrcode.constants.ERROR_CORRECT_M,
        box_size=8,
        border=2,
    )
    qr.add_data(data)
    qr.make(fit=True)
    image = qr.make_image(fill_color="#1B5E20", back_color="white")
    relative = f"qr/{filename}"
    image.save(current_app.config["UPLOAD_DIR"] / relative)
    return relative
