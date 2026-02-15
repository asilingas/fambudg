package steps

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerSearchSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I search transactions with description "([^"]*)"$`, tc.iSearchByDescription)
	ctx.Step(`^I search transactions from "([^"]*)" to "([^"]*)"$`, tc.iSearchByDateRange)
	ctx.Step(`^I should see (\d+) search results$`, tc.iShouldSeeNSearchResults)
}

func (tc *TestContext) iSearchByDescription(description string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	filters := &model.SearchFilters{
		UserID:      user.ID,
		Description: description,
	}

	result, err := tc.ReportService.SearchTransactions(context.Background(), filters)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.SearchResult = result
	tc.LastError = nil
	return nil
}

func (tc *TestContext) iSearchByDateRange(startDate, endDate string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	filters := &model.SearchFilters{
		UserID:    user.ID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := tc.ReportService.SearchTransactions(context.Background(), filters)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.SearchResult = result
	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNSearchResults(expectedCount int) error {
	result, ok := tc.SearchResult.(*model.SearchResult)
	if !ok {
		return fmt.Errorf("no search result")
	}

	if result.TotalCount != expectedCount {
		return fmt.Errorf("expected %d search results, got %d", expectedCount, result.TotalCount)
	}
	return nil
}
