Feature: Recurring Transaction Generation
  As a user with recurring transactions
  I want to auto-generate transaction instances
  So that I don't have to manually create them each period

  Background:
    Given I am logged in as "recuruser@example.com"
    And an account "Main Account" of type "checking" exists
    And a category "Rent" of type "expense" exists

  Scenario: Generate monthly recurring transactions
    Given I have a recurring transaction of -150000 on "2025-10-01" with frequency "monthly"
    When I generate recurring transactions up to "2026-01-15"
    Then 3 recurring transactions should have been generated
    And the generated transactions should have amounts of -150000

  Scenario: No duplicates when generating again
    Given I have a recurring transaction of -5000 on "2025-12-01" with frequency "monthly"
    When I generate recurring transactions up to "2026-02-15"
    Then 2 recurring transactions should have been generated
    When I generate recurring transactions up to "2026-02-15"
    Then 0 recurring transactions should have been generated

  Scenario: Generate weekly recurring transactions
    Given I have a recurring transaction of -2500 on "2026-01-01" with frequency "weekly"
    When I generate recurring transactions up to "2026-01-22"
    Then 3 recurring transactions should have been generated
