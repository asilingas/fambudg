package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/asilingas/fambudg/backend/internal/model"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, userID string, req *model.CreateTransactionRequest) (*model.Transaction, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	transaction := &model.Transaction{}
	query := `
		INSERT INTO transactions (user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query,
		userID, req.AccountID, req.CategoryID, req.Amount, req.Type, req.Description,
		date, req.IsShared, req.IsRecurring, req.RecurringRule, req.Tags, req.TransferToAccountID,
	).Scan(
		&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID,
		&transaction.Amount, &transaction.Type, &transaction.Description, &transaction.Date,
		&transaction.IsShared, &transaction.IsRecurring, &transaction.RecurringRule,
		&transaction.Tags, &transaction.TransferToAccountID, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, id string) (*model.Transaction, error) {
	transaction := &model.Transaction{}
	query := `
		SELECT id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID,
		&transaction.Amount, &transaction.Type, &transaction.Description, &transaction.Date,
		&transaction.IsShared, &transaction.IsRecurring, &transaction.RecurringRule,
		&transaction.Tags, &transaction.TransferToAccountID, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) FindByUserID(ctx context.Context, userID string, filters *model.TransactionFilters) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argPos := 2

	// Build dynamic filters
	if filters.AccountID != "" {
		query += fmt.Sprintf(" AND account_id = $%d", argPos)
		args = append(args, filters.AccountID)
		argPos++
	}

	if filters.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argPos)
		args = append(args, filters.CategoryID)
		argPos++
	}

	if filters.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argPos)
		args = append(args, filters.Type)
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

	query += " ORDER BY date DESC, created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		transaction := &model.Transaction{}
		if err := rows.Scan(
			&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID,
			&transaction.Amount, &transaction.Type, &transaction.Description, &transaction.Date,
			&transaction.IsShared, &transaction.IsRecurring, &transaction.RecurringRule,
			&transaction.Tags, &transaction.TransferToAccountID, &transaction.CreatedAt, &transaction.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) Update(ctx context.Context, id string, req *model.UpdateTransactionRequest) (*model.Transaction, error) {
	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if req.AccountID != nil {
		updates = append(updates, fmt.Sprintf("account_id = $%d", argPos))
		args = append(args, *req.AccountID)
		argPos++
	}

	if req.CategoryID != nil {
		updates = append(updates, fmt.Sprintf("category_id = $%d", argPos))
		args = append(args, *req.CategoryID)
		argPos++
	}

	if req.Amount != nil {
		updates = append(updates, fmt.Sprintf("amount = $%d", argPos))
		args = append(args, *req.Amount)
		argPos++
	}

	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *req.Description)
		argPos++
	}

	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		updates = append(updates, fmt.Sprintf("date = $%d", argPos))
		args = append(args, date)
		argPos++
	}

	if req.IsShared != nil {
		updates = append(updates, fmt.Sprintf("is_shared = $%d", argPos))
		args = append(args, *req.IsShared)
		argPos++
	}

	if req.Tags != nil {
		updates = append(updates, fmt.Sprintf("tags = $%d", argPos))
		args = append(args, req.Tags)
		argPos++
	}

	updates = append(updates, "updated_at = NOW()")

	if len(updates) == 0 {
		return r.FindByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE transactions
		SET %s
		WHERE id = $%d
		RETURNING id, user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id, created_at, updated_at
	`, strings.Join(updates, ", "), argPos)

	transaction := &model.Transaction{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&transaction.ID, &transaction.UserID, &transaction.AccountID, &transaction.CategoryID,
		&transaction.Amount, &transaction.Type, &transaction.Description, &transaction.Date,
		&transaction.IsShared, &transaction.IsRecurring, &transaction.RecurringRule,
		&transaction.Tags, &transaction.TransferToAccountID, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return transaction, nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
