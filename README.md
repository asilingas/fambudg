# Fambudg

Family budget tracking web app with role-based access for parents and children. Track income, expenses, budgets, saving goals, and bills across your whole family.

## Features

- **Transactions** — Track income and expenses with categories, tags, and shared/personal flags
- **Accounts** — Manage checking, savings, credit, and cash accounts with automatic balance updates
- **Budgets** — Set monthly spending limits per category and track actual vs. budgeted
- **Reports** — Dashboard, monthly summaries, category breakdowns, trends, and family spending comparison
- **Saving Goals** — Set targets with deadlines and contribute over time
- **Bill Reminders** — Track recurring bills, mark as paid, auto-advance due dates
- **Transfers** — Move money between accounts
- **Recurring Transactions** — Daily, weekly, monthly, or yearly auto-generation
- **CSV Import/Export** — Bulk import transactions or export for external use
- **Allowances** — Set spending limits for children with automatic tracking
- **Search** — Find transactions by description, date range, amount, category, account, or tags
- **Dark Mode** — Light and dark themes with toggle
- **Internationalization** — English and Lithuanian (EN/LT toggle)
- **Role-Based Access** — Admin, member, and child roles with granular permissions

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go, chi/v5, pgx/v5 |
| Frontend | React 19, TypeScript, Vite, Tailwind CSS 4 |
| UI | Radix UI, shadcn/ui, Recharts, Lucide icons |
| Database | PostgreSQL 17 |
| Auth | JWT (bcrypt passwords) |
| Testing | Vitest + Testing Library (frontend), godog BDD (backend) |

## Prerequisites

- [Go](https://go.dev/) 1.23+
- [Node.js](https://nodejs.org/) 20+
- [Docker](https://www.docker.com/) (for PostgreSQL)

## Getting Started

### 1. Start the database

```bash
docker compose up -d
```

### 2. Set up the backend

```bash
cd backend
cp ../.env.example .env  # edit with your values

# Run migrations
go run cmd/migrate/main.go up

# Seed with sample data (optional)
go run cmd/seed/main.go

# Start the server
go run cmd/server/main.go
# Runs on http://localhost:8080
```

### 3. Start the frontend

```bash
cd frontend
npm install
npm run dev
# Runs on http://localhost:5173
```

### 4. Log in

If you ran the seed command, use:
- **Admin:** admin@family.com / password123
- **Member:** parent@family.com / password123
- **Child:** child@family.com / password123

Otherwise, register at the login page — the first user becomes admin.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `budget` |
| `DB_PASSWORD` | Database password | — |
| `DB_NAME` | Database name | `fambudg` |
| `JWT_SECRET` | Secret for signing JWT tokens | — |
| `SERVER_PORT` | Backend server port | `8080` |

## Project Structure

```
fambudg/
├── backend/
│   ├── cmd/server/          # Entry point, route setup
│   ├── cmd/migrate/         # Migration runner
│   ├── cmd/seed/            # Database seeder
│   ├── internal/
│   │   ├── handler/         # HTTP handlers
│   │   ├── service/         # Business logic
│   │   ├── repository/      # Database queries (pgx)
│   │   ├── model/           # Domain structs
│   │   ├── middleware/      # Auth, RBAC, CORS
│   │   ├── config/          # Environment config
│   │   └── database/        # Connection pool
│   ├── migrations/          # SQL migrations (001–008)
│   └── tests/               # BDD tests (Gherkin + godog)
├── frontend/
│   ├── src/
│   │   ├── pages/           # 14 page components
│   │   ├── components/      # Layout, UI (shadcn/ui)
│   │   ├── context/         # Auth, language providers
│   │   ├── hooks/           # useTheme
│   │   ├── lib/             # API client, types, i18n, utils
│   │   └── test/            # Test utilities
│   └── package.json
├── docker-compose.yml
└── DOCS.md                  # User guide
```

## Running Tests

```bash
# Backend BDD tests (49 scenarios)
cd backend && go test ./tests/... -v

# Frontend tests (96 tests)
cd frontend && npx vitest run
```

## Roles & Permissions

| Feature | Admin | Member | Child |
|---------|-------|--------|-------|
| Accounts | All family | Own only | Own only |
| Transactions | All family | Own only | Own only |
| Categories | Full CRUD | Read + Create | Read only |
| Budgets | Full CRUD | Read only | No access |
| Reports | All family | Own data | Own data |
| Saving Goals | Full CRUD | Read only | No access |
| Bill Reminders | Full CRUD | Read + Pay | No access |
| Transfers | Yes | Yes | No |
| CSV Import/Export | Yes | Yes | No |
| User Management | Yes | No | No |
| Allowances | Manage all | No | View own |

## Architecture

```
Client → Router → Middleware (auth/RBAC) → Handler → Service → Repository → PostgreSQL
```

- Money stored as integers in cents (never floats)
- UUIDs exposed via API, integer IDs used internally for joins
- Positive amounts = income, negative = expense
- All timestamps in UTC

## License

Private project.
