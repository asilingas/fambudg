package model

import "time"

type Account struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Type      string    `json:"type" validate:"required,oneof=checking savings credit cash"`
	Currency  string    `json:"currency" validate:"required,len=3"`
	Balance   int64     `json:"balance"` // in cents
	CreatedAt time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Type     string `json:"type" validate:"required,oneof=checking savings credit cash"`
	Currency string `json:"currency" validate:"required,len=3"`
	Balance  int64  `json:"balance"`
}

type UpdateAccountRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Type     string `json:"type" validate:"omitempty,oneof=checking savings credit cash"`
	Currency string `json:"currency" validate:"omitempty,len=3"`
}
