// Electron main process for Rentacar CRM.
//
// Responsibilities:
//  1. Spawn the bundled Python/Flask backend (packaged with PyInstaller into
//     `resources/backend/rentacar-backend.exe`, or `python run.py --serve`
//     during development).
//  2. Wait until the backend health endpoint responds.
//  3. Create the application window and load the built Vue frontend.
//  4. Cleanly terminate the backend when the app quits.
//  5. Provide a native "Файл" menu (Сохранить/Открыть проект) wired to the API.

const { app, BrowserWindow, Menu, dialog } = require('electron')
const { spawn } = require('child_process')
const path = require('path')
const http = require('http')
const fs = require('fs')

const BACKEND_PORT = 5000
const BACKEND_URL = `http://127.0.0.1:${BACKEND_PORT}`
const isDev = !app.isPackaged

let backendProcess = null
let mainWindow = null

function resourcePath(...segments) {
  // In production, extraResources land in process.resourcesPath.
  const base = isDev ? path.join(__dirname, '..') : process.resourcesPath
  return path.join(base, ...segments)
}

function startBackend() {
  const env = {
    ...process.env,
    RENTACAR_ENV: 'production',
    // Store the database & uploads under the user's app data directory.
    RENTACAR_DATA_DIR: path.join(app.getPath('userData'), 'data'),
    RENTACAR_BACKUP_DIR: path.join(app.getPath('userData'), 'backups'),
  }

  if (isDev) {
    // Development: run the Flask app from source.
    backendProcess = spawn('python', ['run.py', '--serve', '--port', String(BACKEND_PORT)], {
      cwd: path.join(__dirname, '..', 'backend'),
      env,
      shell: true,
    })
  } else {
    // Production: run the PyInstaller-built backend executable.
    const exe = resourcePath('backend', process.platform === 'win32' ? 'rentacar-backend.exe' : 'rentacar-backend')
    backendProcess = spawn(exe, ['--serve', '--port', String(BACKEND_PORT)], { env })
  }

  backendProcess.stdout?.on('data', (d) => console.log(`[backend] ${d}`))
  backendProcess.stderr?.on('data', (d) => console.error(`[backend] ${d}`))
  backendProcess.on('exit', (code) => console.log(`Backend exited with code ${code}`))
}

function waitForBackend(retries = 60) {
  return new Promise((resolve, reject) => {
    const attempt = (left) => {
      http
        .get(`${BACKEND_URL}/api/health`, (res) => {
          if (res.statusCode === 200) resolve()
          else retry(left)
        })
        .on('error', () => retry(left))
    }
    const retry = (left) => {
      if (left <= 0) return reject(new Error('Backend did not start in time'))
      setTimeout(() => attempt(left - 1), 500)
    }
    attempt(retries)
  })
}

function buildMenu() {
  const template = [
    {
      label: 'Файл',
      submenu: [
        {
          label: 'Сохранить проект',
          click: () => mainWindow?.webContents.executeJavaScript('window.dispatchEvent(new Event("rentacar:save-project"))'),
        },
        {
          label: 'Открыть проект',
          click: () => mainWindow?.webContents.executeJavaScript('window.dispatchEvent(new Event("rentacar:open-project"))'),
        },
        { type: 'separator' },
        { role: 'quit', label: 'Выход' },
      ],
    },
    {
      label: 'Правка',
      submenu: [
        { role: 'undo', label: 'Отменить' },
        { role: 'redo', label: 'Повторить' },
        { type: 'separator' },
        { role: 'cut', label: 'Вырезать' },
        { role: 'copy', label: 'Копировать' },
        { role: 'paste', label: 'Вставить' },
      ],
    },
    {
      label: 'Вид',
      submenu: [
        { role: 'reload', label: 'Обновить' },
        { role: 'togglefullscreen', label: 'Полный экран' },
        { role: 'zoomIn', label: 'Увеличить' },
        { role: 'zoomOut', label: 'Уменьшить' },
      ],
    },
  ]
  Menu.setApplicationMenu(Menu.buildFromTemplate(template))
}

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1440,
    height: 900,
    minWidth: 1024,
    minHeight: 700,
    title: 'Rentacar CRM',
    icon: resourcePath('desktop', 'resources', 'icon.ico'),
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
    },
  })

  buildMenu()

  if (isDev) {
    mainWindow.loadURL('http://127.0.0.1:5173')
  } else {
    const indexHtml = path.join(__dirname, 'frontend-dist', 'index.html')
    if (fs.existsSync(indexHtml)) {
      mainWindow.loadFile(indexHtml)
    } else {
      mainWindow.loadURL(BACKEND_URL)
    }
  }
}

app.whenReady().then(async () => {
  startBackend()
  try {
    await waitForBackend()
  } catch (e) {
    dialog.showErrorBox('Ошибка запуска', 'Не удалось запустить сервер приложения.')
  }
  createWindow()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit()
})

app.on('quit', () => {
  if (backendProcess) {
    try {
      backendProcess.kill()
    } catch (e) {
      /* ignore */
    }
  }
})
