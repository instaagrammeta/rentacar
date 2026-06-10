# Документация API — Rentacar CRM

Базовый URL: `http://127.0.0.1:5000/api`

Все защищённые эндпоинты требуют заголовок:

```
Authorization: Bearer <access_token>
```

Ошибки возвращаются в формате `{"error": "сообщение"}` с соответствующим HTTP-кодом
(400 — валидация, 401 — авторизация, 403 — права, 404 — не найдено, 409 —
конфликт, 422 — нарушение бизнес-правила).

## Аутентификация

| Метод | Путь | Описание | Доступ |
|-------|------|----------|--------|
| POST | `/auth/login` | Вход. Тело: `{username, password}` → `{access_token, refresh_token, user}` | публично |
| GET | `/auth/me` | Текущий пользователь | любой |
| GET | `/auth/roles` | Список ролей | любой |
| GET | `/auth/users` | Список пользователей | admin |
| POST | `/auth/users` | Создать пользователя | admin |
| POST | `/auth/users/{id}/password` | Сменить пароль | admin |

## Клиенты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/clients?search=&status=&page=&per_page=` | Список (пагинация/поиск/фильтр) |
| GET | `/clients/{id}` | Карточка клиента |
| GET | `/clients/{id}/history` | История (аренды, платежи, штрафы, ДТП) |
| GET | `/clients/search?phone=&code=&term=` | Поиск VIP/постоянного клиента |
| POST | `/clients` | Создать |
| PUT | `/clients/{id}` | Обновить |
| DELETE | `/clients/{id}` | Удалить |

## Автомобили

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/cars?search=&status=&page=&per_page=` | Список |
| GET | `/cars/available` | Доступные авто |
| GET | `/cars/{id}` | Карточка |
| POST | `/cars` | Создать |
| PUT | `/cars/{id}` | Обновить |
| DELETE | `/cars/{id}` | Удалить |

## Бронирования

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/reservations?status=` | Список |
| POST | `/reservations` | Создать |
| POST | `/reservations/{id}/status` | Сменить статус `{status}` |
| POST | `/reservations/{id}/cancel` | Отменить |

## Аренда и возврат

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/rentals?status=` | Список договоров |
| GET | `/rentals/{id}` | Договор (+ возврат) |
| POST | `/rentals` | Создать договор (расчёт + PDF) |
| POST | `/rentals/{id}/cancel` | Отменить |
| GET | `/rentals/{id}/contract` | Скачать PDF договора |
| POST | `/rentals/{id}/return/preview` | Предрасчёт возврата |
| POST | `/rentals/{id}/return` | Оформить возврат |

## Платежи

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/payments?client_id=&rental_id=` | Список |
| GET | `/payments/{id}` | Платёж |
| POST | `/payments` | Создать (+ квитанция PDF) |
| GET | `/payments/{id}/receipt` | Скачать квитанцию PDF |

## Чёрный список / ДТП

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/blacklist` | Список |
| GET | `/blacklist/reasons` | Причины |
| POST | `/blacklist` | Добавить `{client_id, reason, comment}` |
| DELETE | `/blacklist/{client_id}` | Убрать |
| GET | `/accidents?car_id=` | Список ДТП |
| POST | `/accidents` | Создать |
| PUT | `/accidents/{id}` | Обновить |
| DELETE | `/accidents/{id}` | Удалить |

## Дашборд и отчёты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/dashboard/summary` | Сводные показатели |
| GET | `/dashboard/revenue-by-month?months=12` | Выручка по месяцам |
| GET | `/dashboard/top-cars?limit=5` | Топ авто |
| GET | `/dashboard/rental-statistics` | Статистика аренд |
| GET | `/reports/daily-revenue?date=` | Дневная выручка |
| GET | `/reports/monthly-revenue?year=&month=` | Месячная |
| GET | `/reports/yearly-revenue?year=` | Годовая |
| GET | `/reports/profitable-cars` | Прибыльные авто |
| GET | `/reports/active-rentals` | Активные аренды |
| GET | `/reports/debtors` | Должники |
| GET | `/reports/client-statistics` | Статистика клиентов |
| GET | `/reports/export/{type}` | Экспорт в Excel |

## Настройки, бэкапы, проект, файлы, аудит

| Метод | Путь | Описание | Доступ |
|-------|------|----------|--------|
| GET | `/settings` | Данные компании | любой |
| PUT | `/settings` | Обновить | admin |
| GET | `/backups` | Список копий | admin |
| POST | `/backups` | Создать копию | admin |
| POST | `/project/export` | Экспорт `.rentacar` | admin |
| POST | `/project/import` | Импорт `.rentacar` (multipart `file`) | admin |
| POST | `/uploads` | Загрузка файлов (multipart `files`) | любой |
| GET | `/uploads/{path}` | Отдача файла | публично |
| GET | `/audit?limit=&offset=` | Журнал аудита | admin |
| GET | `/health` | Проверка сервиса | публично |
