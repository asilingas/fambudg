package repository

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetMonthSummary(ctx context.Context, userID string, month, year int) (*model.MonthSummary, error) {
	summary := &model.MonthSummary{
		Month: month,
		Year:  year,
	}

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) AS total_expense
		FROM transactions
		WHERE user_id = $1
			AND EXTRACT(MONTH FROM date) = $2
			AND EXTRACT(YEAR FROM date) = $3
	`

	err := r.db.QueryRow(ctx, query, userID, month, year).Scan(
		&summary.TotalIncome, &summary.TotalExpense,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get month summary: %w", err)
	}

	summary.Net = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}

func (r *ReportRepository) GetRecentTransactions(ctx context.Context, userID string, limit int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY date DESC, created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		t := &model.Transaction{}
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.AccountID, &t.CategoryID,
			&t.Amount, &t.Type, &t.Description, &t.Date,
			&t.IsShared, &t.IsRecurring, &t.RecurringRule,
			&t.Tags, &t.TransferToAccountID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetMonthSummaryAll returns month summary for all users (admin)
func (r *ReportRepository) GetMonthSummaryAll(ctx context.Context, month, year int) (*model.MonthSummary, error) {
	summary := &model.MonthSummary{
		Month: month,
		Year:  year,
	}

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) AS total_expense
		FROM transactions
		WHERE EXTRACT(MONTH FROM date) = $1
			AND EXTRACT(YEAR FROM date) = $2
	`

	err := r.db.QueryRow(ctx, query, month, year).Scan(
		&summary.TotalIncome, &summary.TotalExpense,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get month summary: %w", err)
	}

	summary.Net = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}

// GetRecentTransactionsAll returns recent transactions for all users (admin)
func (r *ReportRepository) GetRecentTransactionsAll(ctx context.Context, limit int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
		FROM transactions
		ORDER BY date DESC, created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		t := &model.Transaction{}
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.AccountID, &t.CategoryID,
			&t.Amount, &t.Type, &t.Description, &t.Date,
			&t.IsShared, &t.IsRecurring, &t.RecurringRule,
			&t.Tags, &t.TransferToAccountID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// GetSpendingByCategoryAll returns spending by category for all users (admin)
func (r *ReportRepository) GetSpendingByCategoryAll(ctx context.Context, month, year int) ([]*model.CategorySpending, error) {
	query := `
		SELECT
			t.category_id,
			c.name AS category_name,
			COALESCE(SUM(ABS(t.amount)), 0) AS total_amount
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE t.amount < 0
			AND EXTRACT(MONTH FROM t.date) = $1
			AND EXTRACT(YEAR FROM t.date) = $2
		GROUP BY t.category_id, c.name
		ORDER BY total_amount DESC
	`

	rows, err := r.db.Query(ctx, query, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending by category: %w", err)
	}
	defer rows.Close()

	var results []*model.CategorySpending
	var grandTotal int64
	for rows.Next() {
		cs := &model.CategorySpending{}
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.TotalAmount); err != nil {
			return nil, fmt.Errorf("failed to scan category spending: %w", err)
		}
		grandTotal += cs.TotalAmount
		results = append(results, cs)
	}

	for _, cs := range results {
		if grandTotal > 0 {
			cs.Percentage = float64(cs.TotalAmount) / float64(grandTotal) * 100
		}
	}

	return results, nil
}

// GetTrendsAll returns trends for all users (admin)
func (r *ReportRepository) GetTrendsAll(ctx context.Context, months int) ([]*model.TrendPoint, error) {
	query := `
		SELECT
			EXTRACT(MONTH FROM date)::int AS month,
			EXTRACT(YEAR FROM date)::int AS year,
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) AS total_expense
		FROM transactions
		WHERE date >= (CURRENT_DATE - make_interval(months => $1))
		GROUP BY EXTRACT(YEAR FROM date), EXTRACT(MONTH FROM date)
		ORDER BY year, month
	`

	rows, err := r.db.Query(ctx, query, months)
	if err != nil {
		return nil, fmt.Errorf("failed to get trends: %w", err)
	}
	defer rows.Close()

	var trends []*model.TrendPoint
	for rows.Next() {
		tp := &model.TrendPoint{}
		if err := rows.Scan(&tp.Month, &tp.Year, &tp.TotalIncome, &tp.TotalExpense); err != nil {
			return nil, fmt.Errorf("failed to scan trend point: %w", err)
		}
		tp.Net = tp.TotalIncome - tp.TotalExpense
		trends = append(trends, tp)
	}

	return trends, nil
}

// SearchTransactionsAll searches all transactions without user filter (admin)
func (r *ReportRepository) SearchTransactionsAll(ctx context.Context, filters *model.SearchFilters) (*model.SearchResult, error) {
	query := `
		SELECT id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
		FROM transactions
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filters.Description != "" {
		query += fmt.Sprintf(" AND description ILIKE $%d", argPos)
		args = append(args, "%"+filters.Description+"%")
		argPos++
	}

	if filters.MinAmount != nil {
		query += fmt.Sprintf(" AND ABS(amount) >= $%d", argPos)
		args = append(args, *filters.MinAmount)
		argPos++
	}

	if filters.MaxAmount != nil {
		query += fmt.Sprintf(" AND ABS(amount) <= $%d", argPos)
		args = append(args, *filters.MaxAmount)
		argPos++
	}

	if filters.StartDate != "" {
		query += fmt.Sprintf(" AND date >= $%d", argPos)
		args = append(args, filters.StartDate)
		argPos++
	}

	if filters.EndDate != "" {
		query += fmt.Sprintf(" AND date <= $%d", argPos)
		args = append(args, filters.EndDate)
		argPos++
	}

	if filters.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argPos)
		args = append(args, filters.CategoryID)
		argPos++
	}

	if filters.AccountID != "" {
		query += fmt.Sprintf(" AND account_id = $%d", argPos)
		args = append(args, filters.AccountID)
		argPos++
	}

	if len(filters.Tags) > 0 {
		query += fmt.Sprintf(" AND tags && $%d", argPos)
		args = append(args, filters.Tags)
		argPos++
	}

	query += " ORDER BY date DESC, created_at DESC LIMIT 100"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		t := &model.Transaction{}
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.AccountID, &t.CategoryID,
			&t.Amount, &t.Type, &t.Description, &t.Date,
			&t.IsShared, &t.IsRecurring, &t.RecurringRule,
			&t.Tags, &t.TransferToAccountID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return &model.SearchResult{
		Transactions: transactions,
		TotalCount:   len(transactions),
	}, nil
}

func (r *ReportRepository) GetSpendingByCategory(ctx context.Context, userID string, month, year int) ([]*model.CategorySpending, error) {
	query := `
		SELECT
			t.category_id,
			c.name AS category_name,
			COALESCE(SUM(ABS(t.amount)), 0) AS total_amount
		FROM transactions t
		JOIN categories c ON c.id = t.category_id
		WHERE t.user_id = $1
			AND t.amount < 0
			AND EXTRACT(MONTH FROM t.date) = $2
			AND EXTRACT(YEAR FROM t.date) = $3
		GROUP BY t.category_id, c.name
		ORDER BY total_amount DESC
	`

	rows, err := r.db.Query(ctx, query, userID, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending by category: %w", err)
	}
	defer rows.Close()

	var results []*model.CategorySpending
	var grandTotal int64
	for rows.Next() {
		cs := &model.CategorySpending{}
		if err := rows.Scan(&cs.CategoryID, &cs.CategoryName, &cs.TotalAmount); err != nil {
			return nil, fmt.Errorf("failed to scan category spending: %w", err)
		}
		grandTotal += cs.TotalAmount
		results = append(results, cs)
	}

	// Calculate percentages
	for _, cs := range results {
		if grandTotal > 0 {
			cs.Percentage = float64(cs.TotalAmount) / float64(grandTotal) * 100
		}
	}

	return results, nil
}

func (r *ReportRepository) GetSpendingByMember(ctx context.Context, month, year int) ([]*model.MemberSpending, error) {
	query := `
		SELECT
			t.user_id,
			u.name AS user_name,
			COALESCE(SUM(CASE WHEN t.amount < 0 THEN ABS(t.amount) ELSE 0 END), 0) AS total_expense,
			COALESCE(SUM(CASE WHEN t.amount > 0 THEN t.amount ELSE 0 END), 0) AS total_income
		FROM transactions t
		JOIN users u ON u.id = t.user_id
		WHERE EXTRACT(MONTH FROM t.date) = $1
			AND EXTRACT(YEAR FROM t.date) = $2
		GROUP BY t.user_id, u.name
		ORDER BY total_expense DESC
	`

	rows, err := r.db.Query(ctx, query, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get spending by member: %w", err)
	}
	defer rows.Close()

	var results []*model.MemberSpending
	for rows.Next() {
		ms := &model.MemberSpending{}
		if err := rows.Scan(&ms.UserID, &ms.UserName, &ms.TotalExpense, &ms.TotalIncome); err != nil {
			return nil, fmt.Errorf("failed to scan member spending: %w", err)
		}
		ms.Net = ms.TotalIncome - ms.TotalExpense
		results = append(results, ms)
	}

	return results, nil
}

func (r *ReportRepository) SearchTransactions(ctx context.Context, filters *model.SearchFilters) (*model.SearchResult, error) {
	query := `
		SELECT id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
	`

	args := []any{filters.UserID}
	argPos := 2

	if filters.Description != "" {
		query += fmt.Sprintf(" AND description ILIKE $%d", argPos)
		args = append(args, "%"+filters.Description+"%")
		argPos++
	}

	if filters.MinAmount != nil {
		query += fmt.Sprintf(" AND ABS(amount) >= $%d", argPos)
		args = append(args, *filters.MinAmount)
		argPos++
	}

	if filters.MaxAmount != nil {
		query += fmt.Sprintf(" AND ABS(amount) <= $%d", argPos)
		args = append(args, *filters.MaxAmount)
		argPos++
	}

	if filters.StartDate != "" {
		query += fmt.Sprintf(" AND date >= $%d", argPos)
		args = append(args, filters.StartDate)
		argPos++
	}

	if filters.EndDate != "" {
		query += fmt.Sprintf(" AND date <= $%d", argPos)
		args = append(args, filters.EndDate)
		argPos++
	}

	if filters.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argPos)
		args = append(args, filters.CategoryID)
		argPos++
	}

	if filters.AccountID != "" {
		query += fmt.Sprintf(" AND account_id = $%d", argPos)
		args = append(args, filters.AccountID)
		argPos++
	}

	if len(filters.Tags) > 0 {
		query += fmt.Sprintf(" AND tags && $%d", argPos)
		args = append(args, filters.Tags)
		argPos++
	}

	query += " ORDER BY date DESC, created_at DESC LIMIT 100"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		t := &model.Transaction{}
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.AccountID, &t.CategoryID,
			&t.Amount, &t.Type, &t.Description, &t.Date,
			&t.IsShared, &t.IsRecurring, &t.RecurringRule,
			&t.Tags, &t.TransferToAccountID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return &model.SearchResult{
		Transactions: transactions,
		TotalCount:   len(transactions),
	}, nil
}

func (r *ReportRepository) GetTrends(ctx context.Context, userID string, months int) ([]*model.TrendPoint, error) {
	query := `
		SELECT
			EXTRACT(MONTH FROM date)::int AS month,
			EXTRACT(YEAR FROM date)::int AS year,
			COALESCE(SUM(CASE WHEN amount > 0 THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN amount < 0 THEN ABS(amount) ELSE 0 END), 0) AS total_expense
		FROM transactions
		WHERE user_id = $1
			AND date >= (CURRENT_DATE - make_interval(months => $2))
		GROUP BY EXTRACT(YEAR FROM date), EXTRACT(MONTH FROM date)
		ORDER BY year, month
	`

	rows, err := r.db.Query(ctx, query, userID, months)
	if err != nil {
		return nil, fmt.Errorf("failed to get trends: %w", err)
	}
	defer rows.Close()

	var trends []*model.TrendPoint
	for rows.Next() {
		tp := &model.TrendPoint{}
		if err := rows.Scan(&tp.Month, &tp.Year, &tp.TotalIncome, &tp.TotalExpense); err != nil {
			return nil, fmt.Errorf("failed to scan trend point: %w", err)
		}
		tp.Net = tp.TotalIncome - tp.TotalExpense
		trends = append(trends, tp)
	}

	return trends, nil
}
