package model

import "time"

type Budget struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"categoryId"`
	Amount     int64     `json:"amount"` // budget limit in cents
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CreateBudgetRequest struct {
	CategoryID string `json:"categoryId" validate:"required"`
	Amount     int64  `json:"amount" validate:"required,gt=0"`
	Month      int    `json:"month" validate:"required,min=1,max=12"`
	Year       int    `json:"year" validate:"required,min=2000"`
}

type UpdateBudgetRequest struct {
	Amount *int64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
}

type BudgetSummary struct {
	CategoryID   string `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	BudgetAmount int64  `json:"budgetAmount"` // budget limit in cents
	ActualAmount int64  `json:"actualAmount"` // actual spending in cents (positive = spent)
	Remaining    int64  `json:"remaining"`    // budget - actual
}

type BudgetFilters struct {
	Month int
	Year  int
}
