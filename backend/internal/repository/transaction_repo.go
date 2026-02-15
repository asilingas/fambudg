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

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

const txnSelectCols = `t.uuid, u.uuid, acc.uuid, cat.uuid, t.amount, t.type, t.description, t.date, t.is_shared, t.is_recurring, t.recurring_rule, t.tags, xfer.uuid, t.created_at, t.updated_at`

const txnJoins = `
	FROM transactions t
	JOIN users u ON u.id = t.user_id
	JOIN accounts acc ON acc.id = t.account_id
	JOIN categories cat ON cat.id = t.category_id
	LEFT JOIN accounts xfer ON xfer.id = t.transfer_to_account_id`

func scanTransaction(row interface{ Scan(dest ...any) error }) (*model.Transaction, error) {
	t := &model.Transaction{}
	err := row.Scan(
		&t.ID, &t.UserID, &t.AccountID, &t.CategoryID,
		&t.Amount, &t.Type, &t.Description, &t.Date,
		&t.IsShared, &t.IsRecurring, &t.RecurringRule,
		&t.Tags, &t.TransferToAccountID, &t.CreatedAt, &t.UpdatedAt,
	)
	return t, err
}

func (r *TransactionRepository) Create(ctx context.Context, userID string, req *model.CreateTransactionRequest) (*model.Transaction, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	query := `
		WITH inserted AS (
			INSERT INTO transactions (user_id, account_id, category_id, amount, type, description, date, is_shared, is_recurring, recurring_rule, tags, transfer_to_account_id)
			VALUES (
				(SELECT id FROM users WHERE uuid = $1),
				(SELECT id FROM accounts WHERE uuid = $2),
				(SELECT id FROM categories WHERE uuid = $3),
				$4, $5, $6, $7, $8, $9, $10, $11,
				(SELECT id FROM accounts WHERE uuid = $12)
			)
			RETURNING *
		)
		SELECT i.uuid, u.uuid, acc.uuid, cat.uuid, i.amount, i.type, i.description, i.date, i.is_shared, i.is_recurring, i.recurring_rule, i.tags, xfer.uuid, i.created_at, i.updated_at
		FROM inserted i
		JOIN users u ON u.id = i.user_id
		JOIN accounts acc ON acc.id = i.account_id
		JOIN categories cat ON cat.id = i.category_id
		LEFT JOIN accounts xfer ON xfer.id = i.transfer_to_account_id
	`

	t, err := scanTransaction(r.db.QueryRow(ctx, query,
		userID, req.AccountID, req.CategoryID, req.Amount, req.Type, req.Description,
		date, req.IsShared, req.IsRecurring, req.RecurringRule, req.Tags, req.TransferToAccountID,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return t, nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, id string) (*model.Transaction, error) {
	query := `SELECT ` + txnSelectCols + txnJoins + ` WHERE t.uuid = $1`

	t, err := scanTransaction(r.db.QueryRow(ctx, query, id))
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return t, nil
}

func (r *TransactionRepository) FindByUserID(ctx context.Context, userID string, filters *model.TransactionFilters) ([]*model.Transaction, error) {
	query := `SELECT ` + txnSelectCols + txnJoins + `
		WHERE t.user_id = (SELECT id FROM users WHERE uuid = $1)
	`

	args := []any{userID}
	argPos := 2

	if filters.AccountID != "" {
		query += fmt.Sprintf(" AND t.account_id = (SELECT id FROM accounts WHERE uuid = $%d)", argPos)
		args = append(args, filters.AccountID)
		argPos++
	}

	if filters.CategoryID != "" {
		query += fmt.Sprintf(" AND t.category_id = (SELECT id FROM categories WHERE uuid = $%d)", argPos)
		args = append(args, filters.CategoryID)
		argPos++
	}

	if filters.Type != "" {
		query += fmt.Sprintf(" AND t.type = $%d", argPos)
		args = append(args, filters.Type)
		argPos++
	}

	if filters.StartDate != "" {
		query += fmt.Sprintf(" AND t.date >= $%d", argPos)
		args = append(args, filters.StartDate)
		argPos++
	}

	if filters.EndDate != "" {
		query += fmt.Sprintf(" AND t.date <= $%d", argPos)
		args = append(args, filters.EndDate)
		argPos++
	}

	if filters.IsShared != nil {
		query += fmt.Sprintf(" AND t.is_shared = $%d", argPos)
		args = append(args, *filters.IsShared)
		argPos++
	}

	query += " ORDER BY t.date DESC, t.created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// FindAll returns all transactions with optional filters (for admin)
func (r *TransactionRepository) FindAll(ctx context.Context, filters *model.TransactionFilters) ([]*model.Transaction, error) {
	query := `SELECT ` + txnSelectCols + txnJoins + ` WHERE 1=1`

	args := []any{}
	argPos := 1

	if filters.AccountID != "" {
		query += fmt.Sprintf(" AND t.account_id = (SELECT id FROM accounts WHERE uuid = $%d)", argPos)
		args = append(args, filters.AccountID)
		argPos++
	}

	if filters.CategoryID != "" {
		query += fmt.Sprintf(" AND t.category_id = (SELECT id FROM categories WHERE uuid = $%d)", argPos)
		args = append(args, filters.CategoryID)
		argPos++
	}

	if filters.Type != "" {
		query += fmt.Sprintf(" AND t.type = $%d", argPos)
		args = append(args, filters.Type)
		argPos++
	}

	if filters.StartDate != "" {
		query += fmt.Sprintf(" AND t.date >= $%d", argPos)
		args = append(args, filters.StartDate)
		argPos++
	}

	if filters.EndDate != "" {
		query += fmt.Sprintf(" AND t.date <= $%d", argPos)
		args = append(args, filters.EndDate)
		argPos++
	}

	if filters.IsShared != nil {
		query += fmt.Sprintf(" AND t.is_shared = $%d", argPos)
		args = append(args, *filters.IsShared)
		argPos++
	}

	query += " ORDER BY t.date DESC, t.created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionRepository) Update(ctx context.Context, id string, req *model.UpdateTransactionRequest) (*model.Transaction, error) {
	// Build dynamic update query
	updates := []string{}
	args := []any{}
	argPos := 1

	if req.AccountID != nil {
		updates = append(updates, fmt.Sprintf("account_id = (SELECT id FROM accounts WHERE uuid = $%d)", argPos))
		args = append(args, *req.AccountID)
		argPos++
	}

	if req.CategoryID != nil {
		updates = append(updates, fmt.Sprintf("category_id = (SELECT id FROM categories WHERE uuid = $%d)", argPos))
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
		WITH updated AS (
			UPDATE transactions
			SET %s
			WHERE uuid = $%d
			RETURNING *
		)
		SELECT up.uuid, u.uuid, acc.uuid, cat.uuid, up.amount, up.type, up.description, up.date, up.is_shared, up.is_recurring, up.recurring_rule, up.tags, xfer.uuid, up.created_at, up.updated_at
		FROM updated up
		JOIN users u ON u.id = up.user_id
		JOIN accounts acc ON acc.id = up.account_id
		JOIN categories cat ON cat.id = up.category_id
		LEFT JOIN accounts xfer ON xfer.id = up.transfer_to_account_id
	`, strings.Join(updates, ", "), argPos)

	t, err := scanTransaction(r.db.QueryRow(ctx, query, args...))
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return t, nil
}

func (r *TransactionRepository) FindRecurring(ctx context.Context, userID string) ([]*model.Transaction, error) {
	query := `SELECT ` + txnSelectCols + txnJoins + `
		WHERE t.user_id = (SELECT id FROM users WHERE uuid = $1) AND t.is_recurring = true AND t.recurring_rule IS NOT NULL
		ORDER BY t.date DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find recurring transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		t, err := scanTransaction(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recurring transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionRepository) FindLatestByTemplate(ctx context.Context, userID, accountID, categoryID, description string) (*model.Transaction, error) {
	query := `SELECT ` + txnSelectCols + txnJoins + `
		WHERE t.user_id = (SELECT id FROM users WHERE uuid = $1)
			AND t.account_id = (SELECT id FROM accounts WHERE uuid = $2)
			AND t.category_id = (SELECT id FROM categories WHERE uuid = $3)
			AND t.description = $4 AND t.is_recurring = false
		ORDER BY t.date DESC
		LIMIT 1
	`

	t, err := scanTransaction(r.db.QueryRow(ctx, query, userID, accountID, categoryID, description))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find latest transaction by template: %w", err)
	}

	return t, nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE uuid = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}
