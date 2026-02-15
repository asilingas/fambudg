Feature: Saving goals management

  Background:
    Given I am logged in as "admin@family.com"

  Scenario: Create a saving goal
    When I create a saving goal with:
      | field        | value          |
      | name         | Family Vacation|
      | targetAmount | 500000         |
      | priority     | 1              |
    Then the saving goal should be created successfully
    And the saving goal should have target amount 500000

  Scenario: List saving goals
    Given the following saving goals exist:
      | name            | targetAmount | priority |
      | Family Vacation | 500000       | 1        |
      | New Car         | 3000000      | 2        |
    When I list saving goals
    Then I should see 2 saving goals

  Scenario: Update a saving goal
    Given a saving goal "Emergency Fund" exists with target 1000000
    When I update the saving goal with name "Rainy Day Fund"
    Then the saving goal should be updated successfully
    And the saving goal name should be "Rainy Day Fund"

  Scenario: Contribute to a saving goal
    Given a saving goal "New Laptop" exists with target 200000
    When I contribute 50000 to the saving goal
    Then the saving goal current amount should be 50000

  Scenario: Auto-complete saving goal when target reached
    Given a saving goal "Small Goal" exists with target 10000
    When I contribute 10000 to the saving goal
    Then the saving goal status should be "completed"
