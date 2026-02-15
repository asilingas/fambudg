package repository

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account
func (r *AccountRepository) Create(ctx context.Context, userID string, req *model.CreateAccountRequest) (*model.Account, error) {
	account := &model.Account{}
	query := `
		INSERT INTO accounts (user_id, name, type, currency, balance)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, type, currency, balance, created_at
	`

	err := r.db.QueryRow(ctx, query, userID, req.Name, req.Type, req.Currency, req.Balance).
		Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Currency, &account.Balance, &account.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// FindByID finds an account by ID
func (r *AccountRepository) FindByID(ctx context.Context, id string) (*model.Account, error) {
	account := &model.Account{}
	query := `
		SELECT id, user_id, name, type, currency, balance, created_at
		FROM accounts
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Currency, &account.Balance, &account.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find account: %w", err)
	}

	return account, nil
}

// FindByUserID finds all accounts for a user
func (r *AccountRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Account, error) {
	query := `
		SELECT id, user_id, name, type, currency, balance, created_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		account := &model.Account{}
		if err := rows.Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Currency, &account.Balance, &account.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// Update updates an account
func (r *AccountRepository) Update(ctx context.Context, id string, req *model.UpdateAccountRequest) (*model.Account, error) {
	query := `
		UPDATE accounts
		SET name = COALESCE(NULLIF($1, ''), name),
		    type = COALESCE(NULLIF($2, ''), type),
		    currency = COALESCE(NULLIF($3, ''), currency)
		WHERE id = $4
		RETURNING id, user_id, name, type, currency, balance, created_at
	`

	account := &model.Account{}
	err := r.db.QueryRow(ctx, query, req.Name, req.Type, req.Currency, id).
		Scan(&account.ID, &account.UserID, &account.Name, &account.Type, &account.Currency, &account.Balance, &account.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}

// Delete deletes an account
func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// UpdateBalance updates an account balance
func (r *AccountRepository) UpdateBalance(ctx context.Context, accountID string, amount int64) error {
	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, amount, accountID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}
