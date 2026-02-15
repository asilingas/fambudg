package repository

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		WITH inserted AS (
			INSERT INTO categories (parent_id, name, type, icon, sort_order)
			VALUES ((SELECT id FROM categories WHERE uuid = $1), $2, $3, $4, $5)
			RETURNING *
		)
		SELECT i.uuid, p.uuid, i.name, i.type, i.icon, i.sort_order
		FROM inserted i LEFT JOIN categories p ON p.id = i.parent_id
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
		SELECT c.uuid, p.uuid, c.name, c.type, c.icon, c.sort_order
		FROM categories c LEFT JOIN categories p ON p.id = c.parent_id
		ORDER BY c.sort_order, c.name
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
		SELECT c.uuid, p.uuid, c.name, c.type, c.icon, c.sort_order
		FROM categories c LEFT JOIN categories p ON p.id = c.parent_id
		WHERE c.uuid = $1
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
		WITH updated AS (
			UPDATE categories
			SET name = COALESCE(NULLIF($1, ''), name),
			    icon = COALESCE(NULLIF($2, ''), icon),
			    sort_order = COALESCE($3, sort_order)
			WHERE uuid = $4
			RETURNING *
		)
		SELECT up.uuid, p.uuid, up.name, up.type, up.icon, up.sort_order
		FROM updated up LEFT JOIN categories p ON p.id = up.parent_id
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
	query := `DELETE FROM categories WHERE uuid = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
