package service

import (
	"context"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) Create(ctx context.Context, userID string, req *model.CreateAccountRequest) (*model.Account, error) {
	return s.accountRepo.Create(ctx, userID, req)
}

func (s *AccountService) GetByID(ctx context.Context, id string) (*model.Account, error) {
	return s.accountRepo.FindByID(ctx, id)
}

func (s *AccountService) GetByUserID(ctx context.Context, userID string) ([]*model.Account, error) {
	return s.accountRepo.FindByUserID(ctx, userID)
}

func (s *AccountService) GetAll(ctx context.Context) ([]*model.Account, error) {
	return s.accountRepo.FindAll(ctx)
}

func (s *AccountService) Update(ctx context.Context, id string, req *model.UpdateAccountRequest) (*model.Account, error) {
	return s.accountRepo.Update(ctx, id, req)
}

func (s *AccountService) Delete(ctx context.Context, id string) error {
	return s.accountRepo.Delete(ctx, id)
}
