# =============================================================================
# Rentacar CRM - Windows build script
# -----------------------------------------------------------------------------
# Produces a one-click Windows installer (.exe) that runs fully offline.
#
# Prerequisites (one-time, requires internet):
#   * Python 3.13+      (https://www.python.org)
#   * Node.js 20+       (https://nodejs.org)
#
# Usage (from a PowerShell prompt in the project root):
#   ./build_windows.ps1
#
# Output:
#   desktop/release/Rentacar CRM Setup <version>.exe
# =============================================================================

$ErrorActionPreference = "Stop"
$root = $PSScriptRoot

Write-Host "==> [1/5] Building Python backend (PyInstaller)" -ForegroundColor Green
Push-Location "$root/backend"
python -m venv .venv
& ".venv/Scripts/python.exe" -m pip install --upgrade pip
& ".venv/Scripts/python.exe" -m pip install -r requirements.txt
& ".venv/Scripts/python.exe" -m pip install pyinstaller
& ".venv/Scripts/python.exe" -m PyInstaller rentacar-backend.spec --noconfirm
Pop-Location

Write-Host "==> [2/5] Building Vue frontend (Vite)" -ForegroundColor Green
Push-Location "$root/frontend"
npm install
npm run build
Pop-Location

Write-Host "==> [3/5] Staging artifacts for Electron" -ForegroundColor Green
$desktop = "$root/desktop"
Remove-Item -Recurse -Force "$desktop/backend-dist" -ErrorAction SilentlyContinue
Remove-Item -Recurse -Force "$desktop/frontend-dist" -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path "$desktop/backend-dist" | Out-Null
Copy-Item -Recurse -Force "$root/backend/dist/rentacar-backend/*" "$desktop/backend-dist/"
Copy-Item -Recurse -Force "$root/frontend/dist/*" "$desktop/frontend-dist/"

Write-Host "==> [4/5] Installing Electron build tooling" -ForegroundColor Green
Push-Location $desktop
npm install

Write-Host "==> [5/5] Building Windows installer (NSIS)" -ForegroundColor Green
npm run dist
Pop-Location

Write-Host "Done. Installer is in desktop/release/." -ForegroundColor Cyan
