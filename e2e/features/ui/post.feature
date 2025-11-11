Feature: Post

  Scenario: Admin can create a post
    Given I go to the home page
    And I sign in as "admin"
    And I click enter your suggestion
    And I type "This is just an example of a feature suggestion in fider" as the description
    And I click submit your feedback
    Then I should be on the show post page
    And I should see "This is just an example of a feature suggestion in fider" as the post title
    And I should see 1 vote(s)

  Scenario: Non-logged in user can view a post
    Given I go to the home page
    And I search for "Feature Request Example"
    And I click on the first post
    Then I should be on the show post page
    And I should see "This is just an example of a feature suggestion in fider" as the post title
    And I should see 1 vote(s)

  Scenario: Non-logged in user can draft a post and submit once signed up
    Given I go to the home page
    And I click enter your suggestion
    And I type "This is a draft post from a new user" as the description
    And I type my email address
    And I click continue with email
    Then I should see the name field
    Given I enter my name as "Matt"
    And I click continue
    Then I should be on the confirmation code page
    Given I enter the confirmation code
    Then I should be on the home page
    And I should see the new post modal
    And I should see "This is a draft post from a new user" as the draft post title