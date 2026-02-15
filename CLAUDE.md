# Fambudg

Family budget tracking web app. Monorepo: Go backend + React frontend (later).

## Current Status

Phases 1–4 complete. Backend fully functional with auth, transactions, budgets, reports, saving goals, bill reminders, CSV import/export, recurring transactions, role-based access control, allowances, and family spending comparison. No frontend yet.

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
│   ├── cmd/
│   │   ├── server/main.go       # entry point, route setup
│   │   └── migrate/main.go      # migration runner
│   ├── internal/
│   │   ├── config/              # env-based configuration
│   │   ├── database/            # pgx pool setup
│   │   ├── model/               # domain structs + request/response types
│   │   ├── repository/          # database access (SQL queries)
│   │   ├── service/             # business logic
│   │   ├── handler/             # HTTP handlers
│   │   └── middleware/          # auth, RBAC, logging, CORS
│   ├── migrations/              # goose SQL migration files (001–008)
│   ├── tests/features/          # godog BDD .feature files (13 files)
│   ├── tests/steps/             # godog step definitions (15 files)
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

## Role-Based Access Control

Three roles: `admin`, `member`, `child`. JWT claims include `user_id` and `role`.

Routes are organized into three middleware groups in `main.go`:

| Group | Roles | Access |
|-------|-------|--------|
| All authenticated | admin, member, child | accounts, transactions, categories (read), reports, search, allowances (list) |
| Admin + Member | admin, member | categories (create), budgets (read), saving goals (read), bill reminders (read + pay), transfers, recurring, import/export |
| Admin only | admin | user management, categories (update/delete), budgets (write), saving goals (write), bill reminders (write), by-member report, allowances (manage) |

Admin sees all family data in List/Report/Search endpoints. Member/child see only own data. Ownership checks on Get/Update/Delete prevent cross-user access.

## API Endpoints

### Public
- `POST /api/auth/register` — first user becomes admin
- `POST /api/auth/login`

### All Authenticated
- `GET /api/auth/me`
- `/api/accounts` — CRUD (admin sees all, others own)
- `GET /api/categories` — list
- `/api/transactions` — CRUD + `?isShared=true/false` filter (admin sees all, others own)
- `/api/reports/dashboard` — accounts + month summary + recent transactions
- `/api/reports/monthly?month=N&year=N`
- `/api/reports/by-category?month=N&year=N`
- `/api/reports/trends?months=N`
- `GET /api/search?description=&startDate=&endDate=&categoryId=&accountId=&minAmount=&maxAmount=&tags=`
- `GET /api/allowances`

### Admin + Member
- `POST /api/categories`
- `GET /api/budgets`, `GET /api/budgets/summary`
- `GET /api/saving-goals`
- `GET /api/bill-reminders`, `GET /api/bill-reminders/upcoming`, `POST /api/bill-reminders/{id}/pay`
- `POST /api/transfers`
- `POST /api/transactions/generate-recurring?upTo=YYYY-MM-DD`
- `POST /api/import/csv`, `GET /api/export/csv`

### Admin Only
- `/api/users` — CRUD (create users with role)
- `PUT/DELETE /api/categories/{id}`
- `/api/budgets` — create/update/delete
- `/api/saving-goals` — create/update/contribute
- `/api/bill-reminders` — create/update/delete
- `GET /api/reports/by-member?month=N&year=N`
- `POST/PUT /api/allowances`

## Database Tables

8 migrations: `users`, `accounts`, `categories`, `transactions`, `budgets`, `saving_goals`, `bill_reminders`, `allowances`

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

# BDD tests (47 scenarios)
go test ./tests/... -v

# All tests
go test ./...
```

## Environment Variables

See `.env.example` in project root for required variables. Actual `.env` lives in `backend/` and is never committed.

## Full Plan

See `.claude/plan.md` for the complete project plan including data models, all API endpoints, and phased build checklist. Phases 1–4 complete, Phase 5 (Cloud Deployment) is next.
