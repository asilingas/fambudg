Feature: Bill reminders management

  Background:
    Given I am logged in as "admin@family.com"
    And a category "Utilities" of type "expense" exists
    And an account "Chase Checking" of type "checking" exists

  Scenario: Create a bill reminder
    When I create a bill reminder with:
      | field       | value       |
      | name        | Electricity |
      | amount      | 15000       |
      | dueDay      | 15          |
      | frequency   | monthly     |
      | nextDueDate | 2026-03-15  |
    Then the bill reminder should be created successfully
    And the bill reminder should have amount 15000

  Scenario: List bill reminders
    Given the following bill reminders exist:
      | name        | amount | dueDay | frequency | nextDueDate |
      | Electricity | 15000  | 15     | monthly   | 2026-03-15  |
      | Internet    | 7999   | 1      | monthly   | 2026-03-01  |
    When I list bill reminders
    Then I should see 2 bill reminders

  Scenario: Update a bill reminder
    Given a bill reminder "Water" exists with amount 5000
    When I update the bill reminder with amount 6000
    Then the bill reminder should be updated successfully
    And the bill reminder should have amount 6000

  Scenario: Delete a bill reminder
    Given a bill reminder "Phone" exists with amount 8000
    When I delete the bill reminder
    Then the bill reminder should be deleted successfully

  Scenario: Pay a bill reminder
    Given a bill reminder "Gas" exists with amount 12000 and account
    When I pay the bill reminder on "2026-02-15"
    Then the bill payment should create a transaction
    And the bill next due date should be advanced
