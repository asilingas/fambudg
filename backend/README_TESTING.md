# Testing Guide

## Overview

The Fambudg backend uses two types of tests:
1. **Unit tests** - Test individual functions and services in isolation
2. **BDD tests** - Integration tests using Gherkin syntax with godog

## Prerequisites

- PostgreSQL running via Docker Compose
- Test database `fambudg_test` (automatically created)
- Go 1.21+

## Setup

### 1. Start PostgreSQL with test database

```bash
# From project root
docker compose up -d

# Run migrations on both databases
cd backend
go run cmd/migrate/main.go up

# Run migrations on test database
DB_NAME=fambudg_test go run cmd/migrate/main.go up
```

### 2. Set environment variables

The test suite automatically uses `fambudg_test` database. Make sure your `.env` file is configured:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=budget
DB_PASSWORD=localdev
DB_NAME=fambudg
JWT_SECRET=test-secret-key
SERVER_PORT=8080
```

## Running Tests

### Run all tests

```bash
cd backend
go test ./...
```

### Run unit tests only

```bash
go test ./internal/...
```

### Run BDD tests only

```bash
go test ./tests/... -v
```

### Run BDD tests with pretty output

```bash
go test ./tests/... --godog.format=pretty
```

### Run specific feature

```bash
go test ./tests/... --godog.format=pretty --godog.tags="@auth"
```

## Test Structure

```
backend/
├── internal/
│   └── service/
│       ├── auth_service.go
│       └── auth_service_test.go      # Unit tests
└── tests/
    ├── features/
    │   ├── auth.feature              # BDD scenarios
    │   └── transactions.feature
    └── steps/
        ├── suite_test.go             # Test suite setup
        ├── auth_steps.go             # Step definitions
        └── transaction_steps.go
```

## Writing Tests

### Unit Tests

Create `*_test.go` files next to the code being tested:

```go
func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test"

    // Act
    result := MyFunction(input)

    // Assert
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

### BDD Tests

1. Write feature file in `tests/features/`:

```gherkin
Feature: My Feature
  Scenario: My scenario
    Given some precondition
    When I do something
    Then I should see result
```

2. Implement step definitions in `tests/steps/`:

```go
func registerMySteps(ctx *godog.ScenarioContext, tc *TestContext) {
    ctx.Step(`^some precondition$`, tc.somePrecondition)
}
```

## Test Database

The test database is automatically:
- Created when Docker Compose starts
- Reset between test scenarios (TRUNCATE all tables)
- Isolated from the main development database

## Continuous Integration

Tests run automatically on every push via GitHub Actions (coming in Phase 5).

## Troubleshooting

### Test database doesn't exist

```bash
docker compose down -v
docker compose up -d
# Wait a few seconds for init script
DB_NAME=fambudg_test go run cmd/migrate/main.go up
```

### Tests fail with connection errors

Check that PostgreSQL is running:
```bash
docker compose ps
docker compose logs db
```

### BDD tests fail to find steps

Make sure all step definitions are registered in `suite_test.go`:
```go
func InitializeScenario(ctx *godog.ScenarioContext) {
    registerAuthSteps(ctx, tc)
    registerTransactionSteps(ctx, tc)
}
```
