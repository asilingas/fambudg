package steps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/asilingas/fambudg/backend/internal/model"
	"github.com/cucumber/godog"
)

func registerSavingGoalSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I create a saving goal with:$`, tc.iCreateSavingGoalWith)
	ctx.Step(`^the saving goal should be created successfully$`, tc.theSavingGoalShouldBeCreated)
	ctx.Step(`^the saving goal should have target amount (\d+)$`, tc.theSavingGoalShouldHaveTargetAmount)
	ctx.Step(`^the following saving goals exist:$`, tc.theFollowingSavingGoalsExist)
	ctx.Step(`^I list saving goals$`, tc.iListSavingGoals)
	ctx.Step(`^I should see (\d+) saving goals$`, tc.iShouldSeeNSavingGoals)
	ctx.Step(`^a saving goal "([^"]*)" exists with target (\d+)$`, tc.aSavingGoalExistsWithTarget)
	ctx.Step(`^I update the saving goal with name "([^"]*)"$`, tc.iUpdateSavingGoalWithName)
	ctx.Step(`^the saving goal should be updated successfully$`, tc.theSavingGoalShouldBeUpdated)
	ctx.Step(`^the saving goal name should be "([^"]*)"$`, tc.theSavingGoalNameShouldBe)
	ctx.Step(`^I contribute (\d+) to the saving goal$`, tc.iContributeToSavingGoal)
	ctx.Step(`^the saving goal current amount should be (\d+)$`, tc.theSavingGoalCurrentAmountShouldBe)
	ctx.Step(`^the saving goal status should be "([^"]*)"$`, tc.theSavingGoalStatusShouldBe)
}

func (tc *TestContext) iCreateSavingGoalWith(table *godog.Table) error {
	data := make(map[string]string)
	for _, row := range table.Rows[1:] {
		data[row.Cells[0].Value] = row.Cells[1].Value
	}

	targetAmount, _ := strconv.ParseInt(data["targetAmount"], 10, 64)
	priority, _ := strconv.Atoi(data["priority"])
	if priority == 0 {
		priority = 1
	}

	req := &model.CreateSavingGoalRequest{
		Name:         data["name"],
		TargetAmount: targetAmount,
		Priority:     priority,
	}

	goal, err := tc.SavingGoalService.Create(context.Background(), req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentSavingGoal = goal
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theSavingGoalShouldBeCreated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected saving goal creation to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentSavingGoal == nil {
		return fmt.Errorf("expected saving goal object, got nil")
	}
	return nil
}

func (tc *TestContext) theSavingGoalShouldHaveTargetAmount(expectedAmount int64) error {
	goal, ok := tc.CurrentSavingGoal.(*model.SavingGoal)
	if !ok {
		return fmt.Errorf("current saving goal is not of correct type")
	}

	if goal.TargetAmount != expectedAmount {
		return fmt.Errorf("expected target amount %d, got %d", expectedAmount, goal.TargetAmount)
	}
	return nil
}

func (tc *TestContext) theFollowingSavingGoalsExist(table *godog.Table) error {
	for _, row := range table.Rows[1:] {
		targetAmount, _ := strconv.ParseInt(row.Cells[1].Value, 10, 64)
		priority, _ := strconv.Atoi(row.Cells[2].Value)

		req := &model.CreateSavingGoalRequest{
			Name:         row.Cells[0].Value,
			TargetAmount: targetAmount,
			Priority:     priority,
		}

		_, err := tc.SavingGoalService.Create(context.Background(), req)
		if err != nil {
			return fmt.Errorf("failed to create saving goal: %w", err)
		}
	}

	return nil
}

func (tc *TestContext) iListSavingGoals() error {
	goals, err := tc.SavingGoalService.GetAll(context.Background())
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.SavingGoalList = make([]any, len(goals))
	for i, g := range goals {
		tc.SavingGoalList[i] = g
	}

	tc.LastError = nil
	return nil
}

func (tc *TestContext) iShouldSeeNSavingGoals(expectedCount int) error {
	actualCount := len(tc.SavingGoalList)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d saving goals, got %d", expectedCount, actualCount)
	}
	return nil
}

func (tc *TestContext) aSavingGoalExistsWithTarget(name string, targetAmount int64) error {
	req := &model.CreateSavingGoalRequest{
		Name:         name,
		TargetAmount: targetAmount,
		Priority:     1,
	}

	goal, err := tc.SavingGoalService.Create(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create saving goal: %w", err)
	}

	tc.CurrentSavingGoal = goal
	return nil
}

func (tc *TestContext) iUpdateSavingGoalWithName(newName string) error {
	goal, ok := tc.CurrentSavingGoal.(*model.SavingGoal)
	if !ok {
		return fmt.Errorf("no current saving goal")
	}

	req := &model.UpdateSavingGoalRequest{
		Name: &newName,
	}

	updated, err := tc.SavingGoalService.Update(context.Background(), goal.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentSavingGoal = updated
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theSavingGoalShouldBeUpdated() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected update to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) theSavingGoalNameShouldBe(expectedName string) error {
	goal, ok := tc.CurrentSavingGoal.(*model.SavingGoal)
	if !ok {
		return fmt.Errorf("current saving goal is not of correct type")
	}

	if goal.Name != expectedName {
		return fmt.Errorf("expected name %q, got %q", expectedName, goal.Name)
	}
	return nil
}

func (tc *TestContext) iContributeToSavingGoal(amount int64) error {
	goal, ok := tc.CurrentSavingGoal.(*model.SavingGoal)
	if !ok {
		return fmt.Errorf("no current saving goal")
	}

	req := &model.ContributeRequest{
		Amount: amount,
	}

	updated, err := tc.SavingGoalService.Contribute(context.Background(), goal.ID, req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentSavingGoal = updated
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theSavingGoalCurrentAmountShouldBe(expectedAmount int64) error {
	goal, ok := tc.CurrentSavingGoal.(*model.SavingGoal)
	if !ok {
		return fmt.Errorf("current saving goal is not of correct type")
	}

	if goal.CurrentAmount != expectedAmount {
		return fmt.Errorf("expected current amount %d, got %d", expectedAmount, goal.CurrentAmount)
	}
	return nil
}

func (tc *TestContext) theSavingGoalStatusShouldBe(expectedStatus string) error {
	goal, ok := tc.CurrentSavingGoal.(*model.SavingGoal)
	if !ok {
		return fmt.Errorf("current saving goal is not of correct type")
	}

	if goal.Status != expectedStatus {
		return fmt.Errorf("expected status %q, got %q", expectedStatus, goal.Status)
	}
	return nil
}
