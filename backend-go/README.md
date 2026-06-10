# Rentacar CRM — Public Backend (Go)

A full rewrite of the Rentacar CRM backend in **Go**, built for online/public
deployment. It keeps **all** functionality of the original Python/Flask backend
(clients, cars, reservations, rental contracts, returns, payments, blacklist,
accidents, dashboard, reports, settings, audit, backups, project import/export,
PDF contracts & receipts, QR codes, Excel exports) while using a modern,
production-ready stack.

## Tech stack

| Concern         | Technology                          |
|-----------------|-------------------------------------|
| Language        | Go 1.23                             |
| HTTP framework  | [Gin](https://github.com/gin-gonic/gin) |
| Database        | PostgreSQL (via [GORM](https://gorm.io)) |
| Cache / tokens  | Redis (caching + JWT logout blacklist) |
| Auth            | JWT (HS256) + bcrypt                 |
| PDF             | go-pdf/fpdf (embedded Cyrillic font) |
| Excel           | excelize                            |
| QR codes        | go-qrcode                           |
| Scheduler       | robfig/cron (daily backups)         |
| Containerisation| Docker + docker-compose             |
| CI/CD           | GitHub Actions                      |

## Quick start (Docker — recommended)

Everything (PostgreSQL + Redis + backend) starts with a single command:

```bash
cd backend-go
docker compose up --build
```

- API: `http://localhost:5000/api`
- Health check: `http://localhost:5000/api/health`
- Default admin login: **admin / admin123** (change it after first login)

On first launch the database is migrated automatically and seeded with the
default admin, demo users (`manager`, `cashier`, `operator` — password is the
username + `123`), company settings and a few demo cars/clients.

## Running locally (without Docker)

Requires Go 1.23+, a PostgreSQL instance and (optionally) Redis.

```bash
cd backend-go
cp .env.example .env        # adjust connection strings
go mod tidy                 # resolves dependencies and creates go.sum
go run ./cmd/server         # start the API
go run ./cmd/seed           # (optional) seed demo data
```

> **Note on `go.sum`:** this code was authored in an offline environment, so
> `go.sum` is generated on first build (`go mod tidy`). The Docker build and CI
> pipeline run `go mod tidy` automatically.

## Configuration

All settings come from environment variables (see `.env.example`):

| Variable                    | Default                          | Description                       |
|-----------------------------|----------------------------------|-----------------------------------|
| `RENTACAR_ENV`              | `development`                    | `development` / `production`      |
| `RENTACAR_PORT`             | `5000`                           | HTTP port                         |
| `RENTACAR_DATABASE_URI`     | `postgres://…`                   | PostgreSQL DSN                    |
| `RENTACAR_REDIS_ADDR`       | `localhost:6379`                 | Redis address                     |
| `RENTACAR_REDIS_ENABLED`    | `true`                           | Toggle Redis usage                |
| `RENTACAR_JWT_SECRET`       | `RENTACAR_SECRET_KEY`            | JWT signing secret                |
| `RENTACAR_JWT_HOURS`        | `12`                             | Access-token lifetime (hours)     |
| `RENTACAR_AUTO_SEED`        | `true`                           | Seed admin + settings on boot     |
| `RENTACAR_SEED_DEMO`        | `true`                           | Also seed demo data               |
| `RENTACAR_ENABLE_BACKUPS`   | `true`                           | Daily 02:00 backups               |
| `RENTACAR_BACKUP_RETENTION` | `30`                             | Backup retention in days          |
| `RENTACAR_CORS_ORIGINS`     | `*`                              | Comma-separated origins or `*`    |

## API overview

All routes are prefixed with `/api`. Authentication uses a Bearer access token
obtained from `POST /api/auth/login`.

- `auth`: `login`, `me`, `logout`, `roles`, `users` (admin)
- `clients`, `cars`, `reservations`, `rentals` (+ contract PDF, returns),
  `payments` (+ receipt PDF), `blacklist`, `accidents`
- `dashboard`: `summary`, `revenue-by-month`, `top-cars`, `rental-statistics`
- `reports`: revenue (daily/monthly/yearly), profitable cars, active rentals,
  debtors, client statistics, and styled Excel `export/{type}`
- `settings`, `backups`, `audit` (admin), `project/export`, `project/import`,
  `uploads`

Role-based access control mirrors the original system (administrator,
rental_manager, cashier, operator).

## CI/CD

`.github/workflows/backend-go.yml` runs on pushes/PRs touching `backend-go/`:

1. **Build & Test** — `go mod tidy`, `gofmt` check, `go vet`, `go build`, `go test`.
2. **Docker Build & Push** — builds the image and pushes it to the GitHub
   Container Registry (`ghcr.io/<owner>/<repo>-backend`) on push events.
