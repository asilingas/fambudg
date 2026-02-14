Feature: User Authentication
  As a user
  I want to register and login
  So that I can access my budget data

  Scenario: User registration with valid data
    When I register with email "test@example.com" password "password123" and name "Test User"
    Then the registration should succeed
    And I should receive a user object with email "test@example.com"

  Scenario: User login with valid credentials
    Given a user exists with email "john@example.com" and password "password123"
    When I login with email "john@example.com" and password "password123"
    Then the login should succeed
    And I should receive a JWT token
    And I should receive a user object

  Scenario: User login with invalid password
    Given a user exists with email "john@example.com" and password "password123"
    When I login with email "john@example.com" and password "wrongpassword"
    Then the login should fail with error "invalid email or password"

  Scenario: User login with non-existent email
    When I login with email "nonexistent@example.com" and password "password123"
    Then the login should fail with error "invalid email or password"
