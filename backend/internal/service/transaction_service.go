package service

import (
	"context"
	"fmt"
	"time"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
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

func (s *TransactionService) GenerateRecurring(ctx context.Context, userID string, upTo time.Time) (*model.GenerateRecurringResponse, error) {
	templates, err := s.transactionRepo.FindRecurring(ctx, userID)
	if err != nil {
		return nil, err
	}

	var generated int
	var errors []string

	for _, tmpl := range templates {
		if tmpl.RecurringRule == nil {
			continue
		}

		// Find the latest generated transaction for this template
		latest, err := s.transactionRepo.FindLatestByTemplate(ctx, tmpl.UserID, tmpl.AccountID, tmpl.CategoryID, tmpl.Description)
		if err != nil {
			errors = append(errors, fmt.Sprintf("template %s: %s", tmpl.ID, err.Error()))
			continue
		}

		// Determine start date: day after the latest generated, or the template date
		startDate := tmpl.Date
		if latest != nil {
			startDate = latest.Date
		}

		// Generate occurrences from startDate to upTo
		nextDate := nextOccurrence(startDate, tmpl.RecurringRule)
		for !nextDate.After(upTo) {
			req := &model.CreateTransactionRequest{
				AccountID:   tmpl.AccountID,
				CategoryID:  tmpl.CategoryID,
				Amount:      tmpl.Amount,
				Type:        tmpl.Type,
				Description: tmpl.Description,
				Date:        nextDate.Format("2006-01-02"),
				IsShared:    tmpl.IsShared,
				IsRecurring: false, // generated copies are not recurring
				Tags:        tmpl.Tags,
			}

			_, err := s.Create(ctx, userID, req)
			if err != nil {
				errors = append(errors, fmt.Sprintf("template %s date %s: %s", tmpl.ID, nextDate.Format("2006-01-02"), err.Error()))
				break
			}

			generated++
			nextDate = nextOccurrence(nextDate, tmpl.RecurringRule)
		}
	}

	return &model.GenerateRecurringResponse{
		Generated:  generated,
		Templates:  len(templates),
		Errors:     errors,
	}, nil
}

func nextOccurrence(from time.Time, rule *model.RecurringRule) time.Time {
	switch rule.Frequency {
	case "daily":
		return from.AddDate(0, 0, 1)
	case "weekly":
		return from.AddDate(0, 0, 7)
	case "monthly":
		next := from.AddDate(0, 1, 0)
		if rule.Day > 0 {
			// Set to the specific day of month
			next = time.Date(next.Year(), next.Month(), rule.Day, 0, 0, 0, 0, next.Location())
		}
		return next
	case "yearly":
		return from.AddDate(1, 0, 0)
	default:
		return from.AddDate(0, 1, 0)
	}
}
