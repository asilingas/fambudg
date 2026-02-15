Feature: Transaction search

  Background:
    Given I am logged in as "admin@family.com"
    And a category "Groceries" of type "expense" exists
    And an account "Chase Checking" of type "checking" exists

  Scenario: Search transactions by description
    Given the following transactions exist:
      | amount | description      | date       |
      | -15000 | Weekly groceries | 2026-02-05 |
      | -10000 | Dinner out       | 2026-02-10 |
      | -5000  | Groceries store  | 2026-02-12 |
    When I search transactions with description "groceries"
    Then I should see 2 search results

  Scenario: Search transactions by date range
    Given the following transactions exist:
      | amount | description | date       |
      | -15000 | Item A      | 2026-01-15 |
      | -10000 | Item B      | 2026-02-10 |
      | -5000  | Item C      | 2026-02-12 |
    When I search transactions from "2026-02-01" to "2026-02-28"
    Then I should see 2 search results
