package service

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
)

type SavingGoalService struct {
	savingGoalRepo *repository.SavingGoalRepository
}

func NewSavingGoalService(savingGoalRepo *repository.SavingGoalRepository) *SavingGoalService {
	return &SavingGoalService{savingGoalRepo: savingGoalRepo}
}

func (s *SavingGoalService) Create(ctx context.Context, req *model.CreateSavingGoalRequest) (*model.SavingGoal, error) {
	return s.savingGoalRepo.Create(ctx, req)
}

func (s *SavingGoalService) GetByID(ctx context.Context, id string) (*model.SavingGoal, error) {
	return s.savingGoalRepo.FindByID(ctx, id)
}

func (s *SavingGoalService) GetAll(ctx context.Context) ([]*model.SavingGoal, error) {
	return s.savingGoalRepo.FindAll(ctx)
}

func (s *SavingGoalService) Update(ctx context.Context, id string, req *model.UpdateSavingGoalRequest) (*model.SavingGoal, error) {
	return s.savingGoalRepo.Update(ctx, id, req)
}

func (s *SavingGoalService) Contribute(ctx context.Context, id string, req *model.ContributeRequest) (*model.SavingGoal, error) {
	goal, err := s.savingGoalRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if goal.Status != "active" {
		return nil, fmt.Errorf("cannot contribute to a %s goal", goal.Status)
	}

	updated, err := s.savingGoalRepo.Contribute(ctx, id, req.Amount)
	if err != nil {
		return nil, err
	}

	// Auto-complete if target reached
	if updated.CurrentAmount >= updated.TargetAmount {
		completedStatus := "completed"
		updateReq := &model.UpdateSavingGoalRequest{Status: &completedStatus}
		return s.savingGoalRepo.Update(ctx, id, updateReq)
	}

	return updated, nil
}
