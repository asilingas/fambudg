Feature: Transaction Management
  As a user
  I want to manage my transactions
  So that I can track my expenses and income

  Background:
    Given I am logged in as "user@example.com"
    And a category "Groceries" of type "expense" exists
    And an account "Chase Checking" of type "checking" exists

  Scenario: Create an expense transaction
    When I create a transaction with:
      | field       | value            |
      | amount      | -4599            |
      | description | Weekly groceries |
      | date        | 2026-02-14       |
    Then the transaction should be created successfully
    And the transaction should have amount -4599
    And the account balance should be updated

  Scenario: Create an income transaction
    Given a category "Salary" of type "income" exists
    When I create a transaction with:
      | field       | value       |
      | amount      | 500000      |
      | description | Monthly pay |
      | date        | 2026-02-01  |
    Then the transaction should be created successfully
    And the transaction should have amount 500000

  Scenario: List transactions with date filter
    Given the following transactions exist:
      | amount | description | date       |
      | -2000  | Coffee      | 2026-01-15 |
      | -5000  | Dinner      | 2026-02-10 |
      | -3000  | Lunch       | 2026-02-12 |
    When I list transactions from "2026-02-01" to "2026-02-28"
    Then I should see 2 transactions

  Scenario: Update a transaction
    Given a transaction exists with amount -1000 and description "Old description"
    When I update the transaction with description "New description"
    Then the transaction should be updated successfully
    And the transaction should have description "New description"

  Scenario: Delete a transaction
    Given a transaction exists with amount -1000
    When I delete the transaction
    Then the transaction should be deleted successfully
    And the account balance should be adjusted
