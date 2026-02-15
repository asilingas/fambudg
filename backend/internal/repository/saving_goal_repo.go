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

type SavingGoalRepository struct {
	db *pgxpool.Pool
}

func NewSavingGoalRepository(db *pgxpool.Pool) *SavingGoalRepository {
	return &SavingGoalRepository{db: db}
}

func (r *SavingGoalRepository) Create(ctx context.Context, req *model.CreateSavingGoalRequest) (*model.SavingGoal, error) {
	goal := &model.SavingGoal{}

	var targetDate *time.Time
	if req.TargetDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return nil, fmt.Errorf("invalid target_date format: %w", err)
		}
		targetDate = &parsed
	}

	query := `
		INSERT INTO saving_goals (name, target_amount, target_date, priority)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, target_amount, current_amount, target_date, priority, status, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		req.Name, req.TargetAmount, targetDate, req.Priority,
	).Scan(
		&goal.ID, &goal.Name, &goal.TargetAmount, &goal.CurrentAmount,
		&goal.TargetDate, &goal.Priority, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create saving goal: %w", err)
	}

	return goal, nil
}

func (r *SavingGoalRepository) FindByID(ctx context.Context, id string) (*model.SavingGoal, error) {
	goal := &model.SavingGoal{}
	query := `
		SELECT id, name, target_amount, current_amount, target_date, priority, status, created_at, updated_at
		FROM saving_goals
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&goal.ID, &goal.Name, &goal.TargetAmount, &goal.CurrentAmount,
		&goal.TargetDate, &goal.Priority, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("saving goal not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find saving goal: %w", err)
	}

	return goal, nil
}

func (r *SavingGoalRepository) FindAll(ctx context.Context) ([]*model.SavingGoal, error) {
	query := `
		SELECT id, name, target_amount, current_amount, target_date, priority, status, created_at, updated_at
		FROM saving_goals
		ORDER BY priority, name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find saving goals: %w", err)
	}
	defer rows.Close()

	var goals []*model.SavingGoal
	for rows.Next() {
		goal := &model.SavingGoal{}
		if err := rows.Scan(
			&goal.ID, &goal.Name, &goal.TargetAmount, &goal.CurrentAmount,
			&goal.TargetDate, &goal.Priority, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan saving goal: %w", err)
		}
		goals = append(goals, goal)
	}

	return goals, nil
}

func (r *SavingGoalRepository) Update(ctx context.Context, id string, req *model.UpdateSavingGoalRequest) (*model.SavingGoal, error) {
	updates := []string{}
	args := []any{}
	argPos := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}

	if req.TargetAmount != nil {
		updates = append(updates, fmt.Sprintf("target_amount = $%d", argPos))
		args = append(args, *req.TargetAmount)
		argPos++
	}

	if req.TargetDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.TargetDate)
		if err != nil {
			return nil, fmt.Errorf("invalid target_date format: %w", err)
		}
		updates = append(updates, fmt.Sprintf("target_date = $%d", argPos))
		args = append(args, parsed)
		argPos++
	}

	if req.Priority != nil {
		updates = append(updates, fmt.Sprintf("priority = $%d", argPos))
		args = append(args, *req.Priority)
		argPos++
	}

	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *req.Status)
		argPos++
	}

	updates = append(updates, "updated_at = NOW()")

	if len(updates) == 1 {
		return r.FindByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE saving_goals
		SET %s
		WHERE id = $%d
		RETURNING id, name, target_amount, current_amount, target_date, priority, status, created_at, updated_at
	`, strings.Join(updates, ", "), argPos)

	goal := &model.SavingGoal{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&goal.ID, &goal.Name, &goal.TargetAmount, &goal.CurrentAmount,
		&goal.TargetDate, &goal.Priority, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("saving goal not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update saving goal: %w", err)
	}

	return goal, nil
}

func (r *SavingGoalRepository) Contribute(ctx context.Context, id string, amount int64) (*model.SavingGoal, error) {
	goal := &model.SavingGoal{}
	query := `
		UPDATE saving_goals
		SET current_amount = current_amount + $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, target_amount, current_amount, target_date, priority, status, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, amount, id).Scan(
		&goal.ID, &goal.Name, &goal.TargetAmount, &goal.CurrentAmount,
		&goal.TargetDate, &goal.Priority, &goal.Status, &goal.CreatedAt, &goal.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("saving goal not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to contribute to saving goal: %w", err)
	}

	return goal, nil
}
