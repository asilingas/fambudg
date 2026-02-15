Feature: Account transfers

  Background:
    Given I am logged in as "admin@family.com"
    And a category "General" of type "expense" exists
    And an account "Chase Checking" of type "checking" exists
    And a second account "Savings Account" of type "savings" exists

  Scenario: Transfer between accounts
    When I transfer 50000 from "Chase Checking" to "Savings Account" with description "Monthly savings"
    Then the transfer should be created successfully
    And the transfer transaction should have amount -50000

  Scenario: Cannot transfer to same account
    When I transfer 10000 from "Chase Checking" to "Chase Checking" with description "Invalid"
    Then the transfer should fail
