Feature: Role-based access control

  Scenario: Admin creates a member user
    Given I am logged in as "admin@family.com"
    When I create a user with email "wife@family.com" password "password123" name "Wife" and role "member"
    Then the user creation should succeed
    And the created user should have role "member"

  Scenario: Admin creates a child user
    Given I am logged in as "admin@family.com"
    When I create a user with email "kid@family.com" password "password123" name "Kid" and role "child"
    Then the user creation should succeed
    And the created user should have role "child"

  Scenario: Admin lists all users
    Given I am logged in as "admin@family.com"
    And a user exists with email "member@family.com" password "password123" name "Member" and role "member"
    When I list all users
    Then I should see 2 users

  Scenario: Admin updates a user role
    Given I am logged in as "admin@family.com"
    And a user exists with email "member@family.com" password "password123" name "Member" and role "member"
    When I update the user role to "child"
    Then the user update should succeed
    And the updated user should have role "child"

  Scenario: Admin deletes a user
    Given I am logged in as "admin@family.com"
    And a user exists with email "member@family.com" password "password123" name "Member" and role "member"
    When I delete the user
    Then the user deletion should succeed

  Scenario: Child can only see own transactions
    Given I am logged in as "admin@family.com"
    And a category "Groceries" exists
    And an account "Main Account" exists
    And a child user "kid@family.com" exists with an account and a transaction of -2000
    When I list transactions for the child user
    Then I should see 1 transactions

  Scenario: Shared vs personal expense filtering
    Given I am logged in as "admin@family.com"
    And a category "Groceries" exists
    And an account "Main Account" exists
    When I create a shared transaction with amount -3000
    And I create a personal transaction with amount -1500
    And I filter transactions by isShared "true"
    Then I should see 1 transactions
    And I filter transactions by isShared "false"
    Then I should see 1 transactions
