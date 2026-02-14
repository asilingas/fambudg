package service

import (
	"context"

	"github.com/yourusername/fambudg/backend/internal/model"
	"github.com/yourusername/fambudg/backend/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) Create(ctx context.Context, req *model.CreateCategoryRequest) (*model.Category, error) {
	return s.categoryRepo.Create(ctx, req)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*model.Category, error) {
	return s.categoryRepo.FindAll(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id string) (*model.Category, error) {
	return s.categoryRepo.FindByID(ctx, id)
}

func (s *CategoryService) Update(ctx context.Context, id string, req *model.UpdateCategoryRequest) (*model.Category, error) {
	return s.categoryRepo.Update(ctx, id, req)
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	return s.categoryRepo.Delete(ctx, id)
}
