package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Transaction struct {
	ID                  string         `json:"id"`
	UserID              string         `json:"userId"`
	AccountID           string         `json:"accountId"`
	CategoryID          string         `json:"categoryId"`
	Amount              int64          `json:"amount"` // in cents, positive=income, negative=expense
	Type                string         `json:"type" validate:"required,oneof=expense income transfer"`
	Description         string         `json:"description,omitempty"`
	Date                time.Time      `json:"date"`
	IsShared            bool           `json:"isShared"`
	IsRecurring         bool           `json:"isRecurring"`
	RecurringRule       *RecurringRule `json:"recurringRule,omitempty"`
	Tags                []string       `json:"tags,omitempty"`
	TransferToAccountID *string        `json:"transferToAccountId,omitempty"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
}

type RecurringRule struct {
	Frequency string `json:"frequency"`           // daily, weekly, monthly, yearly
	Day       int    `json:"day,omitempty"`       // day of month for monthly
	DayOfWeek int    `json:"dayOfWeek,omitempty"` // 0=Sunday, 6=Saturday for weekly
}

// Scan implements sql.Scanner for RecurringRule
func (r *RecurringRule) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, r)
}

// Value implements driver.Valuer for RecurringRule
func (r RecurringRule) Value() (driver.Value, error) {
	return json.Marshal(r)
}

type CreateTransactionRequest struct {
	AccountID           string         `json:"accountId" validate:"required"`
	CategoryID          string         `json:"categoryId" validate:"required"`
	Amount              int64          `json:"amount" validate:"required"`
	Type                string         `json:"type" validate:"required,oneof=expense income transfer"`
	Description         string         `json:"description"`
	Date                string         `json:"date" validate:"required"` // YYYY-MM-DD format
	IsShared            bool           `json:"isShared"`
	IsRecurring         bool           `json:"isRecurring"`
	RecurringRule       *RecurringRule `json:"recurringRule,omitempty"`
	Tags                []string       `json:"tags,omitempty"`
	TransferToAccountID *string        `json:"transferToAccountId,omitempty"`
}

type UpdateTransactionRequest struct {
	AccountID   *string  `json:"accountId,omitempty"`
	CategoryID  *string  `json:"categoryId,omitempty"`
	Amount      *int64   `json:"amount,omitempty"`
	Description *string  `json:"description,omitempty"`
	Date        *string  `json:"date,omitempty"` // YYYY-MM-DD format
	IsShared    *bool    `json:"isShared,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type GenerateRecurringResponse struct {
	Generated int      `json:"generated"`
	Templates int      `json:"templates"`
	Errors    []string `json:"errors,omitempty"`
}

type TransactionFilters struct {
	AccountID  string
	CategoryID string
	StartDate  string // YYYY-MM-DD
	EndDate    string // YYYY-MM-DD
	Type       string
	IsShared   *bool
}
