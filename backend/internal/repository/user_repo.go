package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new admin user via self-registration
func (r *UserRepository) Create(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	return r.createUser(ctx, req.Email, req.Password, req.Name, "admin")
}

// CreateWithRole creates a new user with the specified role
func (r *UserRepository) CreateWithRole(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	return r.createUser(ctx, req.Email, req.Password, req.Name, req.Role)
}

func (r *UserRepository) createUser(ctx context.Context, email, password, name, role string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{}
	query := `
		INSERT INTO users (email, password_hash, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, name, role, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query, email, string(hashedPassword), name, role).
		Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// ListAll returns all users
func (r *UserRepository) ListAll(ctx context.Context) ([]*model.User, error) {
	query := `
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// Update updates a user's name and/or role
func (r *UserRepository) Update(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error) {
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}

	if req.Role != nil {
		updates = append(updates, fmt.Sprintf("role = $%d", argPos))
		args = append(args, *req.Role)
		argPos++
	}

	if len(updates) == 0 {
		return r.FindByID(ctx, id)
	}

	updates = append(updates, "updated_at = NOW()")

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING id, email, name, role, created_at, updated_at
	`, strings.Join(updates, ", "), argPos)

	user := &model.User{}
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, email, name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// VerifyPassword checks if the provided password matches the user's hashed password
func (r *UserRepository) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
