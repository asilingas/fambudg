package steps

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerTrendSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I get trends for the last (\d+) months$`, tc.iGetTrendsForLastNMonths)
	ctx.Step(`^I should see trend data$`, tc.iShouldSeeTrendData)
	ctx.Step(`^the trend data should contain entries for the transaction months$`, tc.theTrendDataShouldContainEntries)
}

func (tc *TestContext) iGetTrendsForLastNMonths(months int) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	trends, err := tc.ReportService.GetTrends(context.Background(), user.ID, months)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.TrendResult = make([]any, len(trends))
	for i, t := range trends {
		tc.TrendResult[i] = t
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeTrendData() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected trend data, got error: %v", tc.LastError)
	}
	if len(tc.TrendResult) == 0 {
		return fmt.Errorf("expected trend data, got empty result")
	}
	return nil
}

func (tc *TestContext) theTrendDataShouldContainEntries() error {
	if len(tc.TrendResult) < 2 {
		return fmt.Errorf("expected at least 2 trend entries, got %d", len(tc.TrendResult))
	}

	// Verify entries have expense data
	for _, item := range tc.TrendResult {
		tp, ok := item.(*model.TrendPoint)
		if !ok {
			return fmt.Errorf("trend entry is not of correct type")
		}
		if tp.Month == 0 || tp.Year == 0 {
			return fmt.Errorf("trend entry has invalid month/year: %d/%d", tp.Month, tp.Year)
		}
	}

	return nil
}
