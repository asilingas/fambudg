package steps

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	"github.com/yourusername/fambudg/backend/internal/model"
)

func registerAuthSteps(ctx *godog.ScenarioContext, tc *TestContext) {
	ctx.Step(`^I register with email "([^"]*)" password "([^"]*)" and name "([^"]*)"$`, tc.iRegisterWithCredentials)
	ctx.Step(`^the registration should succeed$`, tc.theRegistrationShouldSucceed)
	ctx.Step(`^I should receive a user object with email "([^"]*)"$`, tc.iShouldReceiveUserWithEmail)
	ctx.Step(`^a user exists with email "([^"]*)" and password "([^"]*)"$`, tc.aUserExistsWithCredentials)
	ctx.Step(`^I login with email "([^"]*)" and password "([^"]*)"$`, tc.iLoginWithCredentials)
	ctx.Step(`^the login should succeed$`, tc.theLoginShouldSucceed)
	ctx.Step(`^I should receive a JWT token$`, tc.iShouldReceiveJWTToken)
	ctx.Step(`^I should receive a user object$`, tc.iShouldReceiveUserObject)
	ctx.Step(`^the login should fail with error "([^"]*)"$`, tc.theLoginShouldFailWithError)
	ctx.Step(`^I am logged in as "([^"]*)"$`, tc.iAmLoggedInAs)
}

func (tc *TestContext) iRegisterWithCredentials(email, password, name string) error {
	req := &model.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
	}

	user, err := tc.AuthService.Register(context.Background(), req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentUser = user
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theRegistrationShouldSucceed() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected registration to succeed, got error: %v", tc.LastError)
	}
	if tc.CurrentUser == nil {
		return fmt.Errorf("expected user object, got nil")
	}
	return nil
}

func (tc *TestContext) iShouldReceiveUserWithEmail(expectedEmail string) error {
	user, ok := tc.CurrentUser.(*model.User)
	if !ok {
		return fmt.Errorf("current user is not of type *model.User")
	}

	if user.Email != expectedEmail {
		return fmt.Errorf("expected email %s, got %s", expectedEmail, user.Email)
	}
	return nil
}

func (tc *TestContext) aUserExistsWithCredentials(email, password string) error {
	req := &model.RegisterRequest{
		Email:    email,
		Password: password,
		Name:     "Test User",
	}

	user, err := tc.AuthService.Register(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to create test user: %w", err)
	}

	tc.CurrentUser = user
	return nil
}

func (tc *TestContext) iLoginWithCredentials(email, password string) error {
	req := &model.LoginRequest{
		Email:    email,
		Password: password,
	}

	loginResp, err := tc.AuthService.Login(context.Background(), req)
	if err != nil {
		tc.LastError = err
		return nil
	}

	tc.CurrentToken = loginResp.Token
	tc.CurrentUser = &loginResp.User
	tc.LastError = nil
	return nil
}

func (tc *TestContext) theLoginShouldSucceed() error {
	if tc.LastError != nil {
		return fmt.Errorf("expected login to succeed, got error: %v", tc.LastError)
	}
	return nil
}

func (tc *TestContext) iShouldReceiveJWTToken() error {
	if tc.CurrentToken == "" {
		return fmt.Errorf("expected JWT token, got empty string")
	}
	return nil
}

func (tc *TestContext) iShouldReceiveUserObject() error {
	if tc.CurrentUser == nil {
		return fmt.Errorf("expected user object, got nil")
	}
	return nil
}

func (tc *TestContext) theLoginShouldFailWithError(expectedError string) error {
	if tc.LastError == nil {
		return fmt.Errorf("expected error %q, got nil", expectedError)
	}

	if tc.LastError.Error() != expectedError {
		return fmt.Errorf("expected error %q, got %q", expectedError, tc.LastError.Error())
	}
	return nil
}

func (tc *TestContext) iAmLoggedInAs(email string) error {
	// Create and login user
	req := &model.RegisterRequest{
		Email:    email,
		Password: "password123",
		Name:     "Test User",
	}

	user, err := tc.AuthService.Register(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	tc.CurrentUser = user
	return nil
}
