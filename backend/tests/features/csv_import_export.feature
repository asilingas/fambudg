Feature: CSV Import and Export
  As a user
  I want to import and export transactions as CSV
  So that I can manage my data in spreadsheets

  Background:
    Given I am logged in as "admin@family.com"
    And a category "Groceries" of type "expense" exists
    And an account "Chase Checking" of type "checking" exists

  Scenario: Export transactions to CSV
    Given the following transactions exist:
      | amount | description    | date       |
      | -5000  | Weekly shop    | 2026-01-10 |
      | -3000  | Coffee beans   | 2026-01-15 |
    When I export transactions as CSV
    Then the CSV should contain 2 data rows
    And the CSV header should include "date,amount,type,description"

  Scenario: Import transactions from CSV
    When I import the following CSV:
      | date       | amount | type    | description  | is_shared |
      | 2026-01-20 | -4500  | expense | Dinner out   | true      |
      | 2026-01-21 | -1200  | expense | Bus fare     | false     |
    Then 2 transactions should be imported
    When I list transactions from "2026-01-01" to "2026-12-31"
    Then I should see 2 transactions
