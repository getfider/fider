Feature: Prometheus

  Scenario: Prometheus metrics endpoint is accessible
    Given I send a "GET" request to "http://127.0.0.1:4000/metrics"
    Then I should see http status 200
    And I should see "TYPE fider_info gauge" on the response body