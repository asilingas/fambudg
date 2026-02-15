package steps

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerAllowanceSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I create an allowance of (\d+) for the child starting "([^"]*)"$`, tc.iCreateAllowanceForChild)
	ctx.Step(`^the allowance should be created successfully$`, tc.theAllowanceShouldBeCreated)
	ctx.Step(`^the allowance amount should be (\d+)$`, tc.theAllowanceAmountShouldBe)
	ctx.Step(`^a child user "([^"]*)" exists with an account and a category$`, tc.aChildUserExistsWithAccountAndCategory)
	ctx.Step(`^an allowance of (\d+) for the child starting "([^"]*)"$`, tc.anAllowanceExistsForChild)
	ctx.Step(`^the child has a transaction of (-?\d+) on "([^"]*)"$`, tc.theChildHasTransactionOn)
	ctx.Step(`^I get allowances$`, tc.iGetAllowances)
	ctx.Step(`^the allowance spent should be (\d+)$`, tc.theAllowanceSpentShouldBe)
	ctx.Step(`^the allowance remaining should be (\d+)$`, tc.theAllowanceRemainingShouldBe)
	ctx.Step(`^I update the allowance amount to (\d+)$`, tc.iUpdateAllowanceAmountTo)
	ctx.Step(`^the allowance should be updated successfully$`, tc.theAllowanceShouldBeUpdated)
	ctx.Step(`^I should see (\d+) allowances$`, tc.iShouldSeeNAllowances)
}

func (tc *TestContext) iCreateAllowanceForChild(amount int64, periodStart string) error {
	childUser, ok := tc.ChildUser.(*model.User)
	if !ok {
		return fmt.Errorf("no child user set")
	}

	req := &model.CreateAllowanceRequest{
		UserID:      childUser.ID,
		Amount:      amount,
		PeriodStart: periodStart,
	}

	allowance, err := tc.AllowanceService.Create(context.Background(), req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentAllowance = allowance
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theAllowanceShouldBeCreated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected allowance creation to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentAllowance == nil {
		return fmt.Errorf("expected allowance, got nil")
	}
	return nil
}

func (tc *TestContext) theAllowanceAmountShouldBe(expected int64) error {
	allowance, ok := tc.CurrentAllowance.(*model.Allowance)
	if !ok {
		return fmt.Errorf("current allowance is not of type *model.Allowance")
	}
	if allowance.Amount != expected {
		return fmt.Errorf("expected amount %d, got %d", expected, allowance.Amount)
	}
	return nil
}

func (tc *TestContext) aChildUserExistsWithAccountAndCategory(email string) error {
	// Create child user
	if err := tc.aChildUserExists(email); err != nil {
		return err
	}

	childUser, ok := tc.ChildUser.(*model.User)
	if !ok {
		return fmt.Errorf("no child user set")
	}

	// Create an account for the child
	acctReq := &model.CreateAccountRequest{
		Name:     "Kid Account",
		Type:     "checking",
		Currency: "EUR",
	}
	account, err := tc.AccountService.Create(context.Background(), childUser.ID, acctReq)
	if err != nil {
		return fmt.Errorf("failed to create child account: %w", err)
	}
	tc.ChildAccount = account

	// Create a category
	catReq := &model.CreateCategoryRequest{
		Name: "Kid Spending",
		Type: "expense",
	}
	category, err := tc.CategoryService.Create(context.Background(), catReq)
	if err != nil {
		return fmt.Errorf("failed to create child category: %w", err)
	}
	tc.ChildCategory = category

	return nil
}

func (tc *TestContext) anAllowanceExistsForChild(amount int64, periodStart string) error {
	childUser, ok := tc.ChildUser.(*model.User)
	if !ok {
		return fmt.Errorf("no child user set")
	}

	req := &model.CreateAllowanceRequest{
		UserID:      childUser.ID,
		Amount:      amount,
		PeriodStart: periodStart,
	}

	allowance, err := tc.AllowanceService.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create allowance: %w", err)
	}

	tc.CurrentAllowance = allowance
	return nil
}

func (tc *TestContext) theChildHasTransactionOn(amount int64, date string) error {
	childUser, ok := tc.ChildUser.(*model.User)
	if !ok {
		return fmt.Errorf("no child user set")
	}

	account, ok := tc.ChildAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no child account set")
	}

	category, ok := tc.ChildCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no child category set")
	}

	req := &model.CreateTransactionRequest{
		AccountID:  account.ID,
		CategoryID: category.ID,
		Amount:     amount,
		Type:       "expense",
		Date:       date,
	}

	_, err := tc.TransactionService.Create(context.Background(), childUser.ID, req)
	if err != nil {
		return fmt.Errorf("failed to create child transaction: %w", err)
	}

	return nil
}

func (tc *TestContext) iGetAllowances() error {
	allowances, err := tc.AllowanceService.GetAll(context.Background())
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.AllowanceList = make([]any, len(allowances))
	for i, a := range allowances {
		tc.AllowanceList[i] = a
	}

	// If there's exactly one, set it as current for field assertions
	if len(allowances) == 1 {
		tc.CurrentAllowance = allowances[0]
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) theAllowanceSpentShouldBe(expected int64) error {
	allowance, ok := tc.CurrentAllowance.(*model.Allowance)
	if !ok {
		return fmt.Errorf("current allowance is not of type *model.Allowance")
	}
	if allowance.Spent != expected {
		return fmt.Errorf("expected spent %d, got %d", expected, allowance.Spent)
	}
	return nil
}

func (tc *TestContext) theAllowanceRemainingShouldBe(expected int64) error {
	allowance, ok := tc.CurrentAllowance.(*model.Allowance)
	if !ok {
		return fmt.Errorf("current allowance is not of type *model.Allowance")
	}
	if allowance.Remaining != expected {
		return fmt.Errorf("expected remaining %d, got %d", expected, allowance.Remaining)
	}
	return nil
}

func (tc *TestContext) iUpdateAllowanceAmountTo(newAmount int64) error {
	allowance, ok := tc.CurrentAllowance.(*model.Allowance)
	if !ok {
		return fmt.Errorf("no current allowance to update")
	}

	req := &model.UpdateAllowanceRequest{
		Amount: &newAmount,
	}

	updated, err := tc.AllowanceService.Update(context.Background(), allowance.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentAllowance = updated
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theAllowanceShouldBeUpdated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected allowance update to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) iShouldSeeNAllowances(expected int) error {
	actual := len(tc.AllowanceList)
	if actual != expected {
		return fmt.Errorf("expected %d allowances, got %d", expected, actual)
	}
	return nil
}
