# Fambudg — Backend Plan

## Goal
**Fambudg** — a personal family budget web app for tracking expenses and income. Users: me (admin), wife (later), children (later). Purpose: see clear financial picture, plan savings, and track goals for larger purchases.

## Tech Stack
- **Backend:** Go (Golang)
- **Frontend:** React (later)
- **Database:** PostgreSQL
- **Go driver:** pgx/v5
- **HTTP router:** chi/v5
- **Migrations:** goose/v3
- **Auth:** golang-jwt/jwt/v5
- **Validation:** go-playground/validator/v10
- **Testing:** Go standard testing + godog (BDD, Gherkin syntax — like Behat in PHP)
- **Dev environment:** macOS (MacBook Pro), PostgreSQL via Docker
- **DB GUI:** TablePlus (supports both Postgres and MySQL)
- **API testing:** Postman
- **Version control:** Git + GitHub (private repo)
- **Issue tracking:** GitHub Issues + GitHub Projects (Kanban board)
- **Claude Code MCP:** GitHub MCP server for automated issue management

---

## GitHub + Issue Tracking Setup

```bash
# 1. Create GitHub repo
# Go to github.com → New Repository → "fambudg" → Private

# 2. Generate a Personal Access Token
# github.com → Settings → Developer settings → Personal access tokens → Fine-grained tokens
# Scope: "repo" access for your fambudg repository

# 3. Set up GitHub MCP for Claude Code
# Create/edit ~/.claude/claude_code_config.json:
```

```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your-token-here"
      }
    }
  }
}
```

```bash
# 4. Set up GitHub Projects board
# Go to your repo → Projects → New Project → Board
# Columns: Backlog → In Progress → Review → Done
# Create milestones: Phase 1, Phase 2, Phase 3, Phase 4, Phase 5

# 5. Create labels on the repo:
# feature, bugfix, enhancement, documentation, testing
# priority:high, priority:medium, priority:low
# phase:1, phase:2, phase:3, phase:4, phase:5
```

With this setup, Claude Code can automatically:
- Create issues for each task (`"Create issues for all Phase 1 tasks"`)
- Label and assign milestones (`"Label this as bugfix, priority:high"`)
- Close issues from the terminal (`"Close issue #5, it's done"`)
- List open issues (`"Show me all open Phase 1 issues"`)
- Reference issues in commits (`"Commit this with message 'Add transaction CRUD, fixes #3'"`)


---

## Project Structure

```
fambudg/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── model/
│   │   ├── user.go
│   │   ├── account.go
│   │   ├── category.go
│   │   ├── transaction.go
│   │   ├── budget.go
│   │   └── saving_goal.go
│   │   ├── bill_reminder.go
│   │   └── allowance.go
│   ├── repository/
│   │   ├── user_repo.go
│   │   ├── account_repo.go
│   │   ├── category_repo.go
│   │   ├── transaction_repo.go
│   │   ├── budget_repo.go
│   │   └── saving_goal_repo.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── transaction_service.go
│   │   ├── budget_service.go
│   │   └── report_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── transaction_handler.go
│   │   ├── account_handler.go
│   │   ├── category_handler.go
│   │   ├── budget_handler.go
│   │   ├── saving_goal_handler.go
│   │   ├── bill_reminder_handler.go
│   │   ├── transfer_handler.go
│   │   ├── search_handler.go
│   │   └── import_export_handler.go
│   └── middleware/
│       ├── auth.go
│       └── logging.go
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_accounts.sql
│   ├── 003_create_categories.sql
│   ├── 004_create_transactions.sql
│   ├── 005_create_budgets.sql
│   └── 006_create_saving_goals.sql
│   ├── 007_create_bill_reminders.sql
│   └── 008_create_allowances.sql
├── .claude/
│   └── plan.md                # full project plan
├── CLAUDE.md                  # project context for Claude Code
├── docker-compose.yml
├── .env
├── go.mod
└── go.sum
```

---

## Data Model

### Users
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    | gen_random_uuid()                  |
| email        | VARCHAR(255) | unique, login identifier           |
| password_hash| TEXT         | bcrypt hashed                      |
| name         | VARCHAR(100) | display name                       |
| role         | VARCHAR(20)  | admin / member / child             |
| created_at   | TIMESTAMPTZ  |                                    |
| updated_at   | TIMESTAMPTZ  |                                    |

### Accounts (bank accounts, wallets, credit cards)
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    |                                    |
| user_id      | UUID (FK)    | owner                              |
| name         | VARCHAR(100) | "Chase Checking", "Cash", etc.     |
| type         | VARCHAR(20)  | checking / savings / credit / cash |
| currency     | VARCHAR(3)   | USD, EUR, etc.                     |
| balance      | BIGINT       | stored in cents (avoid floats)     |
| created_at   | TIMESTAMPTZ  |                                    |

### Categories
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    |                                    |
| parent_id    | UUID (FK)    | nullable, for subcategories        |
| name         | VARCHAR(100) | "Groceries", "Rent", etc.          |
| type         | VARCHAR(10)  | expense / income                   |
| icon         | VARCHAR(50)  | optional, for frontend             |
| sort_order   | INT          | display ordering                   |

### Transactions
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    |                                    |
| user_id      | UUID (FK)    | who made this transaction          |
| account_id   | UUID (FK)    | from which account                 |
| category_id  | UUID (FK)    |                                    |
| amount       | BIGINT       | in cents. Positive=income, Negative=expense |
| type         | VARCHAR(10)  | expense / income / transfer        |
| description  | TEXT         |                                    |
| date         | DATE         | transaction date                   |
| is_shared    | BOOLEAN      | default true (family vs personal)  |
| is_recurring | BOOLEAN      | default false                      |
| recurring_rule| JSONB       | nullable: {"frequency":"monthly","day":1} |
| tags         | TEXT[]       | optional tags: ["date-night","business"] |
| transfer_to_account_id | UUID (FK) | nullable, for transfers only |
| created_at   | TIMESTAMPTZ  |                                    |
| updated_at   | TIMESTAMPTZ  |                                    |

### Budgets
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    |                                    |
| category_id  | UUID (FK)    |                                    |
| amount       | BIGINT       | budget limit in cents              |
| month        | INT          | 1-12                               |
| year         | INT          |                                    |
| created_at   | TIMESTAMPTZ  |                                    |

### Saving Goals
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    |                                    |
| name         | VARCHAR(200) | "New Car", "Family Vacation"       |
| target_amount| BIGINT       | in cents                           |
| current_amount| BIGINT      | in cents                           |
| target_date  | DATE         | nullable                           |
| priority     | INT          | 1=highest                          |
| status       | VARCHAR(20)  | active / completed / cancelled     |
| created_at   | TIMESTAMPTZ  |                                    |
| updated_at   | TIMESTAMPTZ  |                                    |

### Bill Reminders
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    |                                    |
| name         | VARCHAR(200) | "Electricity", "Internet"          |
| amount       | BIGINT       | expected amount in cents           |
| due_day      | INT          | day of month (1-31)                |
| frequency    | VARCHAR(20)  | monthly / quarterly / yearly       |
| category_id  | UUID (FK)    | nullable                           |
| account_id   | UUID (FK)    | nullable, which account to pay from|
| is_active    | BOOLEAN      | default true                       |
| next_due_date| DATE         | next upcoming due date             |
| created_at   | TIMESTAMPTZ  |                                    |
| updated_at   | TIMESTAMPTZ  |                                    |

### Allowances (for children — Phase 4)
| Column       | Type         | Notes                              |
|--------------|--------------|------------------------------------|
| id           | UUID (PK)    |                                    |
| user_id      | UUID (FK)    | the child                          |
| amount       | BIGINT       | monthly allowance in cents         |
| spent        | BIGINT       | spent this period in cents         |
| period_start | DATE         | current allowance period start     |
| created_at   | TIMESTAMPTZ  |                                    |
| updated_at   | TIMESTAMPTZ  |                                    |

---

## API Endpoints

### Auth
- `POST   /api/auth/register`   — create user (admin only later)
- `POST   /api/auth/login`      — returns JWT
- `GET    /api/auth/me`         — current user profile

### Transactions
- `GET    /api/transactions`          — list (with filters: date range, category, account)
- `POST   /api/transactions`          — create
- `GET    /api/transactions/:id`      — get one
- `PUT    /api/transactions/:id`      — update
- `DELETE /api/transactions/:id`      — delete

### Accounts
- `GET    /api/accounts`              — list all
- `POST   /api/accounts`              — create
- `PUT    /api/accounts/:id`          — update
- `DELETE /api/accounts/:id`          — delete

### Categories
- `GET    /api/categories`            — list (tree structure)
- `POST   /api/categories`            — create
- `PUT    /api/categories/:id`        — update
- `DELETE /api/categories/:id`        — delete

### Budgets
- `GET    /api/budgets`               — list (filter by month/year)
- `POST   /api/budgets`               — create/update budget for category+month
- `GET    /api/budgets/summary`       — budget vs actual per category

### Saving Goals
- `GET    /api/saving-goals`          — list all
- `POST   /api/saving-goals`          — create
- `PUT    /api/saving-goals/:id`      — update
- `POST   /api/saving-goals/:id/contribute` — add money to goal

### Reports
- `GET    /api/reports/dashboard`     — overview: balances, month summary, recent transactions, goal progress
- `GET    /api/reports/monthly`       — income vs expense summary
- `GET    /api/reports/by-category`   — spending breakdown
- `GET    /api/reports/by-member`     — spending per family member
- `GET    /api/reports/trends`        — month-over-month trends

### Transfers
- `POST   /api/transfers`             — transfer between accounts

### Bill Reminders
- `GET    /api/bill-reminders`         — list all (with upcoming due dates)
- `POST   /api/bill-reminders`         — create
- `PUT    /api/bill-reminders/:id`     — update
- `DELETE /api/bill-reminders/:id`     — delete
- `POST   /api/bill-reminders/:id/pay` — mark as paid (creates transaction)

### Search
- `GET    /api/search`                 — search transactions by description, amount range, date range, category, account, tags

### Import / Export
- `POST   /api/import/csv`            — import transactions from CSV
- `GET    /api/export/csv`             — export transactions to CSV

### Allowances (Phase 4)
- `GET    /api/allowances`             — list all
- `POST   /api/allowances`             — set allowance for a child
- `PUT    /api/allowances/:id`         — update allowance

---

## Key Design Decisions

1. **Money stored as BIGINT in cents** — never use float for money. $19.99 = 1999 cents.
2. **UUID primary keys** — safe for distributed systems, no sequential ID guessing.
3. **Repository pattern** — all DB access goes through repository layer. Makes testing and future DB swaps easy.
4. **JWT auth** — stateless, simple. Short-lived access token + refresh token.
5. **All amounts signed** — positive = income, negative = expense. Simplifies summing.
6. **Soft category hierarchy** — one level of parent/child via parent_id. Keep it simple.
7. **Transfers are paired transactions** — a transfer creates two linked transactions (debit from one account, credit to another) so account balances stay correct.
8. **Shared vs personal** — `is_shared` flag on transactions. Shared = family expense, personal = individual. Defaults to shared.
9. **Bill reminders track next_due_date** — auto-advances after payment so upcoming bills are always current.

---

## Build Phases

### Phase 1 — Foundation
- [x] Init Go module, project structure
- [x] Git init + .gitignore + push to GitHub (private repo)
- [x] Set up GitHub Projects board + milestones + labels
- [x] Create GitHub Issues for Phase 1 tasks
- [x] docker-compose.yml for Postgres
- [x] Database migrations (users, accounts, categories, transactions)
- [x] Config loading from .env
- [x] DB connection pool setup
- [x] User registration + login with JWT
- [x] Transaction CRUD endpoints
- [x] Category CRUD endpoints
- [x] Account CRUD endpoints
- [x] Basic middleware (auth, logging, CORS)
- [x] Unit tests for services
- [x] godog BDD setup + first feature files (auth, transactions)
- [x] Test database setup and teardown

### Phase 2 — Budgeting & Reports
- [x] Budget CRUD (set monthly limits per category)
- [x] Budget vs actual summary endpoint
- [x] Dashboard overview endpoint (balances, month summary, recent transactions, goal progress)
- [x] Monthly income/expense report
- [x] Spending by category report
- [x] Spending by family member report
- [x] Search endpoint (by description, amount, date, category, account, tags)

### Phase 3 — Savings, Recurring & Bills
- [x] Saving goals CRUD
- [x] Contribute to saving goal endpoint
- [x] Recurring transaction rules
- [x] Recurring transaction auto-generation (cron or on-demand)
- [x] Month-over-month trend report
- [x] Bill reminders CRUD
- [x] Mark bill as paid (auto-creates transaction)
- [x] Upcoming bills endpoint
- [x] Transfers between accounts
- [x] CSV import (transactions)
- [x] CSV export (transactions)

### Phase 4 — Multi-User & Permissions
- [ ] User roles (admin, member, child)
- [ ] Role-based access control middleware
- [ ] Admin can manage all users
- [ ] Children see limited data (own spending only)
- [ ] Shared vs personal expense flag (is_shared on transactions)
- [ ] Allowances for children (set monthly limit, track spending)
- [ ] Family member spending comparison

### Phase 5 — Cloud Deployment
- [ ] Dockerfile for the Go backend
- [ ] Choose managed Postgres (Neon / Supabase / Railway)
- [ ] Environment-based config (dev vs prod)
- [ ] Health check endpoint
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] HTTPS + domain setup

---

## Local Development Setup

```bash
# 1. Create project
mkdir fambudg && cd fambudg

# 2. Initialize Git
git init
# Create .gitignore (see below)

# 3. Create GitHub repo
# Go to github.com → New Repository → "fambudg" → Private
git remote add origin git@github.com:yourusername/fambudg.git

# 4. Create backend
mkdir -p backend
cd backend
go mod init github.com/yourusername/fambudg/backend
cd ..

# 5. Start Postgres
docker compose up -d

# 6. Run migrations
cd backend
go run cmd/migrate/main.go up

# 7. Start server
go run cmd/server/main.go
# Runs on http://localhost:8080
```

## .gitignore

```
# Binaries
/backend/bin/
*.exe

# Environment
.env

# Database file (if ever using SQLite)
*.db

# IDE
.idea/
.vscode/
*.swp

# OS
.DS_Store

# Go
coverage.out

# Node / React (for frontend later)
node_modules/
/frontend/build/
/frontend/dist/

# Vendor (optional)
# /backend/vendor/
```

## Docker Compose (local dev)

```yaml
services:
  db:
    image: postgres:17
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: fambudg
      POSTGRES_USER: budget
      POSTGRES_PASSWORD: localdev
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

## Environment Variables (.env)

**Important:** `.env` is in `.gitignore` — never commit secrets. Instead, commit a `.env.example` with placeholder values so others know what's needed.

`.env.example` (committed to Git, in project root):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=budget
DB_PASSWORD=changeme
DB_NAME=fambudg
JWT_SECRET=changeme
SERVER_PORT=8080
```

`backend/.env` (local only, not committed):
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=budget
DB_PASSWORD=localdev
DB_NAME=fambudg
JWT_SECRET=change-me-in-production
SERVER_PORT=8080
```

---

## Testing Strategy

### Unit Tests (Go standard `testing` package)
Test individual functions and methods in isolation. Every service and repository gets a `_test.go` file next to it.

```
internal/
├── service/
│   ├── transaction_service.go
│   └── transaction_service_test.go
├── repository/
│   ├── transaction_repo.go
│   └── transaction_repo_test.go
```

Run with: `go test ./...`

### BDD / Integration Tests (godog — Gherkin syntax)
godog is the Go equivalent of Behat. Same Gherkin feature files, same Given/When/Then flow.

```
tests/
├── features/
│   ├── auth.feature
│   ├── transactions.feature
│   ├── categories.feature
│   ├── budgets.feature
│   └── saving_goals.feature
└── steps/
    ├── auth_steps.go
    ├── transaction_steps.go
    └── common_steps.go
```

Example feature file (`transactions.feature`):

```gherkin
Feature: Transaction management

  Background:
    Given I am logged in as "admin@family.com"

  Scenario: Create an expense transaction
    Given a category "Groceries" exists
    And an account "Chase Checking" exists
    When I create a transaction with:
      | amount      | -4599         |
      | description | Weekly groceries |
      | category    | Groceries     |
      | account     | Chase Checking|
      | date        | 2026-02-14    |
    Then the response status should be 201
    And the transaction should have amount -4599

  Scenario: List transactions filtered by date range
    Given the following transactions exist:
      | amount | description | date       |
      | -2000  | Coffee      | 2026-01-15 |
      | -5000  | Dinner      | 2026-02-10 |
      | -3000  | Lunch       | 2026-02-12 |
    When I list transactions from "2026-02-01" to "2026-02-28"
    Then I should see 2 transactions
```

Run with: `go test ./tests/... --godog.format=pretty`

### Test Database
Tests use a separate Postgres database (`fambudg_test`) that gets reset between test runs. Add to docker-compose:

```yaml
services:
  db:
    environment:
      POSTGRES_MULTIPLE_DATABASES: fambudg,fambudg_test
```

### Build Phase for Testing
Testing is added in Phase 1 alongside the first CRUD endpoints:
- [ ] Set up Go unit test structure
- [ ] Set up godog with first feature file
- [ ] Test database setup and teardown
- [ ] CI-friendly test runner (`go test ./...`)

---

## Claude Code — Getting Started

```bash
cd fambudg
claude
```

Prompt for Phase 1:
> "Here's my project plan for Fambudg, a family budget Go backend. Start with Phase 1: init the Go module, create docker-compose.yml, write the database migrations, set up config loading, DB connection, and basic auth + transaction CRUD. Use pgx, chi, goose, and golang-jwt."