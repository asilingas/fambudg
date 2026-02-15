package steps

import (
	"context"
	"fmt"
	"os"

	"github.com/asilingas/fambudg/backend/internal/config"
	"github.com/asilingas/fambudg/backend/internal/repository"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/cucumber/godog"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestContext struct {
	Pool               *pgxpool.Pool
	AuthService        *service.AuthService
	TransactionService *service.TransactionService
	AccountService     *service.AccountService
	CategoryService    *service.CategoryService
	BudgetService      *service.BudgetService
	UserRepo           *repository.UserRepository
	AccountRepo        *repository.AccountRepository
	CategoryRepo       *repository.CategoryRepository
	TransactionRepo    *repository.TransactionRepository
	BudgetRepo         *repository.BudgetRepository

	// Test state
	CurrentUser        any
	CurrentToken       string
	CurrentAccount     any
	CurrentCategory    any
	CurrentTransaction any
	CurrentBudget      any
	TransactionList    []any
	BudgetList         []any
	BudgetSummaryList  []any
	LastError          error
	LastStatusCode     int
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := &TestContext{}

	// Setup
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return ctx, tc.setupTestDatabase()
	})

	// Teardown
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		tc.cleanupTestDatabase()
		return ctx, nil
	})

	// Register step definitions
	registerAuthSteps(ctx, tc)
	registerTransactionSteps(ctx, tc)
	registerBudgetSteps(ctx, tc)
}

func (tc *TestContext) setupTestDatabase() error {
	// Set test environment variables
	os.Setenv("DB_NAME", "fambudg_test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-bdd-tests")

	// Set other defaults if not already set
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "budget")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "localdev")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to test database
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, "fambudg_test",
	)

	tc.Pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Initialize repositories
	tc.UserRepo = repository.NewUserRepository(tc.Pool)
	tc.AccountRepo = repository.NewAccountRepository(tc.Pool)
	tc.CategoryRepo = repository.NewCategoryRepository(tc.Pool)
	tc.TransactionRepo = repository.NewTransactionRepository(tc.Pool)
	tc.BudgetRepo = repository.NewBudgetRepository(tc.Pool)

	// Initialize services
	tc.AuthService = service.NewAuthService(tc.UserRepo, cfg.JWT.Secret)
	tc.AccountService = service.NewAccountService(tc.AccountRepo)
	tc.CategoryService = service.NewCategoryService(tc.CategoryRepo)
	tc.TransactionService = service.NewTransactionService(tc.TransactionRepo, tc.AccountRepo)
	tc.BudgetService = service.NewBudgetService(tc.BudgetRepo)

	return nil
}

func (tc *TestContext) cleanupTestDatabase() {
	if tc.Pool != nil {
		// Clean up all tables
		ctx := context.Background()
		tc.Pool.Exec(ctx, "TRUNCATE budgets, transactions, accounts, categories, users CASCADE")
		tc.Pool.Close()
	}
}
