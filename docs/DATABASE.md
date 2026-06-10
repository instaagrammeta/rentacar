# Схема базы данных — Rentacar CRM

СУБД по умолчанию — **SQLite**. ORM — SQLAlchemy 2.0. Все таблицы содержат
`id` (PK), `created_at`, `updated_at`.

## Диаграмма связей

```
users ──< audit_logs
users ──< rentals (employee_id)
users ──< payments (cashier_id)

clients ──< reservations >── cars
clients ──< rentals       >── cars
clients ──< payments
clients ──< accidents     >── cars
clients ──1 blacklist_entries

cars ──< car_photos
rentals ──1 vehicle_returns
rentals ──< payments
accidents ──< accident_photos
company_settings (одна строка)
```

## Таблицы

### users — сотрудники/учётные записи
| Поле | Тип | Описание |
|------|-----|----------|
| username | str unique | Логин |
| email | str unique? | Почта |
| full_name | str | ФИО |
| password_hash | str | bcrypt-хэш |
| role | enum | administrator / rental_manager / cashier / operator |
| is_active | bool | Активна ли запись |

### clients — клиенты
| Поле | Тип | Описание |
|------|-----|----------|
| client_code | str unique | Код клиента (CL-000001) |
| first_name, last_name | str | Имя, фамилия |
| phone | str (index) | Телефон |
| email | str | Почта |
| date_of_birth | date | Дата рождения (проверка возраста ≥ 21) |
| passport_number | str | Номер паспорта |
| driver_license_number | str | Номер ВУ |
| driver_license_issue_date | date | Дата выдачи ВУ |
| driver_experience_years | int | Стаж (проверка ≥ 1) |
| passport_scan, driver_license_scan | str | Пути к сканам |
| qr_code_path | str | Путь к QR-коду |
| is_vip | bool | VIP-клиент |
| notes | text | Заметки |
| status | enum | active / pending_verification / blacklisted |

### cars — автомобили
| Поле | Тип | Описание |
|------|-----|----------|
| brand, model | str | Марка, модель |
| year | int | Год |
| color | str | Цвет |
| vin | str unique | VIN |
| plate_number | str unique | Гос. номер |
| mileage | int | Пробег |
| daily_price, weekly_price, monthly_price | float | Тарифы |
| deposit_amount | float | Депозит |
| status | enum | available / reserved / rented / maintenance |

### car_photos
`car_id` (FK → cars), `file_path`.

### reservations — бронирования
`client_id` (FK), `car_id` (FK), `start_date`, `end_date`, `deposit`, `notes`,
`status` (reserved / confirmed / cancelled).

### rentals — договоры аренды
| Поле | Тип | Описание |
|------|-----|----------|
| contract_number | str unique | Номер договора (DOG-2026-00001) |
| client_id, car_id | FK | Клиент, авто |
| reservation_id | FK? | Связанная бронь |
| employee_id | FK? → users | Ответственный |
| rental_start, rental_end | date | Период |
| deposit, daily_price, total_price | float | Финансы |
| start_mileage | int | Пробег при выдаче |
| pdf_path | str | Путь к PDF |
| status | enum | active / completed / cancelled |

### vehicle_returns — возвраты (1:1 с rentals)
| Поле | Тип | Описание |
|------|-----|----------|
| rental_id | FK unique | Аренда |
| return_date | datetime | Дата возврата |
| mileage | int | Пробег при возврате |
| fuel_level | int | Топливо, % |
| damages | text | Повреждения |
| extra_days | int | Доп. дни |
| late_fee | float | Плата за просрочку |
| damage_cost | float | Ущерб |
| penalties | float | Штрафы |
| final_payment | float | Итог (+ доплата / − возврат) |

### payments — платежи
| Поле | Тип | Описание |
|------|-----|----------|
| receipt_number | str unique | Номер квитанции (KV-2026-00001) |
| client_id | FK | Клиент |
| rental_id | FK? | Договор |
| cashier_id | FK? → users | Кассир |
| amount | float | Сумма |
| method | enum | cash / bank_transfer / card |
| payment_type | enum | deposit / rental / penalty / damage / refund |
| paid_at | datetime | Дата оплаты |
| receipt_path | str | Путь к квитанции |

### accidents / accident_photos — ДТП
`car_id` (FK), `client_id` (FK?), `rental_id` (FK?), `accident_date`,
`description`, `repair_cost`; фото в `accident_photos`.

### blacklist_entries — чёрный список (1:1 с clients)
`client_id` (FK unique), `reason` (enum), `comment`, `created_by` (FK → users).

### audit_logs — журнал аудита
`user_id` (FK?), `username`, `action`, `entity`, `entity_id`, `ip_address`,
`details`.

### company_settings — настройки (одна строка)
`company_name`, `address`, `phone`, `email`, `logo_path`, `currency`,
`contract_terms`.

## Перечисления (enums)

| Enum | Значения |
|------|----------|
| UserRole | administrator, rental_manager, cashier, operator |
| ClientStatus | active, pending_verification, blacklisted |
| CarStatus | available, reserved, rented, maintenance |
| ReservationStatus | reserved, confirmed, cancelled |
| RentalStatus | active, completed, cancelled |
| PaymentMethod | cash, bank_transfer, card |
| PaymentType | deposit, rental, penalty, damage, refund |
| BlacklistReason | fraud, vehicle_damage, non_payment, serious_violation |
