# Fambudg

Family budget tracking web app. Monorepo: Go backend + React frontend (later).

## Tech Stack

- **Backend:** Go (in `backend/`)
- **Frontend:** React (in `frontend/` — later)
- **Database:** PostgreSQL (via Docker locally, managed cloud in production)
- **Router:** chi/v5
- **DB driver:** pgx/v5
- **Migrations:** goose/v3
- **Auth:** golang-jwt/jwt/v5 (JWT with bcrypt passwords)
- **Validation:** go-playground/validator/v10
- **Testing:** Go standard testing + godog (BDD, Gherkin syntax)

## Project Structure (Monorepo)

```
fambudg/
├── backend/
│   ├── cmd/server/main.go       # entry point
│   ├── internal/
│   │   ├── config/              # env-based configuration
│   │   ├── model/               # domain structs
│   │   ├── repository/          # database access (SQL queries)
│   │   ├── service/             # business logic
│   │   ├── handler/             # HTTP handlers
│   │   └── middleware/          # auth, logging, CORS
│   ├── migrations/              # goose SQL migration files
│   ├── tests/features/          # godog BDD .feature files
│   ├── tests/steps/             # godog step definitions
│   ├── go.mod
│   └── .env
├── frontend/                    # React app (later)
├── .claude/plan.md              # full project plan with phases
├── docker-compose.yml           # Postgres for local dev
├── .env.example
└── .gitignore
```

## Architecture Pattern

Request flow: Router → Middleware → Handler → Service → Repository → Database

- **Handler:** parses HTTP request, validates input, calls service, returns JSON response
- **Service:** business logic, orchestrates repositories, no HTTP awareness
- **Repository:** raw SQL queries via pgx, returns domain models

## Key Conventions

- Money stored as `BIGINT` in cents (never floats). $19.99 = 1999
- UUID primary keys everywhere
- Positive amounts = income, negative = expense
- All timestamps as `TIMESTAMPTZ` in UTC
- Table names: snake_case, plural (e.g., `transactions`, `saving_goals`)
- Go files: snake_case (e.g., `transaction_repo.go`)
- JSON responses: camelCase keys (e.g., `categoryId`, `createdAt`)
- Errors returned as `{"error": "message"}` with appropriate HTTP status codes
- All endpoints prefixed with `/api/`

## Running Locally

```bash
# Start Postgres (from project root)
docker compose up -d

# Run migrations (from backend/)
cd backend
go run cmd/migrate/main.go up

# Start server (from backend/)
go run cmd/server/main.go
# Runs on http://localhost:8080
```

## Running Tests

```bash
# From backend/ directory

# Unit tests
go test ./internal/...

# BDD tests
go test ./tests/... --godog.format=pretty

# All tests
go test ./...
```

## Environment Variables

See `.env.example` in project root for required variables. Actual `.env` lives in `backend/` and is never committed.

## Full Plan

See `.claude/plan.md` for the complete project plan including data models, all API endpoints, and phased build checklist.