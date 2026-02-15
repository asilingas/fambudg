package steps

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerRoleSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I create a user with email "([^"]*)" password "([^"]*)" name "([^"]*)" and role "([^"]*)"$`, tc.iCreateUserWithRole)
	ctx.Step(`^the user creation should succeed$`, tc.theUserCreationShouldSucceed)
	ctx.Step(`^the created user should have role "([^"]*)"$`, tc.theCreatedUserShouldHaveRole)
	ctx.Step(`^a user exists with email "([^"]*)" password "([^"]*)" name "([^"]*)" and role "([^"]*)"$`, tc.aUserExistsWithRole)
	ctx.Step(`^I list all users$`, tc.iListAllUsers)
	ctx.Step(`^I should see (\d+) users$`, tc.iShouldSeeNUsers)
	ctx.Step(`^I update the user role to "([^"]*)"$`, tc.iUpdateTheUserRoleTo)
	ctx.Step(`^the user update should succeed$`, tc.theUserUpdateShouldSucceed)
	ctx.Step(`^the updated user should have role "([^"]*)"$`, tc.theUpdatedUserShouldHaveRole)
	ctx.Step(`^I delete the user$`, tc.iDeleteTheUser)
	ctx.Step(`^the user deletion should succeed$`, tc.theUserDeletionShouldSucceed)
	ctx.Step(`^a child user "([^"]*)" exists with an account and a transaction of (-?\d+)$`, tc.aChildUserExistsWithAccountAndTransaction)
	ctx.Step(`^I list transactions for the child user$`, tc.iListTransactionsForChildUser)
	ctx.Step(`^I create a shared transaction with amount (-?\d+)$`, tc.iCreateSharedTransaction)
	ctx.Step(`^I create a personal transaction with amount (-?\d+)$`, tc.iCreatePersonalTransaction)
	ctx.Step(`^I filter transactions by isShared "([^"]*)"$`, tc.iFilterTransactionsByIsShared)
	ctx.Step(`^a child user "([^"]*)" exists$`, tc.aChildUserExists)
}

func (tc *TestContext) iCreateUserWithRole(email, password, name, role string) error {
	req := &model.CreateUserRequest{
		Email:    email,
		Password: password,
		Name:     name,
		Role:     role,
	}

	user, err := tc.AuthService.CreateUser(context.Background(), req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CreatedUser = user
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theUserCreationShouldSucceed() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected user creation to succeed, got error: %v", tc.LastError)
	}
	if tc.CreatedUser == nil {
		return fmt.Errorf("expected created user, got nil")
	}
	return nil
}

func (tc *TestContext) theCreatedUserShouldHaveRole(expectedRole string) error {
	user, ok := tc.CreatedUser.(*model.User)
	if !ok {
		return fmt.Errorf("created user is not of type *model.User")
	}
	if user.Role != expectedRole {
		return fmt.Errorf("expected role %s, got %s", expectedRole, user.Role)
	}
	return nil
}

func (tc *TestContext) aUserExistsWithRole(email, password, name, role string) error {
	req := &model.CreateUserRequest{
		Email:    email,
		Password: password,
		Name:     name,
		Role:     role,
	}

	user, err := tc.AuthService.CreateUser(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	tc.CreatedUser = user
	return nil
}

func (tc *TestContext) iListAllUsers() error {
	users, err := tc.AuthService.ListUsers(context.Background())
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.UserList = make([]any, len(users))
	for i, u := range users {
		tc.UserList[i] = u
	}
	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNUsers(expected int) error {
	actual := len(tc.UserList)
	if actual != expected {
		return fmt.Errorf("expected %d users, got %d", expected, actual)
	}
	return nil
}

func (tc *TestContext) iUpdateTheUserRoleTo(newRole string) error {
	user, ok := tc.CreatedUser.(*model.User)
	if !ok {
		return fmt.Errorf("no created user to update")
	}

	req := &model.UpdateUserRequest{
		Role: &newRole,
	}

	updated, err := tc.AuthService.UpdateUser(context.Background(), user.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CreatedUser = updated
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theUserUpdateShouldSucceed() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected update to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) theUpdatedUserShouldHaveRole(expectedRole string) error {
	return tc.theCreatedUserShouldHaveRole(expectedRole)
}

func (tc *TestContext) iDeleteTheUser() error {
	user, ok := tc.CreatedUser.(*model.User)
	if !ok {
		return fmt.Errorf("no created user to delete")
	}

	err := tc.AuthService.DeleteUser(context.Background(), user.ID)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) theUserDeletionShouldSucceed() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected deletion to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) aChildUserExists(email string) error {
	req := &model.CreateUserRequest{
		Email:    email,
		Password: "password123",
		Name:     "Child",
		Role:     "child",
	}

	user, err := tc.AuthService.CreateUser(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create child user: %w", err)
	}

	tc.ChildUser = user
	return nil
}

func (tc *TestContext) aChildUserExistsWithAccountAndTransaction(email string, amount int64) error {
	// Create child user
	req := &model.CreateUserRequest{
		Email:    email,
		Password: "password123",
		Name:     "Child",
		Role:     "child",
	}

	childUser, err := tc.AuthService.CreateUser(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create child user: %w", err)
	}

	tc.ChildUser = childUser

	// Create account for child
	account, ok := tc.CurrentAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no current account set")
	}

	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category set")
	}

	// Create transaction for the child
	txReq := &model.CreateTransactionRequest{
		AccountID:  account.ID,
		CategoryID: category.ID,
		Amount:     amount,
		Type:       "expense",
		Date:       "2026-02-10",
	}

	_, err = tc.TransactionService.Create(context.Background(), childUser.ID, txReq)
	if err != nil {
		return fmt.Errorf("failed to create child transaction: %w", err)
	}

	return nil
}

func (tc *TestContext) iListTransactionsForChildUser() error {
	childUser, ok := tc.ChildUser.(*model.User)
	if !ok {
		return fmt.Errorf("no child user set")
	}

	transactions, err := tc.TransactionService.GetByUserID(context.Background(), childUser.ID, &model.TransactionFilters{})
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.TransactionList = make([]any, len(transactions))
	for i, t := range transactions {
		tc.TransactionList[i] = t
	}
	tc.LastError = nil
	return nil
}

func (tc *TestContext) iCreateSharedTransaction(amount int64) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user set")
	}

	account, ok := tc.CurrentAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no current account set")
	}

	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category set")
	}

	req := &model.CreateTransactionRequest{
		AccountID:  account.ID,
		CategoryID: category.ID,
		Amount:     amount,
		Type:       "expense",
		Date:       "2026-02-10",
		IsShared:   true,
	}

	_, err := tc.TransactionService.Create(context.Background(), user.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iCreatePersonalTransaction(amount int64) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user set")
	}

	account, ok := tc.CurrentAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no current account set")
	}

	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category set")
	}

	req := &model.CreateTransactionRequest{
		AccountID:  account.ID,
		CategoryID: category.ID,
		Amount:     amount,
		Type:       "expense",
		Date:       "2026-02-10",
		IsShared:   false,
	}

	_, err := tc.TransactionService.Create(context.Background(), user.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iFilterTransactionsByIsShared(isSharedStr string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user set")
	}

	isShared := isSharedStr == "true"
	filters := &model.TransactionFilters{
		IsShared: &isShared,
	}

	transactions, err := tc.TransactionService.GetByUserID(context.Background(), user.ID, filters)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.TransactionList = make([]any, len(transactions))
	for i, t := range transactions {
		tc.TransactionList[i] = t
	}
	tc.LastError = nil
	return nil
}
