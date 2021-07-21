Feature: Logout

  Background:
    Given I go to the signup page
    And I create a new account using my email
    And I activate my account
    Then I should be on the home page

  Scenario: User is logged in after signup
    Given I go to the home page
    Then I should be logged in

  Scenario: User can logout
    Given I go to the home page
    And I expand the user menu
    And I click on sign out
    Then I should be on the home page
    And I should not be logged in
