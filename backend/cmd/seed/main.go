package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asilingas/fambudg/backend/internal/config"
	"github.com/asilingas/fambudg/backend/internal/database"
	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
	"github.com/asilingas/fambudg/backend/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Seeding database: %s\n", cfg.Database.DBName)

	// Connect to database
	pool, err := database.NewPool(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Truncate all tables
	fmt.Print("Truncating tables... ")
	truncate(pool)
	fmt.Println("done")

	// Initialize repositories
	userRepo := repository.NewUserRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	transactionRepo := repository.NewTransactionRepository(pool)
	budgetRepo := repository.NewBudgetRepository(pool)
	savingGoalRepo := repository.NewSavingGoalRepository(pool)
	billReminderRepo := repository.NewBillReminderRepository(pool)
	allowanceRepo := repository.NewAllowanceRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	accountService := service.NewAccountService(accountRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)
	budgetService := service.NewBudgetService(budgetRepo)
	savingGoalService := service.NewSavingGoalService(savingGoalRepo)
	billReminderService := service.NewBillReminderService(billReminderRepo, transactionRepo, accountRepo)
	allowanceService := service.NewAllowanceService(allowanceRepo)

	ctx := context.Background()
	now := time.Now()
	thisYear, thisMonth, _ := now.Date()
	// --- Users ---
	fmt.Print("Seeding users... ")

	admin, err := authService.Register(ctx, &model.RegisterRequest{
		Email:    "admin@family.com",
		Password: "password123",
		Name:     "Arturas",
	})
	if err != nil {
		log.Fatalf("Failed to register admin: %v", err)
	}

	member, err := authService.CreateUser(ctx, &model.CreateUserRequest{
		Email:    "wife@family.com",
		Password: "password123",
		Name:     "Greta",
		Role:     "member",
	})
	if err != nil {
		log.Fatalf("Failed to create member: %v", err)
	}

	morta, err := authService.CreateUser(ctx, &model.CreateUserRequest{
		Email:    "morta@family.com",
		Password: "password123",
		Name:     "Morta",
		Role:     "child",
	})
	if err != nil {
		log.Fatalf("Failed to create child Morta: %v", err)
	}

	agota, err := authService.CreateUser(ctx, &model.CreateUserRequest{
		Email:    "agota@family.com",
		Password: "password123",
		Name:     "Agota",
		Role:     "child",
	})
	if err != nil {
		log.Fatalf("Failed to create child Agota: %v", err)
	}

	pija, err := authService.CreateUser(ctx, &model.CreateUserRequest{
		Email:    "pija@family.com",
		Password: "password123",
		Name:     "Pija",
		Role:     "child",
	})
	if err != nil {
		log.Fatalf("Failed to create child Pija: %v", err)
	}

	juozas, err := authService.CreateUser(ctx, &model.CreateUserRequest{
		Email:    "juozas@family.com",
		Password: "password123",
		Name:     "Juozas",
		Role:     "child",
	})
	if err != nil {
		log.Fatalf("Failed to create child Juozas: %v", err)
	}

	fmt.Println("6 users created")

	// --- Categories ---
	fmt.Print("Seeding categories... ")

	type catDef struct {
		Name string
		Type string
		Icon string
	}
	catDefs := []catDef{
		{"Groceries", "expense", "shopping-cart"},
		{"Rent", "expense", "home"},
		{"Utilities", "expense", "zap"},
		{"Transport", "expense", "car"},
		{"Entertainment", "expense", "film"},
		{"Dining Out", "expense", "utensils"},
		{"Salary", "income", "briefcase"},
		{"Freelance", "income", "laptop"},
	}

	categories := make(map[string]*model.Category) // name -> category
	for i, cd := range catDefs {
		cat, err := categoryService.Create(ctx, &model.CreateCategoryRequest{
			Name:      cd.Name,
			Type:      cd.Type,
			Icon:      cd.Icon,
			SortOrder: i + 1,
		})
		if err != nil {
			log.Fatalf("Failed to create category %s: %v", cd.Name, err)
		}
		categories[cd.Name] = cat
	}

	fmt.Printf("%d categories created\n", len(catDefs))

	// --- Accounts ---
	fmt.Print("Seeding accounts... ")

	type acctDef struct {
		UserID   string
		UserName string
		Name     string
		Type     string
	}
	acctDefs := []acctDef{
		{admin.ID, "Arturas", "Main Checking", "checking"},
		{admin.ID, "Arturas", "Savings", "savings"},
		{member.ID, "Greta", "Greta's Account", "checking"},
		{morta.ID, "Morta", "Morta's Account", "cash"},
		{agota.ID, "Agota", "Agota's Account", "cash"},
		{pija.ID, "Pija", "Pija's Account", "cash"},
		{juozas.ID, "Juozas", "Juozas's Account", "cash"},
	}

	accounts := make(map[string]*model.Account) // name -> account
	for _, ad := range acctDefs {
		acct, err := accountService.Create(ctx, ad.UserID, &model.CreateAccountRequest{
			Name:     ad.Name,
			Type:     ad.Type,
			Currency: "EUR",
			Balance:  0,
		})
		if err != nil {
			log.Fatalf("Failed to create account %s: %v", ad.Name, err)
		}
		accounts[ad.Name] = acct
	}

	fmt.Printf("%d accounts created\n", len(acctDefs))

	// --- Transactions ---
	fmt.Print("Seeding transactions... ")

	type txnDef struct {
		UserID      string
		AccountName string
		Category    string
		Amount      int64
		Type        string
		Description string
		Day         int
		MonthOffset int // 0 = current month, -1 = previous month
		IsShared    bool
		Tags        []string
	}

	txnDefs := []txnDef{
		// Admin income - current month
		{admin.ID, "Main Checking", "Salary", 350000, "income", "Monthly salary", 1, 0, false, []string{"income", "regular"}},
		{admin.ID, "Main Checking", "Freelance", 80000, "income", "Website project", 5, 0, false, []string{"freelance"}},
		// Admin expenses - current month
		{admin.ID, "Main Checking", "Groceries", -12500, "expense", "Maxima weekly shop", 3, 0, true, []string{"groceries"}},
		{admin.ID, "Main Checking", "Groceries", -8700, "expense", "Lidl groceries", 10, 0, true, []string{"groceries"}},
		{admin.ID, "Main Checking", "Rent", -85000, "expense", "Monthly rent", 1, 0, true, []string{"housing", "regular"}},
		{admin.ID, "Main Checking", "Utilities", -12000, "expense", "Electricity bill", 8, 0, true, []string{"bills"}},
		{admin.ID, "Main Checking", "Transport", -4500, "expense", "Fuel", 6, 0, false, []string{"car"}},
		{admin.ID, "Main Checking", "Entertainment", -2500, "expense", "Netflix subscription", 4, 0, true, []string{"subscription"}},
		{admin.ID, "Main Checking", "Dining Out", -3800, "expense", "Restaurant dinner", 7, 0, true, []string{"dining"}},
		// Admin expenses - previous month
		{admin.ID, "Main Checking", "Salary", 350000, "income", "Monthly salary", 1, -1, false, []string{"income", "regular"}},
		{admin.ID, "Main Checking", "Groceries", -15200, "expense", "Maxima weekly shop", 5, -1, true, []string{"groceries"}},
		{admin.ID, "Main Checking", "Groceries", -9300, "expense", "Rimi groceries", 15, -1, true, []string{"groceries"}},
		{admin.ID, "Main Checking", "Rent", -85000, "expense", "Monthly rent", 1, -1, true, []string{"housing", "regular"}},
		{admin.ID, "Main Checking", "Dining Out", -5200, "expense", "Family dinner out", 20, -1, true, []string{"dining"}},
		// Greta income + expenses
		{member.ID, "Greta's Account", "Salary", 280000, "income", "Monthly salary", 1, 0, false, []string{"income", "regular"}},
		{member.ID, "Greta's Account", "Groceries", -6500, "expense", "IKI groceries", 4, 0, true, []string{"groceries"}},
		{member.ID, "Greta's Account", "Entertainment", -1500, "expense", "Spotify subscription", 5, 0, false, []string{"subscription"}},
		{member.ID, "Greta's Account", "Dining Out", -4200, "expense", "Lunch with colleagues", 9, 0, false, []string{"dining", "work"}},
		// Morta expenses
		{morta.ID, "Morta's Account", "Entertainment", -1200, "expense", "Mobile game", 3, 0, false, []string{"gaming"}},
		{morta.ID, "Morta's Account", "Dining Out", -800, "expense", "School canteen", 6, 0, false, []string{"school"}},
		// Agota expenses
		{agota.ID, "Agota's Account", "Entertainment", -900, "expense", "Art supplies", 4, 0, false, []string{"hobby"}},
		{agota.ID, "Agota's Account", "Dining Out", -650, "expense", "School canteen", 7, 0, false, []string{"school"}},
		// Pija expenses
		{pija.ID, "Pija's Account", "Entertainment", -1100, "expense", "Books", 5, 0, false, []string{"books"}},
		{pija.ID, "Pija's Account", "Dining Out", -750, "expense", "School canteen", 8, 0, false, []string{"school"}},
		// Juozas expenses
		{juozas.ID, "Juozas's Account", "Entertainment", -1400, "expense", "Football gear", 2, 0, false, []string{"sports"}},
		{juozas.ID, "Juozas's Account", "Dining Out", -850, "expense", "School canteen", 5, 0, false, []string{"school"}},
	}

	txnCount := 0
	for _, td := range txnDefs {
		txnMonth := thisMonth + time.Month(td.MonthOffset)
		txnYear := thisYear
		if txnMonth < 1 {
			txnMonth += 12
			txnYear--
		}
		date := fmt.Sprintf("%04d-%02d-%02d", txnYear, txnMonth, td.Day)

		_, err := transactionService.Create(ctx, td.UserID, &model.CreateTransactionRequest{
			AccountID:   accounts[td.AccountName].ID,
			CategoryID:  categories[td.Category].ID,
			Amount:      td.Amount,
			Type:        td.Type,
			Description: td.Description,
			Date:        date,
			IsShared:    td.IsShared,
			Tags:        td.Tags,
		})
		if err != nil {
			log.Fatalf("Failed to create transaction '%s': %v", td.Description, err)
		}
		txnCount++
	}

	fmt.Printf("%d transactions created\n", txnCount)

	// --- Budgets ---
	fmt.Print("Seeding budgets... ")

	type budgetDef struct {
		Category string
		Amount   int64
	}
	budgetDefs := []budgetDef{
		{"Groceries", 40000},
		{"Dining Out", 15000},
		{"Entertainment", 10000},
		{"Transport", 8000},
	}

	for _, bd := range budgetDefs {
		_, err := budgetService.Create(ctx, &model.CreateBudgetRequest{
			CategoryID: categories[bd.Category].ID,
			Amount:     bd.Amount,
			Month:      int(thisMonth),
			Year:       thisYear,
		})
		if err != nil {
			log.Fatalf("Failed to create budget for %s: %v", bd.Category, err)
		}
	}

	fmt.Printf("%d budgets created\n", len(budgetDefs))

	// --- Saving Goals ---
	fmt.Print("Seeding saving goals... ")

	targetDate1 := fmt.Sprintf("%04d-12-31", thisYear)
	targetDate2 := fmt.Sprintf("%04d-06-30", thisYear+1)

	goal1, err := savingGoalService.Create(ctx, &model.CreateSavingGoalRequest{
		Name:         "Family Vacation",
		TargetAmount: 300000,
		TargetDate:   &targetDate1,
		Priority:     1,
	})
	if err != nil {
		log.Fatalf("Failed to create saving goal: %v", err)
	}
	_, err = savingGoalService.Contribute(ctx, goal1.ID, &model.ContributeRequest{Amount: 120000})
	if err != nil {
		log.Fatalf("Failed to contribute to saving goal: %v", err)
	}

	goal2, err := savingGoalService.Create(ctx, &model.CreateSavingGoalRequest{
		Name:         "New Laptop",
		TargetAmount: 150000,
		TargetDate:   &targetDate2,
		Priority:     2,
	})
	if err != nil {
		log.Fatalf("Failed to create saving goal: %v", err)
	}
	_, err = savingGoalService.Contribute(ctx, goal2.ID, &model.ContributeRequest{Amount: 50000})
	if err != nil {
		log.Fatalf("Failed to contribute to saving goal: %v", err)
	}

	fmt.Println("2 saving goals created (with contributions)")

	// --- Bill Reminders ---
	fmt.Print("Seeding bill reminders... ")

	type billDef struct {
		Name      string
		Amount    int64
		DueDay    int
		Category  string
		AccountID string
	}
	billDefs := []billDef{
		{"Rent", 85000, 1, "Rent", "Main Checking"},
		{"Internet", 3500, 15, "Utilities", "Main Checking"},
		{"Electricity", 12000, 20, "Utilities", "Main Checking"},
	}

	for _, bd := range billDefs {
		catID := categories[bd.Category].ID
		acctID := accounts[bd.AccountID].ID
		nextDue := fmt.Sprintf("%04d-%02d-%02d", thisYear, thisMonth, bd.DueDay)

		_, err := billReminderService.Create(ctx, &model.CreateBillReminderRequest{
			Name:        bd.Name,
			Amount:      bd.Amount,
			DueDay:      bd.DueDay,
			Frequency:   "monthly",
			CategoryID:  &catID,
			AccountID:   &acctID,
			NextDueDate: nextDue,
		})
		if err != nil {
			log.Fatalf("Failed to create bill reminder %s: %v", bd.Name, err)
		}
	}

	fmt.Printf("%d bill reminders created\n", len(billDefs))

	// --- Allowances ---
	fmt.Print("Seeding allowances... ")

	periodStart := fmt.Sprintf("%04d-%02d-01", thisYear, thisMonth)
	childAllowances := []struct {
		userID string
		name   string
		amount int64
	}{
		{morta.ID, "Morta", 5000},
		{agota.ID, "Agota", 4000},
		{pija.ID, "Pija", 4000},
		{juozas.ID, "Juozas", 5000},
	}

	for _, ca := range childAllowances {
		_, err = allowanceService.Create(ctx, &model.CreateAllowanceRequest{
			UserID:      ca.userID,
			Amount:      ca.amount,
			PeriodStart: periodStart,
		})
		if err != nil {
			log.Fatalf("Failed to create allowance for %s: %v", ca.name, err)
		}
	}

	fmt.Printf("%d allowances created\n", len(childAllowances))

	// --- Summary ---
	fmt.Println()
	fmt.Println("Seed complete!")
	fmt.Printf("  Users:          6\n")
	fmt.Printf("  Categories:     %d\n", len(catDefs))
	fmt.Printf("  Accounts:       %d\n", len(acctDefs))
	fmt.Printf("  Transactions:   %d\n", txnCount)
	fmt.Printf("  Budgets:        %d\n", len(budgetDefs))
	fmt.Printf("  Saving Goals:   2\n")
	fmt.Printf("  Bill Reminders: %d\n", len(billDefs))
	fmt.Printf("  Allowances:     %d\n", len(childAllowances))
	fmt.Println()
	fmt.Println("Login as admin@family.com / password123")

}

func truncate(pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE transactions, allowances, bill_reminders, saving_goals, budgets, accounts, categories, users CASCADE")
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v", err)
	}
}
