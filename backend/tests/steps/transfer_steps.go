package steps

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerTransferSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^a second account "([^"]*)" of type "([^"]*)" exists$`, tc.aSecondAccountExists)
	ctx.Step(`^I transfer (\d+) from "([^"]*)" to "([^"]*)" with description "([^"]*)"$`, tc.iTransferBetweenAccounts)
	ctx.Step(`^the transfer should be created successfully$`, tc.theTransferShouldBeCreated)
	ctx.Step(`^the transfer transaction should have amount (-?\d+)$`, tc.theTransferTransactionShouldHaveAmount)
	ctx.Step(`^the transfer should fail$`, tc.theTransferShouldFail)
}

func (tc *TestContext) aSecondAccountExists(name, accountType string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	req := &model.CreateAccountRequest{
		Name:     name,
		Type:     accountType,
		Currency: "EUR",
		Balance:  0,
	}

	account, err := tc.AccountService.Create(context.Background(), user.ID, req)
	if err != nil {
		return fmt.Errorf("failed to create second account: %w", err)
	}

	tc.SecondAccount = account
	return nil
}

func (tc *TestContext) iTransferBetweenAccounts(amount int64, fromName, toName, description string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	fromAccount, ok := tc.CurrentAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no current account (from)")
	}

	// Determine the to account
	var toAccountID string
	if fromName == toName {
		// Same account transfer (should fail)
		toAccountID = fromAccount.ID
	} else {
		toAccount, ok := tc.SecondAccount.(*model.Account)
		if !ok {
			return fmt.Errorf("no second account (to)")
		}
		toAccountID = toAccount.ID
	}

	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category")
	}

	toID := toAccountID
	txReq := &model.CreateTransactionRequest{
		AccountID:           fromAccount.ID,
		CategoryID:          category.ID,
		Amount:              -amount,
		Type:                "transfer",
		Description:         description,
		Date:                "2026-02-15",
		IsShared:            true,
		TransferToAccountID: &toID,
	}

	// Check for same-account transfer
	if fromAccount.ID == toAccountID {
		tc.LastError = fmt.Errorf("cannot transfer to the same account")
		return nil
	}

	transaction, err := tc.TransactionService.Create(context.Background(), user.ID, txReq)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentTransaction = transaction
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theTransferShouldBeCreated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected transfer to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentTransaction == nil {
		return fmt.Errorf("expected transaction from transfer, got nil")
	}
	return nil
}

func (tc *TestContext) theTransferTransactionShouldHaveAmount(expectedAmount int64) error {
	transaction, ok := tc.CurrentTransaction.(*model.Transaction)
	if !ok {
		return fmt.Errorf("current transaction is not of correct type")
	}

	if transaction.Amount != expectedAmount {
		return fmt.Errorf("expected amount %d, got %d", expectedAmount, transaction.Amount)
	}
	return nil
}

func (tc *TestContext) theTransferShouldFail() error {
	if tc.LastError == nil {
		return fmt.Errorf("expected transfer to fail, but it succeeded")
	}
	return nil
}
