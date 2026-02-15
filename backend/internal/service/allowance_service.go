package service

import (
	"context"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
)

type AllowanceService struct {
	allowanceRepo *repository.AllowanceRepository
}

func NewAllowanceService(allowanceRepo *repository.AllowanceRepository) *AllowanceService {
	return &AllowanceService{allowanceRepo: allowanceRepo}
}

func (s *AllowanceService) Create(ctx context.Context, req *model.CreateAllowanceRequest) (*model.Allowance, error) {
	allowance, err := s.allowanceRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.enrichWithSpending(ctx, allowance)
}

func (s *AllowanceService) GetByID(ctx context.Context, id string) (*model.Allowance, error) {
	allowance, err := s.allowanceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.enrichWithSpending(ctx, allowance)
}

func (s *AllowanceService) GetByUserID(ctx context.Context, userID string) (*model.Allowance, error) {
	allowance, err := s.allowanceRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.enrichWithSpending(ctx, allowance)
}

func (s *AllowanceService) GetAll(ctx context.Context) ([]*model.Allowance, error) {
	allowances, err := s.allowanceRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	for i, a := range allowances {
		enriched, err := s.enrichWithSpending(ctx, a)
		if err != nil {
			return nil, err
		}
		allowances[i] = enriched
	}

	return allowances, nil
}

func (s *AllowanceService) Update(ctx context.Context, id string, req *model.UpdateAllowanceRequest) (*model.Allowance, error) {
	allowance, err := s.allowanceRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return s.enrichWithSpending(ctx, allowance)
}

// enrichWithSpending calculates spent and remaining from transactions
func (s *AllowanceService) enrichWithSpending(ctx context.Context, allowance *model.Allowance) (*model.Allowance, error) {
	spent, err := s.allowanceRepo.GetSpentInPeriod(ctx, allowance.UserID, allowance.PeriodStart)
	if err != nil {
		return nil, err
	}

	allowance.Spent = spent
	allowance.Remaining = allowance.Amount - spent
	return allowance, nil
}
