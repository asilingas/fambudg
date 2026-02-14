package service

import (
	"testing"

	"github.com/yourusername/fambudg/backend/internal/model"
)

func TestGenerateToken(t *testing.T) {
	authService := &AuthService{
		jwtSecret: "test-secret-key",
	}

	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "admin",
	}

	token, err := authService.generateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	claims, err := authService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims["user_id"] != user.ID {
		t.Errorf("Expected user_id %s, got %v", user.ID, claims["user_id"])
	}

	if claims["email"] != user.Email {
		t.Errorf("Expected email %s, got %v", user.Email, claims["email"])
	}

	if claims["role"] != user.Role {
		t.Errorf("Expected role %s, got %v", user.Role, claims["role"])
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	authService := &AuthService{
		jwtSecret: "test-secret-key",
	}

	_, err := authService.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	authService1 := &AuthService{
		jwtSecret: "secret-1",
	}

	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  "admin",
	}

	token, err := authService1.generateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	authService2 := &AuthService{
		jwtSecret: "secret-2",
	}

	_, err = authService2.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error when validating token with wrong secret, got nil")
	}
}
