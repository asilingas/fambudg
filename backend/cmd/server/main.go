package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/asilingas/fambudg/backend/internal/config"
	"github.com/asilingas/fambudg/backend/internal/database"
	"github.com/asilingas/fambudg/backend/internal/handler"
	"github.com/asilingas/fambudg/backend/internal/middleware"
	"github.com/asilingas/fambudg/backend/internal/repository"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	pool, err := database.NewPool(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	transactionRepo := repository.NewTransactionRepository(pool)
	budgetRepo := repository.NewBudgetRepository(pool)
	reportRepo := repository.NewReportRepository(pool)
	savingGoalRepo := repository.NewSavingGoalRepository(pool)
	billReminderRepo := repository.NewBillReminderRepository(pool)
	allowanceRepo := repository.NewAllowanceRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	accountService := service.NewAccountService(accountRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)
	budgetService := service.NewBudgetService(budgetRepo)
	reportService := service.NewReportService(reportRepo, accountRepo)
	savingGoalService := service.NewSavingGoalService(savingGoalRepo)
	billReminderService := service.NewBillReminderService(billReminderRepo, transactionRepo, accountRepo)
	allowanceService := service.NewAllowanceService(allowanceRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService)
	accountHandler := handler.NewAccountHandler(accountService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	budgetHandler := handler.NewBudgetHandler(budgetService)
	reportHandler := handler.NewReportHandler(reportService)
	savingGoalHandler := handler.NewSavingGoalHandler(savingGoalService)
	billReminderHandler := handler.NewBillReminderHandler(billReminderService)
	transferHandler := handler.NewTransferHandler(transactionService)
	importExportHandler := handler.NewImportExportHandler(transactionService)
	allowanceHandler := handler.NewAllowanceHandler(allowanceService)

	// Create router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)

	// Public routes
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	// All authenticated users (admin, member, child)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))

		// Auth
		r.Get("/api/auth/me", authHandler.GetMe)

		// Accounts (admin sees all, others see own; ownership checks in handler)
		r.Get("/api/accounts", accountHandler.List)
		r.Post("/api/accounts", accountHandler.Create)
		r.Get("/api/accounts/{id}", accountHandler.Get)
		r.Put("/api/accounts/{id}", accountHandler.Update)
		r.Delete("/api/accounts/{id}", accountHandler.Delete)

		// Categories (read for all)
		r.Get("/api/categories", categoryHandler.List)

		// Transactions (admin sees all, others see own; ownership checks in handler)
		r.Get("/api/transactions", transactionHandler.List)
		r.Post("/api/transactions", transactionHandler.Create)
		r.Get("/api/transactions/{id}", transactionHandler.Get)
		r.Put("/api/transactions/{id}", transactionHandler.Update)
		r.Delete("/api/transactions/{id}", transactionHandler.Delete)

		// Reports (admin sees all data, others see own)
		r.Get("/api/reports/dashboard", reportHandler.Dashboard)
		r.Get("/api/reports/monthly", reportHandler.Monthly)
		r.Get("/api/reports/by-category", reportHandler.ByCategory)
		r.Get("/api/reports/trends", reportHandler.Trends)

		// Search (admin searches all, others search own)
		r.Get("/api/search", reportHandler.Search)

		// Allowances (list: admin sees all, child sees own)
		r.Get("/api/allowances", allowanceHandler.List)
	})

	// Admin + Member routes (read access to budgets/goals/reminders + write categories)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))
		r.Use(middleware.RequireRole("admin", "member"))

		// Categories create (admin + member)
		r.Post("/api/categories", categoryHandler.Create)

		// Budgets read (admin + member)
		r.Get("/api/budgets", budgetHandler.List)
		r.Get("/api/budgets/summary", budgetHandler.Summary)

		// Saving Goals read (admin + member)
		r.Get("/api/saving-goals", savingGoalHandler.List)

		// Bill Reminders read + pay (admin + member)
		r.Get("/api/bill-reminders", billReminderHandler.List)
		r.Get("/api/bill-reminders/upcoming", billReminderHandler.Upcoming)
		r.Post("/api/bill-reminders/{id}/pay", billReminderHandler.Pay)

		// Transfers (admin + member)
		r.Post("/api/transfers", transferHandler.Create)

		// Recurring transactions (admin + member)
		r.Post("/api/transactions/generate-recurring", transactionHandler.GenerateRecurring)

		// Import / Export (admin + member)
		r.Post("/api/import/csv", importExportHandler.ImportCSV)
		r.Get("/api/export/csv", importExportHandler.ExportCSV)
	})

	// Admin-only routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))
		r.Use(middleware.RequireRole("admin"))

		// User management
		r.Get("/api/users", userHandler.List)
		r.Post("/api/users", userHandler.Create)
		r.Put("/api/users/{id}", userHandler.Update)
		r.Delete("/api/users/{id}", userHandler.Delete)

		// Categories update/delete (admin only)
		r.Put("/api/categories/{id}", categoryHandler.Update)
		r.Delete("/api/categories/{id}", categoryHandler.Delete)

		// Budgets write (admin only)
		r.Post("/api/budgets", budgetHandler.Create)
		r.Put("/api/budgets/{id}", budgetHandler.Update)
		r.Delete("/api/budgets/{id}", budgetHandler.Delete)

		// Saving Goals write (admin only)
		r.Post("/api/saving-goals", savingGoalHandler.Create)
		r.Put("/api/saving-goals/{id}", savingGoalHandler.Update)
		r.Post("/api/saving-goals/{id}/contribute", savingGoalHandler.Contribute)

		// Bill Reminders write (admin only)
		r.Post("/api/bill-reminders", billReminderHandler.Create)
		r.Put("/api/bill-reminders/{id}", billReminderHandler.Update)
		r.Delete("/api/bill-reminders/{id}", billReminderHandler.Delete)

		// Family member spending comparison
		r.Get("/api/reports/by-member", reportHandler.ByMember)

		// Allowances management (admin only)
		r.Post("/api/allowances", allowanceHandler.Create)
		r.Put("/api/allowances/{id}", allowanceHandler.Update)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
