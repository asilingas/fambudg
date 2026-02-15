Feature: Reports and Dashboard

  Background:
    Given I am logged in as "admin@family.com"
    And a category "Groceries" of type "expense" exists
    And an account "Chase Checking" of type "checking" exists

  Scenario: Dashboard returns month summary and recent transactions
    Given the following transactions exist:
      | amount  | description      | date       |
      | 500000  | Monthly salary   | 2026-02-01 |
      | -15000  | Weekly groceries | 2026-02-05 |
      | -10000  | More groceries   | 2026-02-12 |
    When I get the dashboard for month 2 and year 2026
    Then the month summary should have income 500000 and expense 25000
    And I should see 3 recent transactions

  Scenario: Monthly income vs expense report
    Given the following transactions exist:
      | amount  | description    | date       |
      | 500000  | Salary         | 2026-02-01 |
      | -20000  | Groceries      | 2026-02-10 |
      | -5000   | Coffee         | 2026-02-15 |
    When I get the monthly report for month 2 and year 2026
    Then the monthly report should have income 500000 and expense 25000

  Scenario: Spending by category report
    Given the following transactions exist:
      | amount | description | date       |
      | -20000 | Groceries   | 2026-02-10 |
    And a category "Transport" of type "expense" exists
    And the following transactions exist:
      | amount | description | date       |
      | -10000 | Bus tickets | 2026-02-12 |
    When I get the category report for month 2 and year 2026
    Then I should see 2 category entries
    And the category "Groceries" should have total 20000
