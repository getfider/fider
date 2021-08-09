Feature: Post

  Scenario: Admin can create a post
    Given I go to the home page
    And I sign in as "admin"
    And I type "Feature Request Example" as the title
    And I type "This is just an example of a feature suggestion in fider" as the description
    And I click submit new post
    Then I should be on the show post page
    And I should see "Feature Request Example" as the post title
    And I should see 1 vote(s)