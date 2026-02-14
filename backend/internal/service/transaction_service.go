package service

import (
	"context"

	"github.com/yourusername/fambudg/backend/internal/model"
	"github.com/yourusername/fambudg/backend/internal/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (s *TransactionService) Create(ctx context.Context, userID string, req *model.CreateTransactionRequest) (*model.Transaction, error) {
	// Create transaction
	transaction, err := s.transactionRepo.Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// Update account balance
	if err := s.accountRepo.UpdateBalance(ctx, req.AccountID, req.Amount); err != nil {
		return nil, err
	}

	// If it's a transfer, update the destination account as well
	if req.Type == "transfer" && req.TransferToAccountID != nil {
		if err := s.accountRepo.UpdateBalance(ctx, *req.TransferToAccountID, -req.Amount); err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

func (s *TransactionService) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	return s.transactionRepo.FindByID(ctx, id)
}

func (s *TransactionService) GetByUserID(ctx context.Context, userID string, filters *model.TransactionFilters) ([]*model.Transaction, error) {
	return s.transactionRepo.FindByUserID(ctx, userID, filters)
}

func (s *TransactionService) Update(ctx context.Context, id string, req *model.UpdateTransactionRequest) (*model.Transaction, error) {
	// Get original transaction to calculate balance difference
	original, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update transaction
	updated, err := s.transactionRepo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	// If amount or account changed, adjust balances
	if req.Amount != nil || req.AccountID != nil {
		// Reverse original balance change
		if err := s.accountRepo.UpdateBalance(ctx, original.AccountID, -original.Amount); err != nil {
			return nil, err
		}

		// Apply new balance change
		newAccountID := original.AccountID
		if req.AccountID != nil {
			newAccountID = *req.AccountID
		}

		newAmount := original.Amount
		if req.Amount != nil {
			newAmount = *req.Amount
		}

		if err := s.accountRepo.UpdateBalance(ctx, newAccountID, newAmount); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *TransactionService) Delete(ctx context.Context, id string) error {
	// Get transaction to reverse balance change
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Reverse balance change
	if err := s.accountRepo.UpdateBalance(ctx, transaction.AccountID, -transaction.Amount); err != nil {
		return err
	}

	// If it's a transfer, reverse the destination account balance too
	if transaction.Type == "transfer" && transaction.TransferToAccountID != nil {
		if err := s.accountRepo.UpdateBalance(ctx, *transaction.TransferToAccountID, transaction.Amount); err != nil {
			return err
		}
	}

	// Delete transaction
	return s.transactionRepo.Delete(ctx, id)
}
