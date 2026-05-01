# kilas-admin

Simple admin dashboard/CMS for Kilas, built with Go (Gin + GORM) and React (Vite).

## Features

- Admin authentication (`users` table in `kilas_admin_db`)
- Sidebar-based CMS UI (React)
- Full CRUD for Kilas entities:
  - `users`
  - `transactions`
  - `products`
  - `decks`
  - `cards`
  - `ai_generation_history`
  - `issues`
- Dashboard summary counters for all major entities
- Public endpoint to submit issue reports

## Project Structure

- `backend`: Go API server
- `frontend`: React admin dashboard

## Database

Backend uses two databases:

- `kilas_admin_db` (admin DB):
  - `users` (admin users that can login)
  - `issues` (user-reported transaction issues)

- `kilas_db` (main Kilas app DB):
  - `users`, `transactions`, `products`, `decks`, `cards`, `ai_generation_histories`
  - used by `/api/admin/:entity` CRUD endpoints

## Backend Setup

1. Copy env template:
   - `cp .env.example .env`
2. Update DB and secrets in `.env`
   - `DATABASE_NAME` must be `kilas_admin_db`
   - `KILAS_DATABASE_NAME` should point to your main app DB (default `kilas_db`)
3. Ensure MySQL database exists:
   - `CREATE DATABASE kilas_admin_db;`
4. Run backend:
   - `go mod tidy`
   - `go run .`

Default backend port: `8081`

### Seed Admin User

Set these in `.env` before first run:

- `ADMIN_SEED_NAME`
- `ADMIN_SEED_EMAIL`
- `ADMIN_SEED_PASSWORD`

Backend will auto-create the admin if email does not exist yet.

## Frontend Setup

1. Copy env template:
   - `cp .env.example .env`
2. Install deps and run:
   - `npm install`
   - `npm run dev`

Default frontend API target: `http://localhost:8081/api`

## API Endpoints

- `POST /api/auth/login`
- `GET /api/auth/me`
- `POST /api/issues/report` (public user report endpoint)

Protected admin endpoints:

- `GET /api/admin/summary`
- `GET /api/admin/:entity`
- `GET /api/admin/:entity/:id`
- `POST /api/admin/:entity`
- `PUT /api/admin/:entity/:id`
- `DELETE /api/admin/:entity/:id`

Supported `:entity` values:

- `users`
- `transactions`
- `products`
- `decks`
- `cards`
- `ai_generation_history`
- `issues`
