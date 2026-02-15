package service

import (
	"context"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
)

type BudgetService struct {
	budgetRepo *repository.BudgetRepository
}

func NewBudgetService(budgetRepo *repository.BudgetRepository) *BudgetService {
	return &BudgetService{budgetRepo: budgetRepo}
}

func (s *BudgetService) Create(ctx context.Context, req *model.CreateBudgetRequest) (*model.Budget, error) {
	return s.budgetRepo.Create(ctx, req)
}

func (s *BudgetService) GetByID(ctx context.Context, id string) (*model.Budget, error) {
	return s.budgetRepo.FindByID(ctx, id)
}

func (s *BudgetService) GetAll(ctx context.Context, filters *model.BudgetFilters) ([]*model.Budget, error) {
	return s.budgetRepo.FindAll(ctx, filters)
}

func (s *BudgetService) Update(ctx context.Context, id string, req *model.UpdateBudgetRequest) (*model.Budget, error) {
	return s.budgetRepo.Update(ctx, id, req)
}

func (s *BudgetService) Delete(ctx context.Context, id string) error {
	return s.budgetRepo.Delete(ctx, id)
}

func (s *BudgetService) GetSummary(ctx context.Context, month, year int) ([]*model.BudgetSummary, error) {
	return s.budgetRepo.GetSummary(ctx, month, year)
}
