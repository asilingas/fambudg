package model

import "time"

type BillReminder struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Amount      int64      `json:"amount"`
	DueDay      int        `json:"dueDay"`
	Frequency   string     `json:"frequency"`
	CategoryID  *string    `json:"categoryId,omitempty"`
	AccountID   *string    `json:"accountId,omitempty"`
	IsActive    bool       `json:"isActive"`
	NextDueDate time.Time  `json:"nextDueDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type CreateBillReminderRequest struct {
	Name        string  `json:"name" validate:"required,max=200"`
	Amount      int64   `json:"amount" validate:"required,gt=0"`
	DueDay      int     `json:"dueDay" validate:"required,min=1,max=31"`
	Frequency   string  `json:"frequency" validate:"required,oneof=monthly quarterly yearly"`
	CategoryID  *string `json:"categoryId,omitempty"`
	AccountID   *string `json:"accountId,omitempty"`
	NextDueDate string  `json:"nextDueDate" validate:"required"`
}

type UpdateBillReminderRequest struct {
	Name       *string `json:"name,omitempty" validate:"omitempty,max=200"`
	Amount     *int64  `json:"amount,omitempty" validate:"omitempty,gt=0"`
	DueDay     *int    `json:"dueDay,omitempty" validate:"omitempty,min=1,max=31"`
	Frequency  *string `json:"frequency,omitempty" validate:"omitempty,oneof=monthly quarterly yearly"`
	CategoryID *string `json:"categoryId,omitempty"`
	AccountID  *string `json:"accountId,omitempty"`
	IsActive   *bool   `json:"isActive,omitempty"`
}

type PayBillRequest struct {
	AccountID string `json:"accountId" validate:"required"`
	Date      string `json:"date" validate:"required"`
}

type TransferRequest struct {
	FromAccountID string `json:"fromAccountId" validate:"required"`
	ToAccountID   string `json:"toAccountId" validate:"required"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Description   string `json:"description"`
	Date          string `json:"date" validate:"required"`
}
