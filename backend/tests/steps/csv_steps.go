package steps

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerCSVSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I export transactions as CSV$`, tc.iExportTransactionsAsCSV)
	ctx.Step(`^the CSV should contain (\d+) data rows$`, tc.theCSVShouldContainNDataRows)
	ctx.Step(`^the CSV header should include "([^"]*)"$`, tc.theCSVHeaderShouldInclude)
	ctx.Step(`^I import the following CSV:$`, tc.iImportTheFollowingCSV)
	ctx.Step(`^(\d+) transactions should be imported$`, tc.nTransactionsShouldBeImported)
}

func (tc *TestContext) iExportTransactionsAsCSV() error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	transactions, err := tc.TransactionService.GetByUserID(context.Background(), user.ID, &model.TransactionFilters{})
	if err != nil {
		tc.LastError = err
		return nil
	}

	// Build CSV in memory
	var lines []string
	lines = append(lines, "date,amount,type,description,category_id,account_id,is_shared")
	for _, t := range transactions {
		line := fmt.Sprintf("%s,%d,%s,%s,%s,%s,%t",
			t.Date.Format("2006-01-02"),
			t.Amount,
			t.Type,
			t.Description,
			t.CategoryID,
			t.AccountID,
			t.IsShared,
		)
		lines = append(lines, line)
	}

	tc.ExportedCSV = lines
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theCSVShouldContainNDataRows(expectedCount int) error {
	// Subtract header row
	dataRows := len(tc.ExportedCSV) - 1
	if dataRows != expectedCount {
		return fmt.Errorf("expected %d CSV data rows, got %d", expectedCount, dataRows)
	}
	return nil
}

func (tc *TestContext) theCSVHeaderShouldInclude(expectedFields string) error {
	if len(tc.ExportedCSV) == 0 {
		return fmt.Errorf("no CSV data available")
	}

	header := tc.ExportedCSV[0]
	for _, field := range strings.Split(expectedFields, ",") {
		if !strings.Contains(header, field) {
			return fmt.Errorf("CSV header missing field %q, got: %s", field, header)
		}
	}
	return nil
}

func (tc *TestContext) iImportTheFollowingCSV(table *godog.Table) error {
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

	var imported int
	for _, row := range table.Rows[1:] {
		amount, _ := strconv.ParseInt(row.Cells[1].Value, 10, 64)
		isShared, _ := strconv.ParseBool(row.Cells[4].Value)

		req := &model.CreateTransactionRequest{
			Date:        row.Cells[0].Value,
			Amount:      amount,
			Type:        row.Cells[2].Value,
			Description: row.Cells[3].Value,
			CategoryID:  category.ID,
			AccountID:   account.ID,
			IsShared:    isShared,
		}

		_, err := tc.TransactionService.Create(context.Background(), user.ID, req)
		if err != nil {
			tc.LastError = err
			return nil
		}
		imported++
	}

	tc.ImportedCount = imported
	tc.LastError = nil
	return nil
}

func (tc *TestContext) nTransactionsShouldBeImported(expected int) error {
	if tc.LastError != nil {
		return fmt.Errorf("expected success, got error: %v", tc.LastError)
	}
	if tc.ImportedCount != expected {
		return fmt.Errorf("expected %d imported, got %d", expected, tc.ImportedCount)
	}
	return nil
}
