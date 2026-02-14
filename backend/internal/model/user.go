package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name" validate:"required,min=2,max=100"`
	Role         string    `json:"role" validate:"required,oneof=admin member child"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
