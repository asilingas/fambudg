package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AllowanceRepository struct {
	db *pgxpool.Pool
}

func NewAllowanceRepository(db *pgxpool.Pool) *AllowanceRepository {
	return &AllowanceRepository{db: db}
}

func (r *AllowanceRepository) Create(ctx context.Context, req *model.CreateAllowanceRequest) (*model.Allowance, error) {
	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		return nil, fmt.Errorf("invalid period_start date format: %w", err)
	}

	allowance := &model.Allowance{}
	query := `
		INSERT INTO allowances (user_id, amount, period_start)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, amount, period_start, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query, req.UserID, req.Amount, periodStart).
		Scan(&allowance.ID, &allowance.UserID, &allowance.Amount, &allowance.PeriodStart, &allowance.CreatedAt, &allowance.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create allowance: %w", err)
	}

	return allowance, nil
}

func (r *AllowanceRepository) FindByID(ctx context.Context, id string) (*model.Allowance, error) {
	allowance := &model.Allowance{}
	query := `
		SELECT id, user_id, amount, period_start, created_at, updated_at
		FROM allowances
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&allowance.ID, &allowance.UserID, &allowance.Amount, &allowance.PeriodStart, &allowance.CreatedAt, &allowance.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("allowance not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find allowance: %w", err)
	}

	return allowance, nil
}

func (r *AllowanceRepository) FindByUserID(ctx context.Context, userID string) (*model.Allowance, error) {
	allowance := &model.Allowance{}
	query := `
		SELECT id, user_id, amount, period_start, created_at, updated_at
		FROM allowances
		WHERE user_id = $1
	`

	err := r.db.QueryRow(ctx, query, userID).
		Scan(&allowance.ID, &allowance.UserID, &allowance.Amount, &allowance.PeriodStart, &allowance.CreatedAt, &allowance.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("allowance not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find allowance: %w", err)
	}

	return allowance, nil
}

func (r *AllowanceRepository) ListAll(ctx context.Context) ([]*model.Allowance, error) {
	query := `
		SELECT id, user_id, amount, period_start, created_at, updated_at
		FROM allowances
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list allowances: %w", err)
	}
	defer rows.Close()

	var allowances []*model.Allowance
	for rows.Next() {
		a := &model.Allowance{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.Amount, &a.PeriodStart, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan allowance: %w", err)
		}
		allowances = append(allowances, a)
	}

	return allowances, nil
}

func (r *AllowanceRepository) Update(ctx context.Context, id string, req *model.UpdateAllowanceRequest) (*model.Allowance, error) {
	updates := []string{}
	args := []any{}
	argPos := 1

	if req.Amount != nil {
		updates = append(updates, fmt.Sprintf("amount = $%d", argPos))
		args = append(args, *req.Amount)
		argPos++
	}

	if req.PeriodStart != nil {
		periodStart, err := time.Parse("2006-01-02", *req.PeriodStart)
		if err != nil {
			return nil, fmt.Errorf("invalid period_start date format: %w", err)
		}
		updates = append(updates, fmt.Sprintf("period_start = $%d", argPos))
		args = append(args, periodStart)
		argPos++
	}

	if len(updates) == 0 {
		return r.FindByID(ctx, id)
	}

	updates = append(updates, "updated_at = NOW()")

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE allowances
		SET %s
		WHERE id = $%d
		RETURNING id, user_id, amount, period_start, created_at, updated_at
	`, strings.Join(updates, ", "), argPos)

	allowance := &model.Allowance{}
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&allowance.ID, &allowance.UserID, &allowance.Amount, &allowance.PeriodStart, &allowance.CreatedAt, &allowance.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("allowance not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update allowance: %w", err)
	}

	return allowance, nil
}

// GetSpentInPeriod calculates spending for a user from period_start to end of that month
func (r *AllowanceRepository) GetSpentInPeriod(ctx context.Context, userID string, periodStart time.Time) (int64, error) {
	// Calculate period end: one month after period_start
	periodEnd := periodStart.AddDate(0, 1, 0)

	query := `
		SELECT COALESCE(SUM(ABS(amount)), 0)
		FROM transactions
		WHERE user_id = $1
			AND amount < 0
			AND date >= $2
			AND date < $3
	`

	var spent int64
	err := r.db.QueryRow(ctx, query, userID, periodStart, periodEnd).Scan(&spent)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate spent amount: %w", err)
	}

	return spent, nil
}
