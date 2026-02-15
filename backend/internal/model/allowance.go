package model

import "time"

type Allowance struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Amount      int64     `json:"amount"`
	Spent       int64     `json:"spent"`
	Remaining   int64     `json:"remaining"`
	PeriodStart time.Time `json:"periodStart"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateAllowanceRequest struct {
	UserID      string `json:"userId" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	PeriodStart string `json:"periodStart" validate:"required"` // YYYY-MM-DD
}

type UpdateAllowanceRequest struct {
	Amount      *int64  `json:"amount,omitempty" validate:"omitempty,gt=0"`
	PeriodStart *string `json:"periodStart,omitempty"` // YYYY-MM-DD
}
