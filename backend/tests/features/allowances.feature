Feature: Allowances for children

  Scenario: Admin sets allowance for a child
    Given I am logged in as "admin@family.com"
    And a child user "kid@family.com" exists
    When I create an allowance of 5000 for the child starting "2026-02-01"
    Then the allowance should be created successfully
    And the allowance amount should be 5000

  Scenario: Allowance spending is calculated from transactions
    Given I am logged in as "admin@family.com"
    And a child user "kid@family.com" exists with an account and a category
    And an allowance of 10000 for the child starting "2026-02-01"
    And the child has a transaction of -3000 on "2026-02-10"
    When I get allowances
    Then the allowance spent should be 3000
    And the allowance remaining should be 7000

  Scenario: Admin updates an allowance
    Given I am logged in as "admin@family.com"
    And a child user "kid@family.com" exists
    And an allowance of 5000 for the child starting "2026-02-01"
    When I update the allowance amount to 8000
    Then the allowance should be updated successfully
    And the allowance amount should be 8000

  Scenario: Admin lists all allowances
    Given I am logged in as "admin@family.com"
    And a child user "kid1@family.com" exists
    And an allowance of 5000 for the child starting "2026-02-01"
    And a child user "kid2@family.com" exists
    And an allowance of 3000 for the child starting "2026-02-01"
    When I get allowances
    Then I should see 2 allowances
