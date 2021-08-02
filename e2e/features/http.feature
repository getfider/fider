Feature: HTTP

  Scenario: Cache Control is set on valid resources
    Given I send a "GET" request to "/assets/assets.json"
    Then I should see http status 200
    And I should see a "Cache-Control" header with value "public, max-age=31536000"

  Scenario: Cache Control is not set on invalid resources
    Given I send a "GET" request to "/assets/invalid.js"
    Then I should see http status 404
    And I should see a "Cache-Control" header with value "no-cache, no-store"