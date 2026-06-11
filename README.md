# Wheelzie — Rentacar CRM (v2, public online edition)

A car-rental management system (CRM) for the public web. Version 2 is a clean,
containerised rewrite:

- **Backend** — Go (Gin) + PostgreSQL + Redis
- **Frontend** — Nuxt 4 + Vue 3 + TypeScript + TailwindCSS + Pinia
- **Infra** — Docker Compose + GitHub Actions CI/CD

> The legacy Python/Flask backend, the old Vue/Vite frontend and the Electron
> desktop wrapper have been removed on this branch. Only the clean v2 code
> remains.

## Run the whole project with one command

```bash
docker compose up --build
```

This starts four services:

| Service    | URL / Port                  | Notes                              |
|------------|-----------------------------|------------------------------------|
| frontend   | http://localhost:3000       | Nuxt 4 SPA (Wheelzie UI)           |
| backend    | http://localhost:5000/api   | Go + Gin REST API                  |
| postgres   | localhost:5432              | data store                         |
| redis      | localhost:6379              | cache + JWT logout blacklist       |

Default login: **admin / admin123** (change it after the first sign-in).

On first launch the database is migrated and seeded automatically (admin, demo
users, company settings, demo cars & clients).

## Configuration

Copy `.env.example` files if you want to override defaults. The most important
variables when deploying to a real server:

```env
# docker-compose (root) reads these
POSTGRES_PASSWORD=strong-password
RENTACAR_JWT_SECRET=long-random-secret
RENTACAR_SECRET_KEY=another-random-secret

# URL of the API as seen from the user's browser (set to your domain)
NUXT_PUBLIC_API_BASE=https://api.example.com/api
RENTACAR_CORS_ORIGINS=https://app.example.com
```

## Project structure

```
.
├── backend-go/        # Go (Gin + GORM + Redis) backend
├── frontend/          # Nuxt 4 + Vue 3 + Tailwind + Pinia frontend
├── docker-compose.yml # full stack (postgres + redis + backend + frontend)
└── .github/workflows/ # CI/CD for backend and frontend
```

See [`backend-go/README.md`](backend-go/README.md) for backend details.

## CI/CD

GitHub Actions builds, lints/tests and produces Docker images for both apps on
every push to `v-2` / `main`:

- `backend-go.yml` — gofmt, vet, build, test, Docker image → GHCR
- `frontend.yml` — install, Nuxt build, Docker image → GHCR

## Deploying on a server

1. Install Docker + Docker Compose on the server.
2. `git clone -b v-2 https://github.com/instaagrammeta/rentacar.git`
3. Set the environment variables above (a root `.env` file works).
4. `docker compose up --build -d`
5. Put Nginx + HTTPS (certbot) in front of ports `3000` (frontend) and
   `5000` (backend) for your domain.

Default credentials should be changed immediately in production.
