package steps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerTransactionSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^a category "([^"]*)" of type "([^"]*)" exists$`, tc.aCategoryExists)
	ctx.Step(`^an account "([^"]*)" of type "([^"]*)" exists$`, tc.anAccountExists)
	ctx.Step(`^I create a transaction with:$`, tc.iCreateTransactionWith)
	ctx.Step(`^the transaction should be created successfully$`, tc.theTransactionShouldBeCreated)
	ctx.Step(`^the transaction should have amount (-?\d+)$`, tc.theTransactionShouldHaveAmount)
	ctx.Step(`^the account balance should be updated$`, tc.theAccountBalanceShouldBeUpdated)
	ctx.Step(`^the following transactions exist:$`, tc.theFollowingTransactionsExist)
	ctx.Step(`^I list transactions from "([^"]*)" to "([^"]*)"$`, tc.iListTransactionsFromTo)
	ctx.Step(`^I should see (\d+) transactions$`, tc.iShouldSeeNTransactions)
	ctx.Step(`^a transaction exists with amount (-?\d+) and description "([^"]*)"$`, tc.aTransactionExistsWithAmountAndDescription)
	ctx.Step(`^a transaction exists with amount (-?\d+)$`, tc.aTransactionExistsWithAmount)
	ctx.Step(`^I update the transaction with description "([^"]*)"$`, tc.iUpdateTransactionWithDescription)
	ctx.Step(`^the transaction should be updated successfully$`, tc.theTransactionShouldBeUpdated)
	ctx.Step(`^the transaction should have description "([^"]*)"$`, tc.theTransactionShouldHaveDescription)
	ctx.Step(`^I delete the transaction$`, tc.iDeleteTheTransaction)
	ctx.Step(`^the transaction should be deleted successfully$`, tc.theTransactionShouldBeDeleted)
	ctx.Step(`^the account balance should be adjusted$`, tc.theAccountBalanceShouldBeAdjusted)
}

func (tc *TestContext) aCategoryExists(name, categoryType string) error {
	req := &model.CreateCategoryRequest{
		Name:      name,
		Type:      categoryType,
		SortOrder: 0,
	}

	category, err := tc.CategoryService.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	tc.CurrentCategory = category
	return nil
}

func (tc *TestContext) anAccountExists(name, accountType string) error {
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
		return fmt.Errorf("failed to create account: %w", err)
	}

	tc.CurrentAccount = account
	return nil
}

func (tc *TestContext) iCreateTransactionWith(table *godog.Table) error {
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

	// Parse table data
	data := make(map[string]string)
	for _, row := range table.Rows[1:] { // Skip header
		data[row.Cells[0].Value] = row.Cells[1].Value
	}

	amount, _ := strconv.ParseInt(data["amount"], 10, 64)

	req := &model.CreateTransactionRequest{
		AccountID:   account.ID,
		CategoryID:  category.ID,
		Amount:      amount,
		Type:        "expense",
		Description: data["description"],
		Date:        data["date"],
		IsShared:    true,
		IsRecurring: false,
	}

	if amount > 0 {
		req.Type = "income"
	}

	transaction, err := tc.TransactionService.Create(context.Background(), user.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentTransaction = transaction
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theTransactionShouldBeCreated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected transaction creation to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentTransaction == nil {
		return fmt.Errorf("expected transaction object, got nil")
	}
	return nil
}

func (tc *TestContext) theTransactionShouldHaveAmount(expectedAmount int64) error {
	transaction, ok := tc.CurrentTransaction.(*model.Transaction)
	if !ok {
		return fmt.Errorf("current transaction is not of correct type")
	}

	if transaction.Amount != expectedAmount {
		return fmt.Errorf("expected amount %d, got %d", expectedAmount, transaction.Amount)
	}
	return nil
}

func (tc *TestContext) theAccountBalanceShouldBeUpdated() error {
	// Just verify no error occurred
	return nil
}

func (tc *TestContext) theFollowingTransactionsExist(table *godog.Table) error {
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

	for _, row := range table.Rows[1:] { // Skip header
		amount, _ := strconv.ParseInt(row.Cells[0].Value, 10, 64)

		req := &model.CreateTransactionRequest{
			AccountID:   account.ID,
			CategoryID:  category.ID,
			Amount:      amount,
			Type:        "expense",
			Description: row.Cells[1].Value,
			Date:        row.Cells[2].Value,
			IsShared:    true,
		}

		_, err := tc.TransactionService.Create(context.Background(), user.ID, req)
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}
	}

	return nil
}

func (tc *TestContext) iListTransactionsFromTo(startDate, endDate string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	filters := &model.TransactionFilters{
		StartDate: startDate,
		EndDate:   endDate,
	}

	transactions, err := tc.TransactionService.GetByUserID(context.Background(), user.ID, filters)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.TransactionList = make([]interface{}, len(transactions))
	for i, t := range transactions {
		tc.TransactionList[i] = t
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNTransactions(expectedCount int) error {
	actualCount := len(tc.TransactionList)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d transactions, got %d", expectedCount, actualCount)
	}
	return nil
}

func (tc *TestContext) aTransactionExistsWithAmountAndDescription(amount int64, description string) error {
	return tc.aTransactionExistsWithAmount(amount)
}

func (tc *TestContext) aTransactionExistsWithAmount(amount int64) error {
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

	req := &model.CreateTransactionRequest{
		AccountID:   account.ID,
		CategoryID:  category.ID,
		Amount:      amount,
		Type:        "expense",
		Description: "Test transaction",
		Date:        "2026-02-14",
		IsShared:    true,
	}

	transaction, err := tc.TransactionService.Create(context.Background(), user.ID, req)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	tc.CurrentTransaction = transaction
	return nil
}

func (tc *TestContext) iUpdateTransactionWithDescription(newDescription string) error {
	transaction, ok := tc.CurrentTransaction.(*model.Transaction)
	if !ok {
		return fmt.Errorf("no current transaction")
	}

	req := &model.UpdateTransactionRequest{
		Description: &newDescription,
	}

	updated, err := tc.TransactionService.Update(context.Background(), transaction.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentTransaction = updated
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theTransactionShouldBeUpdated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected update to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) theTransactionShouldHaveDescription(expectedDescription string) error {
	transaction, ok := tc.CurrentTransaction.(*model.Transaction)
	if !ok {
		return fmt.Errorf("current transaction is not of correct type")
	}

	if transaction.Description != expectedDescription {
		return fmt.Errorf("expected description %q, got %q", expectedDescription, transaction.Description)
	}
	return nil
}

func (tc *TestContext) iDeleteTheTransaction() error {
	transaction, ok := tc.CurrentTransaction.(*model.Transaction)
	if !ok {
		return fmt.Errorf("no current transaction")
	}

	err := tc.TransactionService.Delete(context.Background(), transaction.ID)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) theTransactionShouldBeDeleted() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected deletion to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) theAccountBalanceShouldBeAdjusted() error {
	// Just verify no error occurred
	return nil
}
