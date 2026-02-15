package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/asilingas/fambudg/backend/internal/model"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, req *model.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{}
	query := `
		INSERT INTO categories (parent_id, name, type, icon, sort_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, parent_id, name, type, icon, sort_order
	`

	err := r.db.QueryRow(ctx, query, req.ParentID, req.Name, req.Type, req.Icon, req.SortOrder).
		Scan(&category.ID, &category.ParentID, &category.Name, &category.Type, &category.Icon, &category.SortOrder)

	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	query := `
		SELECT id, parent_id, name, type, icon, sort_order
		FROM categories
		ORDER BY sort_order, name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find categories: %w", err)
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		category := &model.Category{}
		if err := rows.Scan(&category.ID, &category.ParentID, &category.Name, &category.Type, &category.Icon, &category.SortOrder); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id string) (*model.Category, error) {
	category := &model.Category{}
	query := `
		SELECT id, parent_id, name, type, icon, sort_order
		FROM categories
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).
		Scan(&category.ID, &category.ParentID, &category.Name, &category.Type, &category.Icon, &category.SortOrder)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id string, req *model.UpdateCategoryRequest) (*model.Category, error) {
	category := &model.Category{}
	query := `
		UPDATE categories
		SET name = COALESCE(NULLIF($1, ''), name),
		    icon = COALESCE(NULLIF($2, ''), icon),
		    sort_order = COALESCE($3, sort_order)
		WHERE id = $4
		RETURNING id, parent_id, name, type, icon, sort_order
	`

	err := r.db.QueryRow(ctx, query, req.Name, req.Icon, req.SortOrder, id).
		Scan(&category.ID, &category.ParentID, &category.Name, &category.Type, &category.Icon, &category.SortOrder)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
