# Руководство по установке — Rentacar CRM

## Вариант A. Готовый установщик (для конечных пользователей)

1. Получите файл `Rentacar CRM Setup <версия>.exe`.
2. Дважды кликните по нему — установка пройдёт в один клик.
3. На рабочем столе появится ярлык **Rentacar CRM**.
4. Запустите приложение. Интернет не требуется — всё работает локально.
5. Войдите под учётной записью администратора: **admin / admin123**.

Данные (база SQLite, загруженные файлы, резервные копии) хранятся в профиле
пользователя: `%APPDATA%/Rentacar CRM/`.

## Вариант B. Сборка установщика из исходного кода

### Требования (нужен интернет только на этапе сборки)

- Python 3.13+ — https://www.python.org (включите «Add to PATH»)
- Node.js 20+ — https://nodejs.org

### Сборка

В PowerShell из корня проекта:

```powershell
./build_windows.ps1
```

Скрипт выполнит:

1. Сборку backend в автономный `.exe` (PyInstaller).
2. Сборку frontend (Vite) в статические файлы.
3. Подготовку артефактов и сборку установщика Electron (NSIS).

Готовый установщик: `desktop/release/`.

## Вариант C. Запуск из исходного кода (разработка)

### Backend

```bash
cd backend
python -m venv .venv
# Windows:
.venv\Scripts\activate
# Linux/macOS:
source .venv/bin/activate

pip install -r requirements.txt
python seed.py        # инициализация БД + демо-данные
python run.py         # сервер на http://127.0.0.1:5000
```

### Frontend

```bash
cd frontend
npm install
npm run dev           # http://127.0.0.1:5173
```

### Десктоп через PyWebView (без Node)

```bash
cd frontend && npm install && npm run build && cd ..
pip install pywebview waitress
python desktop/desktop_app.py
```

## Переменные окружения (необязательно)

| Переменная | Назначение | По умолчанию |
|------------|------------|--------------|
| `RENTACAR_ENV` | `development` / `production` / `testing` | `development` |
| `RENTACAR_SECRET_KEY` | Секретный ключ Flask | `change-me-in-production` |
| `RENTACAR_JWT_SECRET` | Секрет для JWT | = SECRET_KEY |
| `RENTACAR_DATABASE_URI` | Строка подключения к БД | локальный SQLite |
| `RENTACAR_DATA_DIR` | Каталог данных | `backend/data` |
| `RENTACAR_BACKUP_DIR` | Каталог резервных копий | `backups/` |
| `RENTACAR_ENABLE_BACKUPS` | Автобэкапы (`true`/`false`) | `true` |

## Решение проблем

- **Порт 5000 занят** — закройте конкурирующее приложение или измените порт
  (`python run.py --port 5050`).
- **Кириллица в PDF выглядит некорректно** — установите шрифт DejaVu Sans или
  убедитесь, что в системе доступен Arial; backend подхватит его автоматически.
