# -*- mode: python ; coding: utf-8 -*-
"""PyInstaller spec for the Rentacar CRM backend.

Build the standalone backend executable with::

    pyinstaller rentacar-backend.spec

The resulting `dist/rentacar-backend(.exe)` is bundled by Electron Builder via
the `extraResources` entry in `desktop/package.json`.
"""
from PyInstaller.utils.hooks import collect_submodules, collect_data_files

hidden_imports = (
    collect_submodules("app")
    + collect_submodules("flask_sqlalchemy")
    + collect_submodules("flask_jwt_extended")
    + collect_submodules("openpyxl")
    + collect_submodules("reportlab")
    + collect_submodules("passlib")
    + ["waitress", "qrcode", "PIL", "apscheduler", "pandas"]
)

datas = collect_data_files("reportlab")

a = Analysis(
    ["run.py"],
    pathex=["."],
    binaries=[],
    datas=datas,
    hiddenimports=hidden_imports,
    hookspath=[],
    runtime_hooks=[],
    excludes=["tkinter"],
    noarchive=False,
)

pyz = PYZ(a.pure, a.zipped_data)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.zipfiles,
    a.datas,
    [],
    name="rentacar-backend",
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    console=True,
    disable_windowed_traceback=False,
    target_arch=None,
)
