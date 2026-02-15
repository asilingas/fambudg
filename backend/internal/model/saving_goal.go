package model

import "time"

type SavingGoal struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	TargetAmount  int64      `json:"targetAmount"`
	CurrentAmount int64      `json:"currentAmount"`
	TargetDate    *time.Time `json:"targetDate,omitempty"`
	Priority      int        `json:"priority"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type CreateSavingGoalRequest struct {
	Name         string  `json:"name" validate:"required,max=200"`
	TargetAmount int64   `json:"targetAmount" validate:"required,gt=0"`
	TargetDate   *string `json:"targetDate,omitempty"`
	Priority     int     `json:"priority" validate:"min=1"`
}

type UpdateSavingGoalRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,max=200"`
	TargetAmount *int64  `json:"targetAmount,omitempty" validate:"omitempty,gt=0"`
	TargetDate   *string `json:"targetDate,omitempty"`
	Priority     *int    `json:"priority,omitempty" validate:"omitempty,min=1"`
	Status       *string `json:"status,omitempty" validate:"omitempty,oneof=active completed cancelled"`
}

type ContributeRequest struct {
	Amount int64 `json:"amount" validate:"required,gt=0"`
}
