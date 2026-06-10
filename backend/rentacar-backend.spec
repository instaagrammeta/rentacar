# -*- mode: python ; coding: utf-8 -*-
"""PyInstaller spec for the Rentacar CRM backend.

Build the standalone backend executable with::

    pyinstaller rentacar-backend.spec

This produces a *one-folder* build at ``dist/rentacar-backend/`` containing
``rentacar-backend.exe`` together with all dependencies. The whole folder is
staged into ``desktop/backend-dist/`` by ``build_windows.ps1`` and bundled by
Electron Builder via the ``extraResources`` entry in ``desktop/package.json``
(so the exe ends up at ``resources/backend/rentacar-backend.exe`` at runtime,
matching ``desktop/main.js``).
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

# One-folder build: the EXE excludes binaries which are collected into the
# output directory by COLLECT below. This yields ``dist/rentacar-backend/``.
exe = EXE(
    pyz,
    a.scripts,
    [],
    exclude_binaries=True,
    name="rentacar-backend",
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    console=True,
    disable_windowed_traceback=False,
    target_arch=None,
)

coll = COLLECT(
    exe,
    a.binaries,
    a.zipfiles,
    a.datas,
    strip=False,
    upx=True,
    upx_exclude=[],
    name="rentacar-backend",
)
