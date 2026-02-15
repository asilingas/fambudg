package steps

import (
	"context"
	"fmt"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerReportSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I get the dashboard for month (\d+) and year (\d+)$`, tc.iGetDashboard)
	ctx.Step(`^the month summary should have income (\d+) and expense (\d+)$`, tc.theMonthSummaryShouldHave)
	ctx.Step(`^I should see (\d+) recent transactions$`, tc.iShouldSeeNRecentTransactions)
	ctx.Step(`^I get the monthly report for month (\d+) and year (\d+)$`, tc.iGetMonthlyReport)
	ctx.Step(`^the monthly report should have income (\d+) and expense (\d+)$`, tc.theMonthlyReportShouldHave)
	ctx.Step(`^I get the category report for month (\d+) and year (\d+)$`, tc.iGetCategoryReport)
	ctx.Step(`^I should see (\d+) category entries$`, tc.iShouldSeeNCategoryEntries)
	ctx.Step(`^the category "([^"]*)" should have total (\d+)$`, tc.theCategoryShouldHaveTotal)
}

func (tc *TestContext) iGetDashboard(month, year int) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	dashboard, err := tc.ReportService.GetDashboard(context.Background(), user.ID, month, year)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.DashboardResult = dashboard
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theMonthSummaryShouldHave(expectedIncome, expectedExpense int64) error {
	dashboard, ok := tc.DashboardResult.(*model.DashboardResponse)
	if !ok {
		return fmt.Errorf("no dashboard result")
	}

	if dashboard.MonthSummary.TotalIncome != expectedIncome {
		return fmt.Errorf("expected income %d, got %d", expectedIncome, dashboard.MonthSummary.TotalIncome)
	}
	if dashboard.MonthSummary.TotalExpense != expectedExpense {
		return fmt.Errorf("expected expense %d, got %d", expectedExpense, dashboard.MonthSummary.TotalExpense)
	}
	return nil
}

func (tc *TestContext) iShouldSeeNRecentTransactions(expectedCount int) error {
	dashboard, ok := tc.DashboardResult.(*model.DashboardResponse)
	if !ok {
		return fmt.Errorf("no dashboard result")
	}

	actualCount := len(dashboard.RecentTransactions)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d recent transactions, got %d", expectedCount, actualCount)
	}
	return nil
}

func (tc *TestContext) iGetMonthlyReport(month, year int) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	summary, err := tc.ReportService.GetMonthlySummary(context.Background(), user.ID, month, year)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.MonthlyReportResult = summary
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theMonthlyReportShouldHave(expectedIncome, expectedExpense int64) error {
	summary, ok := tc.MonthlyReportResult.(*model.MonthSummary)
	if !ok {
		return fmt.Errorf("no monthly report result")
	}

	if summary.TotalIncome != expectedIncome {
		return fmt.Errorf("expected income %d, got %d", expectedIncome, summary.TotalIncome)
	}
	if summary.TotalExpense != expectedExpense {
		return fmt.Errorf("expected expense %d, got %d", expectedExpense, summary.TotalExpense)
	}
	return nil
}

func (tc *TestContext) iGetCategoryReport(month, year int) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	spending, err := tc.ReportService.GetSpendingByCategory(context.Background(), user.ID, month, year)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CategoryReportResult = make([]any, len(spending))
	for i, s := range spending {
		tc.CategoryReportResult[i] = s
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNCategoryEntries(expectedCount int) error {
	actualCount := len(tc.CategoryReportResult)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d category entries, got %d", expectedCount, actualCount)
	}
	return nil
}

func (tc *TestContext) theCategoryShouldHaveTotal(categoryName string, expectedTotal int64) error {
	for _, item := range tc.CategoryReportResult {
		cs, ok := item.(*model.CategorySpending)
		if !ok {
			continue
		}
		if cs.CategoryName == categoryName {
			if cs.TotalAmount != expectedTotal {
				return fmt.Errorf("expected total %d for %q, got %d", expectedTotal, categoryName, cs.TotalAmount)
			}
			return nil
		}
	}
	return fmt.Errorf("category %q not found in report", categoryName)
}
