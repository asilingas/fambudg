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

type BillReminderRepository struct {
	db *pgxpool.Pool
}

func NewBillReminderRepository(db *pgxpool.Pool) *BillReminderRepository {
	return &BillReminderRepository{db: db}
}

func (r *BillReminderRepository) Create(ctx context.Context, req *model.CreateBillReminderRequest) (*model.BillReminder, error) {
	bill := &model.BillReminder{}

	nextDueDate, err := time.Parse("2006-01-02", req.NextDueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid next_due_date format: %w", err)
	}

	query := `
		WITH inserted AS (
			INSERT INTO bill_reminders (name, amount, due_day, frequency, category_id, account_id, next_due_date)
			VALUES ($1, $2, $3, $4,
				(SELECT id FROM categories WHERE uuid = $5),
				(SELECT id FROM accounts WHERE uuid = $6),
				$7)
			RETURNING *
		)
		SELECT i.uuid, i.name, i.amount, i.due_day, i.frequency,
			c.uuid, a.uuid, i.is_active, i.next_due_date, i.created_at, i.updated_at
		FROM inserted i
		LEFT JOIN categories c ON c.id = i.category_id
		LEFT JOIN accounts a ON a.id = i.account_id
	`

	err = r.db.QueryRow(ctx, query,
		req.Name, req.Amount, req.DueDay, req.Frequency, req.CategoryID, req.AccountID, nextDueDate,
	).Scan(
		&bill.ID, &bill.Name, &bill.Amount, &bill.DueDay, &bill.Frequency,
		&bill.CategoryID, &bill.AccountID, &bill.IsActive, &bill.NextDueDate,
		&bill.CreatedAt, &bill.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create bill reminder: %w", err)
	}

	return bill, nil
}

func (r *BillReminderRepository) FindByID(ctx context.Context, id string) (*model.BillReminder, error) {
	bill := &model.BillReminder{}
	query := `
		SELECT br.uuid, br.name, br.amount, br.due_day, br.frequency,
			c.uuid, a.uuid, br.is_active, br.next_due_date, br.created_at, br.updated_at
		FROM bill_reminders br
		LEFT JOIN categories c ON c.id = br.category_id
		LEFT JOIN accounts a ON a.id = br.account_id
		WHERE br.uuid = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&bill.ID, &bill.Name, &bill.Amount, &bill.DueDay, &bill.Frequency,
		&bill.CategoryID, &bill.AccountID, &bill.IsActive, &bill.NextDueDate,
		&bill.CreatedAt, &bill.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bill reminder not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find bill reminder: %w", err)
	}

	return bill, nil
}

func (r *BillReminderRepository) FindAll(ctx context.Context) ([]*model.BillReminder, error) {
	query := `
		SELECT br.uuid, br.name, br.amount, br.due_day, br.frequency,
			c.uuid, a.uuid, br.is_active, br.next_due_date, br.created_at, br.updated_at
		FROM bill_reminders br
		LEFT JOIN categories c ON c.id = br.category_id
		LEFT JOIN accounts a ON a.id = br.account_id
		ORDER BY br.next_due_date ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find bill reminders: %w", err)
	}
	defer rows.Close()

	var bills []*model.BillReminder
	for rows.Next() {
		bill := &model.BillReminder{}
		if err := rows.Scan(
			&bill.ID, &bill.Name, &bill.Amount, &bill.DueDay, &bill.Frequency,
			&bill.CategoryID, &bill.AccountID, &bill.IsActive, &bill.NextDueDate,
			&bill.CreatedAt, &bill.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan bill reminder: %w", err)
		}
		bills = append(bills, bill)
	}

	return bills, nil
}

func (r *BillReminderRepository) FindUpcoming(ctx context.Context, days int) ([]*model.BillReminder, error) {
	query := `
		SELECT br.uuid, br.name, br.amount, br.due_day, br.frequency,
			c.uuid, a.uuid, br.is_active, br.next_due_date, br.created_at, br.updated_at
		FROM bill_reminders br
		LEFT JOIN categories c ON c.id = br.category_id
		LEFT JOIN accounts a ON a.id = br.account_id
		WHERE br.is_active = true AND br.next_due_date <= CURRENT_DATE + make_interval(days => $1)
		ORDER BY br.next_due_date ASC
	`

	rows, err := r.db.Query(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to find upcoming bill reminders: %w", err)
	}
	defer rows.Close()

	var bills []*model.BillReminder
	for rows.Next() {
		bill := &model.BillReminder{}
		if err := rows.Scan(
			&bill.ID, &bill.Name, &bill.Amount, &bill.DueDay, &bill.Frequency,
			&bill.CategoryID, &bill.AccountID, &bill.IsActive, &bill.NextDueDate,
			&bill.CreatedAt, &bill.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan upcoming bill reminder: %w", err)
		}
		bills = append(bills, bill)
	}

	return bills, nil
}

func (r *BillReminderRepository) Update(ctx context.Context, id string, req *model.UpdateBillReminderRequest) (*model.BillReminder, error) {
	updates := []string{}
	args := []any{}
	argPos := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}

	if req.Amount != nil {
		updates = append(updates, fmt.Sprintf("amount = $%d", argPos))
		args = append(args, *req.Amount)
		argPos++
	}

	if req.DueDay != nil {
		updates = append(updates, fmt.Sprintf("due_day = $%d", argPos))
		args = append(args, *req.DueDay)
		argPos++
	}

	if req.Frequency != nil {
		updates = append(updates, fmt.Sprintf("frequency = $%d", argPos))
		args = append(args, *req.Frequency)
		argPos++
	}

	if req.CategoryID != nil {
		updates = append(updates, fmt.Sprintf("category_id = (SELECT id FROM categories WHERE uuid = $%d)", argPos))
		args = append(args, *req.CategoryID)
		argPos++
	}

	if req.AccountID != nil {
		updates = append(updates, fmt.Sprintf("account_id = (SELECT id FROM accounts WHERE uuid = $%d)", argPos))
		args = append(args, *req.AccountID)
		argPos++
	}

	if req.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *req.IsActive)
		argPos++
	}

	updates = append(updates, "updated_at = NOW()")

	if len(updates) == 1 {
		return r.FindByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		WITH updated AS (
			UPDATE bill_reminders
			SET %s
			WHERE uuid = $%d
			RETURNING *
		)
		SELECT up.uuid, up.name, up.amount, up.due_day, up.frequency,
			c.uuid, a.uuid, up.is_active, up.next_due_date, up.created_at, up.updated_at
		FROM updated up
		LEFT JOIN categories c ON c.id = up.category_id
		LEFT JOIN accounts a ON a.id = up.account_id
	`, strings.Join(updates, ", "), argPos)

	bill := &model.BillReminder{}
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&bill.ID, &bill.Name, &bill.Amount, &bill.DueDay, &bill.Frequency,
		&bill.CategoryID, &bill.AccountID, &bill.IsActive, &bill.NextDueDate,
		&bill.CreatedAt, &bill.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bill reminder not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update bill reminder: %w", err)
	}

	return bill, nil
}

func (r *BillReminderRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM bill_reminders WHERE uuid = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bill reminder: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("bill reminder not found")
	}

	return nil
}

func (r *BillReminderRepository) AdvanceNextDueDate(ctx context.Context, id string, frequency string, currentDueDate time.Time) error {
	var nextDate time.Time
	switch frequency {
	case "monthly":
		nextDate = currentDueDate.AddDate(0, 1, 0)
	case "quarterly":
		nextDate = currentDueDate.AddDate(0, 3, 0)
	case "yearly":
		nextDate = currentDueDate.AddDate(1, 0, 0)
	default:
		nextDate = currentDueDate.AddDate(0, 1, 0)
	}

	query := `UPDATE bill_reminders SET next_due_date = $1, updated_at = NOW() WHERE uuid = $2`
	_, err := r.db.Exec(ctx, query, nextDate, id)
	if err != nil {
		return fmt.Errorf("failed to advance next due date: %w", err)
	}

	return nil
}
