Feature: Budget management

  Background:
    Given I am logged in as "admin@family.com"
    And a category "Groceries" of type "expense" exists

  Scenario: Create a monthly budget
    When I create a budget with:
      | field  | value |
      | month  | 2     |
      | year   | 2026  |
      | amount | 50000 |
    Then the budget should be created successfully
    And the budget should have amount 50000

  Scenario: List budgets filtered by month and year
    Given the following budgets exist:
      | month | year | amount |
      | 1     | 2026 | 40000  |
      | 2     | 2026 | 50000  |
      | 3     | 2026 | 45000  |
    When I list budgets for month 2 and year 2026
    Then I should see 1 budgets

  Scenario: Update a budget amount
    Given a budget exists with amount 50000 for month 2 and year 2026
    When I update the budget with amount 60000
    Then the budget should be updated successfully
    And the budget should have amount 60000

  Scenario: Delete a budget
    Given a budget exists with amount 50000 for month 2 and year 2026
    When I delete the budget
    Then the budget should be deleted successfully

  Scenario: Budget vs actual summary
    Given a budget exists with amount 50000 for month 2 and year 2026
    And an account "Chase Checking" of type "checking" exists
    And the following transactions exist:
      | amount | description      | date       |
      | -15000 | Weekly groceries | 2026-02-05 |
      | -10000 | More groceries   | 2026-02-12 |
    When I get the budget summary for month 2 and year 2026
    Then I should see 1 budget summaries
    And the budget summary for "Groceries" should have budget 50000 and actual 25000
