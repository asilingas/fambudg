package repository

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BudgetRepository struct {
	db *pgxpool.Pool
}

func NewBudgetRepository(db *pgxpool.Pool) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Create(ctx context.Context, req *model.CreateBudgetRequest) (*model.Budget, error) {
	budget := &model.Budget{}
	query := `
		WITH inserted AS (
			INSERT INTO budgets (category_id, amount, month, year)
			VALUES ((SELECT id FROM categories WHERE uuid = $1), $2, $3, $4)
			RETURNING *
		)
		SELECT i.uuid, c.uuid, i.amount, i.month, i.year, i.created_at
		FROM inserted i JOIN categories c ON c.id = i.category_id
	`

	err := r.db.QueryRow(ctx, query,
		req.CategoryID, req.Amount, req.Month, req.Year,
	).Scan(
		&budget.ID, &budget.CategoryID, &budget.Amount,
		&budget.Month, &budget.Year, &budget.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return budget, nil
}

func (r *BudgetRepository) FindByID(ctx context.Context, id string) (*model.Budget, error) {
	budget := &model.Budget{}
	query := `
		SELECT b.uuid, c.uuid, b.amount, b.month, b.year, b.created_at
		FROM budgets b JOIN categories c ON c.id = b.category_id
		WHERE b.uuid = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&budget.ID, &budget.CategoryID, &budget.Amount,
		&budget.Month, &budget.Year, &budget.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("budget not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find budget: %w", err)
	}

	return budget, nil
}

func (r *BudgetRepository) FindAll(ctx context.Context, filters *model.BudgetFilters) ([]*model.Budget, error) {
	query := `
		SELECT b.uuid, c.uuid, b.amount, b.month, b.year, b.created_at
		FROM budgets b JOIN categories c ON c.id = b.category_id
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filters.Month > 0 {
		query += fmt.Sprintf(" AND b.month = $%d", argPos)
		args = append(args, filters.Month)
		argPos++
	}

	if filters.Year > 0 {
		query += fmt.Sprintf(" AND b.year = $%d", argPos)
		args = append(args, filters.Year)
		argPos++
	}

	query += " ORDER BY b.year DESC, b.month DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find budgets: %w", err)
	}
	defer rows.Close()

	var budgets []*model.Budget
	for rows.Next() {
		budget := &model.Budget{}
		if err := rows.Scan(
			&budget.ID, &budget.CategoryID, &budget.Amount,
			&budget.Month, &budget.Year, &budget.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (r *BudgetRepository) Update(ctx context.Context, id string, req *model.UpdateBudgetRequest) (*model.Budget, error) {
	if req.Amount == nil {
		return r.FindByID(ctx, id)
	}

	budget := &model.Budget{}
	query := `
		WITH updated AS (
			UPDATE budgets
			SET amount = $1
			WHERE uuid = $2
			RETURNING *
		)
		SELECT up.uuid, c.uuid, up.amount, up.month, up.year, up.created_at
		FROM updated up JOIN categories c ON c.id = up.category_id
	`

	err := r.db.QueryRow(ctx, query, *req.Amount, id).Scan(
		&budget.ID, &budget.CategoryID, &budget.Amount,
		&budget.Month, &budget.Year, &budget.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("budget not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return budget, nil
}

func (r *BudgetRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM budgets WHERE uuid = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("budget not found")
	}

	return nil
}

func (r *BudgetRepository) GetSummary(ctx context.Context, month, year int) ([]*model.BudgetSummary, error) {
	query := `
		SELECT
			c.uuid,
			c.name AS category_name,
			b.amount AS budget_amount,
			COALESCE(ABS(SUM(CASE WHEN t.amount < 0 THEN t.amount ELSE 0 END)), 0) AS actual_amount
		FROM budgets b
		JOIN categories c ON c.id = b.category_id
		LEFT JOIN transactions t ON t.category_id = b.category_id
			AND EXTRACT(MONTH FROM t.date) = $1
			AND EXTRACT(YEAR FROM t.date) = $2
		WHERE b.month = $1 AND b.year = $2
		GROUP BY c.uuid, c.name, b.amount
		ORDER BY c.name
	`

	rows, err := r.db.Query(ctx, query, month, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget summary: %w", err)
	}
	defer rows.Close()

	var summaries []*model.BudgetSummary
	for rows.Next() {
		s := &model.BudgetSummary{}
		if err := rows.Scan(
			&s.CategoryID, &s.CategoryName, &s.BudgetAmount, &s.ActualAmount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan budget summary: %w", err)
		}
		s.Remaining = s.BudgetAmount - s.ActualAmount
		summaries = append(summaries, s)
	}

	return summaries, nil
}
