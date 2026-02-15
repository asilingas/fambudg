package steps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerBillReminderSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I create a bill reminder with:$`, tc.iCreateBillReminderWith)
	ctx.Step(`^the bill reminder should be created successfully$`, tc.theBillReminderShouldBeCreated)
	ctx.Step(`^the bill reminder should have amount (\d+)$`, tc.theBillReminderShouldHaveAmount)
	ctx.Step(`^the following bill reminders exist:$`, tc.theFollowingBillRemindersExist)
	ctx.Step(`^I list bill reminders$`, tc.iListBillReminders)
	ctx.Step(`^I should see (\d+) bill reminders$`, tc.iShouldSeeNBillReminders)
	ctx.Step(`^a bill reminder "([^"]*)" exists with amount (\d+)$`, tc.aBillReminderExistsWithAmount)
	ctx.Step(`^I update the bill reminder with amount (\d+)$`, tc.iUpdateBillReminderWithAmount)
	ctx.Step(`^the bill reminder should be updated successfully$`, tc.theBillReminderShouldBeUpdated)
	ctx.Step(`^I delete the bill reminder$`, tc.iDeleteTheBillReminder)
	ctx.Step(`^the bill reminder should be deleted successfully$`, tc.theBillReminderShouldBeDeleted)
	ctx.Step(`^a bill reminder "([^"]*)" exists with amount (\d+) and account$`, tc.aBillReminderExistsWithAmountAndAccount)
	ctx.Step(`^I pay the bill reminder on "([^"]*)"$`, tc.iPayTheBillReminder)
	ctx.Step(`^the bill payment should create a transaction$`, tc.theBillPaymentShouldCreateTransaction)
	ctx.Step(`^the bill next due date should be advanced$`, tc.theBillNextDueDateShouldBeAdvanced)
}

func (tc *TestContext) iCreateBillReminderWith(table *godog.Table) error {
	data := make(map[string]string)
	for _, row := range table.Rows[1:] {
		data[row.Cells[0].Value] = row.Cells[1].Value
	}

	amount, _ := strconv.ParseInt(data["amount"], 10, 64)
	dueDay, _ := strconv.Atoi(data["dueDay"])

	req := &model.CreateBillReminderRequest{
		Name:        data["name"],
		Amount:      amount,
		DueDay:      dueDay,
		Frequency:   data["frequency"],
		NextDueDate: data["nextDueDate"],
	}

	bill, err := tc.BillReminderService.Create(context.Background(), req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentBillReminder = bill
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theBillReminderShouldBeCreated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected bill reminder creation to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentBillReminder == nil {
		return fmt.Errorf("expected bill reminder object, got nil")
	}
	return nil
}

func (tc *TestContext) theBillReminderShouldHaveAmount(expectedAmount int64) error {
	bill, ok := tc.CurrentBillReminder.(*model.BillReminder)
	if !ok {
		return fmt.Errorf("current bill reminder is not of correct type")
	}

	if bill.Amount != expectedAmount {
		return fmt.Errorf("expected amount %d, got %d", expectedAmount, bill.Amount)
	}
	return nil
}

func (tc *TestContext) theFollowingBillRemindersExist(table *godog.Table) error {
	for _, row := range table.Rows[1:] {
		amount, _ := strconv.ParseInt(row.Cells[1].Value, 10, 64)
		dueDay, _ := strconv.Atoi(row.Cells[2].Value)

		req := &model.CreateBillReminderRequest{
			Name:        row.Cells[0].Value,
			Amount:      amount,
			DueDay:      dueDay,
			Frequency:   row.Cells[3].Value,
			NextDueDate: row.Cells[4].Value,
		}

		_, err := tc.BillReminderService.Create(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create bill reminder: %w", err)
		}
	}

	return nil
}

func (tc *TestContext) iListBillReminders() error {
	bills, err := tc.BillReminderService.GetAll(context.Background())
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.BillReminderList = make([]any, len(bills))
	for i, b := range bills {
		tc.BillReminderList[i] = b
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNBillReminders(expectedCount int) error {
	actualCount := len(tc.BillReminderList)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d bill reminders, got %d", expectedCount, actualCount)
	}
	return nil
}

func (tc *TestContext) aBillReminderExistsWithAmount(name string, amount int64) error {
	req := &model.CreateBillReminderRequest{
		Name:        name,
		Amount:      amount,
		DueDay:      15,
		Frequency:   "monthly",
		NextDueDate: "2026-03-15",
	}

	bill, err := tc.BillReminderService.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create bill reminder: %w", err)
	}

	tc.CurrentBillReminder = bill
	return nil
}

func (tc *TestContext) iUpdateBillReminderWithAmount(newAmount int64) error {
	bill, ok := tc.CurrentBillReminder.(*model.BillReminder)
	if !ok {
		return fmt.Errorf("no current bill reminder")
	}

	req := &model.UpdateBillReminderRequest{
		Amount: &newAmount,
	}

	updated, err := tc.BillReminderService.Update(context.Background(), bill.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentBillReminder = updated
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theBillReminderShouldBeUpdated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected update to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) iDeleteTheBillReminder() error {
	bill, ok := tc.CurrentBillReminder.(*model.BillReminder)
	if !ok {
		return fmt.Errorf("no current bill reminder")
	}

	err := tc.BillReminderService.Delete(context.Background(), bill.ID)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) theBillReminderShouldBeDeleted() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected deletion to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) aBillReminderExistsWithAmountAndAccount(name string, amount int64) error {
	account, ok := tc.CurrentAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no current account")
	}

	category, ok := tc.CurrentCategory.(*model.Category)
	if !ok {
		return fmt.Errorf("no current category")
	}

	req := &model.CreateBillReminderRequest{
		Name:        name,
		Amount:      amount,
		DueDay:      15,
		Frequency:   "monthly",
		CategoryID:  &category.ID,
		AccountID:   &account.ID,
		NextDueDate: "2026-03-15",
	}

	bill, err := tc.BillReminderService.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create bill reminder: %w", err)
	}

	tc.CurrentBillReminder = bill
	return nil
}

func (tc *TestContext) iPayTheBillReminder(date string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("no current user")
	}

	bill, ok := tc.CurrentBillReminder.(*model.BillReminder)
	if !ok {
		return fmt.Errorf("no current bill reminder")
	}

	account, ok := tc.CurrentAccount.(*model.Account)
	if !ok {
		return fmt.Errorf("no current account")
	}

	req := &model.PayBillRequest{
		AccountID: account.ID,
		Date:      date,
	}

	transaction, err := tc.BillReminderService.Pay(context.Background(), user.ID, bill.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentTransaction = transaction
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theBillPaymentShouldCreateTransaction() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected bill payment to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentTransaction == nil {
		return fmt.Errorf("expected transaction from bill payment, got nil")
	}
	return nil
}

func (tc *TestContext) theBillNextDueDateShouldBeAdvanced() error {
	bill, ok := tc.CurrentBillReminder.(*model.BillReminder)
	if !ok {
		return fmt.Errorf("no current bill reminder")
	}

	// Re-fetch the bill to get updated next_due_date
	updated, err := tc.BillReminderService.GetByID(context.Background(), bill.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch updated bill: %w", err)
	}

	if !updated.NextDueDate.After(bill.NextDueDate) {
		return fmt.Errorf("expected next due date to be advanced, got %v (original: %v)", updated.NextDueDate, bill.NextDueDate)
	}

	return nil
}
