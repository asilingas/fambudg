Feature: Spending trends

  Background:
    Given I am logged in as "admin@family.com"
    And a category "Groceries" of type "expense" exists
    And an account "Chase Checking" of type "checking" exists

  Scenario: Get month-over-month trends
    Given the following transactions exist:
      | amount | description      | date       |
      | -5000  | January groceries | 2026-01-15 |
      | -3000  | February groceries| 2026-02-10 |
    When I get trends for the last 6 months
    Then I should see trend data
    And the trend data should contain entries for the transaction months
