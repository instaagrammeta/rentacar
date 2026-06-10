# Документация для разработчиков — Rentacar CRM

## Архитектура

Проект следует принципам **чистой архитектуры** с разделением ответственности:

```
HTTP (Vue SPA)
   │  REST/JSON + JWT
   ▼
API blueprints  (app/api)        — валидация запроса, сериализация ответа
   ▼
Service Layer   (app/services)   — бизнес-логика, правила, расчёты, транзакции
   ▼
Repository      (app/repositories) — доступ к данным (Repository Pattern)
   ▼
Models          (app/models)     — SQLAlchemy ORM
   ▼
SQLite / SQLAlchemy
```

### Принципы

- **API не содержит бизнес-логики** — только разбор запроса и вызов сервиса.
- **Сервисы не обращаются к ORM-сессии напрямую** — только через репозитории.
- **Репозитории** инкапсулируют запросы; общий CRUD в `BaseRepository`.
- **Типизация**: используются type hints во всём backend.
- **Обработка ошибок**: доменные исключения `app/utils/errors.py`
  (`ValidationError`, `NotFoundError`, `BusinessRuleError`, …) транслируются в
  единообразные JSON-ответы с русскими сообщениями.
- **Логирование**: ротация логов в `data/logs/rentacar.log` + журнал аудита в БД.

## Backend: ключевые модули

| Путь | Назначение |
|------|------------|
| `app/config.py` | Конфигурация (dev/test/prod), пути к данным, бизнес-константы |
| `app/extensions.py` | Синглтоны расширений (db, jwt, migrate, cors) |
| `app/__init__.py` | Application factory `create_app()` |
| `app/models/` | Доменные модели и перечисления (`enums.py`) |
| `app/repositories/` | Репозитории + синглтоны |
| `app/services/` | Бизнес-логика по модулям |
| `app/api/` | REST blueprints |
| `app/utils/security.py` | Хэширование паролей (bcrypt), RBAC-декораторы |
| `app/utils/pricing.py` | Расчёт стоимости аренды |
| `app/services/pdf_service.py` | Генерация PDF (ReportLab) |
| `app/services/report_service.py` | Excel-отчёты (OpenPyXL) |
| `app/services/project_service.py` | Формат `.rentacar` (экспорт/импорт) |
| `app/scheduler.py` | Ежедневные авто-бэкапы (APScheduler) |

### Роли и права

Матрица `ROLE_PERMISSIONS` в `app/utils/security.py` задаёт доступ к модулям.
Эндпоинты защищаются декораторами `@permission_required("clients")` или
`@role_required(UserRole.ADMINISTRATOR)`. Фронтенд дублирует матрицу в
`stores/auth.js`, чтобы скрывать недоступные пункты меню.

## Frontend: структура

| Путь | Назначение |
|------|------------|
| `src/main.js` | Инициализация Vue, Pinia, Router, Vuetify, i18n |
| `src/plugins/vuetify.js` | Тема Green + White |
| `src/i18n/ru.js` | Все русские строки интерфейса |
| `src/api/` | Axios-клиент с JWT и обёртки эндпоинтов |
| `src/stores/` | Pinia: `auth` (токен/роль), `app` (уведомления) |
| `src/router/` | Маршруты + guard'ы (auth + права) |
| `src/layouts/MainLayout.vue` | Каркас: боковое меню, меню «Файл» |
| `src/views/` | Экраны 12 модулей + вход/пользователи/аудит |

## Миграции

```bash
cd backend
export FLASK_APP=run.py
flask db migrate -m "change"
flask db upgrade
```

При запуске `create_app()` также вызывает `db.create_all()` — это удобно для
десктопной сборки, где миграции не запускаются.

## Тестирование вручную

```bash
# health
curl http://127.0.0.1:5000/api/health
# login
curl -X POST http://127.0.0.1:5000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## Стиль кода

- Python: PEP 8, type hints, docstrings. Рекомендуются `ruff`, `black`, `mypy`.
- JS/Vue: Composition API (`<script setup>`), ESLint/Prettier по желанию.

## Сборка десктопа

См. `build_windows.ps1` и `desktop/`. Backend упаковывается PyInstaller'ом
(`backend/rentacar-backend.spec`), frontend — Vite, всё объединяется
Electron Builder (NSIS). Альтернатива — `desktop/desktop_app.py` (PyWebView).
