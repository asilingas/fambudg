package steps

import (
	"context"
	"fmt"
	"time"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerRecurringSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I have a recurring transaction of (-?\d+) on "([^"]*)" with frequency "([^"]*)"$`, tc.iHaveARecurringTransaction)
	ctx.Step(`^I generate recurring transactions up to "([^"]*)"$`, tc.iGenerateRecurringTransactionsUpTo)
	ctx.Step(`^(\d+) recurring transactions should have been generated$`, tc.nRecurringTransactionsShouldBeGenerated)
	ctx.Step(`^the generated transactions should have amounts of (-?\d+)$`, tc.theGeneratedTransactionsShouldHaveAmounts)
}

func (tc *TestContext) iHaveARecurringTransaction(amount int64, dateStr, frequency string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	account, ok := tc.CurrentAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no current account")
	}

	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category")
	}

	rule := &model.RecurringRule{
		Frequency: frequency,
	}

	req := &model.CreateTransactionRequest{
		AccountID:     account.ID,
		CategoryID:    category.ID,
		Amount:        amount,
		Type:          "expense",
		Description:   "Recurring " + frequency,
		Date:          dateStr,
		IsRecurring:   true,
		RecurringRule: rule,
	}

	transaction, err := tc.TransactionService.Create(context.Background(), user.ID, req)
	if err != nil {
		return fmt.Errorf("failed to create recurring transaction: %w", err)
	}

	tc.CurrentTransaction = transaction
	return nil
}

func (tc *TestContext) iGenerateRecurringTransactionsUpTo(dateStr string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	upTo, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	result, err := tc.TransactionService.GenerateRecurring(context.Background(), user.ID, upTo)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.RecurringResult = result
	tc.LastError = nil
	return nil
}

func (tc *TestContext) nRecurringTransactionsShouldBeGenerated(expected int) error {
	if tc.LastError != nil {
		return fmt.Errorf("expected success, got error: %v", tc.LastError)
	}

	result, ok := tc.RecurringResult.(*model.GenerateRecurringResponse)
	if !ok {
		return fmt.Errorf("no recurring result available")
	}

	if result.Generated != expected {
		return fmt.Errorf("expected %d generated, got %d", expected, result.Generated)
	}

	return nil
}

func (tc *TestContext) theGeneratedTransactionsShouldHaveAmounts(expectedAmount int64) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	// Fetch all non-recurring transactions to verify amounts
	transactions, err := tc.TransactionService.GetByUserID(context.Background(), user.ID, &model.TransactionFilters{})
	if err != nil {
		return fmt.Errorf("failed to get transactions: %w", err)
	}

	for _, t := range transactions {
		if !t.IsRecurring && t.Amount != expectedAmount {
			return fmt.Errorf("expected amount %d, got %d for transaction %s", expectedAmount, t.Amount, t.ID)
		}
	}

	return nil
}
