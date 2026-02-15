package service

import (
	"context"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/asilingas/fambudg/backend/internal/repository"
)

type BillReminderService struct {
	billReminderRepo *repository.BillReminderRepository
	transactionRepo  *repository.TransactionRepository
	accountRepo      *repository.AccountRepository
}

func NewBillReminderService(
	billReminderRepo *repository.BillReminderRepository,
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.AccountRepository,
) *BillReminderService {
	return &BillReminderService{
		billReminderRepo: billReminderRepo,
		transactionRepo:  transactionRepo,
		accountRepo:      accountRepo,
	}
}

func (s *BillReminderService) Create(ctx context.Context, req *model.CreateBillReminderRequest) (*model.BillReminder, error) {
	return s.billReminderRepo.Create(ctx, req)
}

func (s *BillReminderService) GetByID(ctx context.Context, id string) (*model.BillReminder, error) {
	return s.billReminderRepo.FindByID(ctx, id)
}

func (s *BillReminderService) GetAll(ctx context.Context) ([]*model.BillReminder, error) {
	return s.billReminderRepo.FindAll(ctx)
}

func (s *BillReminderService) GetUpcoming(ctx context.Context, days int) ([]*model.BillReminder, error) {
	return s.billReminderRepo.FindUpcoming(ctx, days)
}

func (s *BillReminderService) Update(ctx context.Context, id string, req *model.UpdateBillReminderRequest) (*model.BillReminder, error) {
	return s.billReminderRepo.Update(ctx, id, req)
}

func (s *BillReminderService) Delete(ctx context.Context, id string) error {
	return s.billReminderRepo.Delete(ctx, id)
}

func (s *BillReminderService) Pay(ctx context.Context, userID string, id string, req *model.PayBillRequest) (*model.Transaction, error) {
	bill, err := s.billReminderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create a transaction for the payment
	txReq := &model.CreateTransactionRequest{
		AccountID:   req.AccountID,
		Amount:      -bill.Amount, // expense
		Type:        "expense",
		Description: "Bill payment: " + bill.Name,
		Date:        req.Date,
		IsShared:    true,
	}

	if bill.CategoryID != nil {
		txReq.CategoryID = *bill.CategoryID
	}

	transaction, err := s.transactionRepo.Create(ctx, userID, txReq)
	if err != nil {
		return nil, err
	}

	// Update account balance
	if err := s.accountRepo.UpdateBalance(ctx, req.AccountID, -bill.Amount); err != nil {
		return nil, err
	}

	// Advance the next due date
	if err := s.billReminderRepo.AdvanceNextDueDate(ctx, id, bill.Frequency, bill.NextDueDate); err != nil {
		return nil, err
	}

	return transaction, nil
}
