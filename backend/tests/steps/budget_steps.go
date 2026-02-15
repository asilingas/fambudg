package steps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerBudgetSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I create a budget with:$`, tc.iCreateBudgetWith)
	ctx.Step(`^the budget should be created successfully$`, tc.theBudgetShouldBeCreated)
	ctx.Step(`^the budget should have amount (\d+)$`, tc.theBudgetShouldHaveAmount)
	ctx.Step(`^the following budgets exist:$`, tc.theFollowingBudgetsExist)
	ctx.Step(`^I list budgets for month (\d+) and year (\d+)$`, tc.iListBudgetsForMonthAndYear)
	ctx.Step(`^I should see (\d+) budgets$`, tc.iShouldSeeNBudgets)
	ctx.Step(`^a budget exists with amount (\d+) for month (\d+) and year (\d+)$`, tc.aBudgetExistsWithAmount)
	ctx.Step(`^I update the budget with amount (\d+)$`, tc.iUpdateBudgetWithAmount)
	ctx.Step(`^the budget should be updated successfully$`, tc.theBudgetShouldBeUpdated)
	ctx.Step(`^I delete the budget$`, tc.iDeleteTheBudget)
	ctx.Step(`^the budget should be deleted successfully$`, tc.theBudgetShouldBeDeleted)
	ctx.Step(`^I get the budget summary for month (\d+) and year (\d+)$`, tc.iGetBudgetSummary)
	ctx.Step(`^I should see (\d+) budget summaries$`, tc.iShouldSeeNBudgetSummaries)
	ctx.Step(`^the budget summary for "([^"]*)" should have budget (\d+) and actual (\d+)$`, tc.theBudgetSummaryShouldHave)
}

func (tc *TestContext) iCreateBudgetWith(table *godog.Table) error {
	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category")
	}

	data := make(map[string]string)
	for _, row := range table.Rows[1:] {
		data[row.Cells[0].Value] = row.Cells[1].Value
	}

	amount, _ := strconv.ParseInt(data["amount"], 10, 64)
	month, _ := strconv.Atoi(data["month"])
	year, _ := strconv.Atoi(data["year"])

	req := &model.CreateBudgetRequest{
		CategoryID: category.ID,
		Amount:     amount,
		Month:      month,
		Year:       year,
	}

	budget, err := tc.BudgetService.Create(context.Background(), req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentBudget = budget
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theBudgetShouldBeCreated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected budget creation to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentBudget == nil {
		return fmt.Errorf("expected budget object, got nil")
	}
	return nil
}

func (tc *TestContext) theBudgetShouldHaveAmount(expectedAmount int64) error {
	budget, ok := tc.CurrentBudget.(*model.Budget)
	if !ok {
		return fmt.Errorf("current budget is not of correct type")
	}

	if budget.Amount != expectedAmount {
		return fmt.Errorf("expected amount %d, got %d", expectedAmount, budget.Amount)
	}
	return nil
}

func (tc *TestContext) theFollowingBudgetsExist(table *godog.Table) error {
	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category")
	}

	for _, row := range table.Rows[1:] {
		month, _ := strconv.Atoi(row.Cells[0].Value)
		year, _ := strconv.Atoi(row.Cells[1].Value)
		amount, _ := strconv.ParseInt(row.Cells[2].Value, 10, 64)

		req := &model.CreateBudgetRequest{
			CategoryID: category.ID,
			Amount:     amount,
			Month:      month,
			Year:       year,
		}

		_, err := tc.BudgetService.Create(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create budget: %w", err)
		}
	}

	return nil
}

func (tc *TestContext) iListBudgetsForMonthAndYear(month, year int) error {
	filters := &model.BudgetFilters{
		Month: month,
		Year:  year,
	}

	budgets, err := tc.BudgetService.GetAll(context.Background(), filters)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.BudgetList = make([]any, len(budgets))
	for i, b := range budgets {
		tc.BudgetList[i] = b
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNBudgets(expectedCount int) error {
	actualCount := len(tc.BudgetList)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d budgets, got %d", expectedCount, actualCount)
	}
	return nil
}

func (tc *TestContext) aBudgetExistsWithAmount(amount int64, month, year int) error {
	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category")
	}

	req := &model.CreateBudgetRequest{
		CategoryID: category.ID,
		Amount:     amount,
		Month:      month,
		Year:       year,
	}

	budget, err := tc.BudgetService.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create budget: %w", err)
	}

	tc.CurrentBudget = budget
	return nil
}

func (tc *TestContext) iUpdateBudgetWithAmount(newAmount int64) error {
	budget, ok := tc.CurrentBudget.(*model.Budget)
	if !ok {
		return fmt.Errorf("no current budget")
	}

	req := &model.UpdateBudgetRequest{
		Amount: &newAmount,
	}

	updated, err := tc.BudgetService.Update(context.Background(), budget.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentBudget = updated
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theBudgetShouldBeUpdated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected update to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) iDeleteTheBudget() error {
	budget, ok := tc.CurrentBudget.(*model.Budget)
	if !ok {
		return fmt.Errorf("no current budget")
	}

	err := tc.BudgetService.Delete(context.Background(), budget.ID)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) theBudgetShouldBeDeleted() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected deletion to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) iGetBudgetSummary(month, year int) error {
	summaries, err := tc.BudgetService.GetSummary(context.Background(), month, year)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.BudgetSummaryList = make([]any, len(summaries))
	for i, s := range summaries {
		tc.BudgetSummaryList[i] = s
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNBudgetSummaries(expectedCount int) error {
	actualCount := len(tc.BudgetSummaryList)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d budget summaries, got %d", expectedCount, actualCount)
	}
	return nil
}

func (tc *TestContext) theBudgetSummaryShouldHave(categoryName string, expectedBudget, expectedActual int64) error {
	for _, item := range tc.BudgetSummaryList {
		summary, ok := item.(*model.BudgetSummary)
		if !ok {
			continue
		}
		if summary.CategoryName == categoryName {
			if summary.BudgetAmount != expectedBudget {
				return fmt.Errorf("expected budget %d for %q, got %d", expectedBudget, categoryName, summary.BudgetAmount)
			}
			if summary.ActualAmount != expectedActual {
				return fmt.Errorf("expected actual %d for %q, got %d", expectedActual, categoryName, summary.ActualAmount)
			}
			return nil
		}
	}
	return fmt.Errorf("budget summary for %q not found", categoryName)
}
