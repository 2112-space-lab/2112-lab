Feature: App Service registers mapping of satellite at each new propagation
  App Service registers mapping of satellite at each new propagation

  Background: App and Redis services are setup and running
    Given I wait "0s" and set now time as checkpoint "begin"
    And a App database is created and migrated for service "TEST"
    And I register App service default scenario environment variable overrides:
      | Key       | Value |
      | LOG_LEVEL | debug |
    And a App service is created for service "TEST"
    And I subscribe as consumer "TEST" with registered callbacks:
      | EventType | Action | ActionHandlerArgs | ActionDelay |
    And I wait "5s" and set now time as checkpoint "ready"

  @App
  @WIP
  Scenario: App registers mapping of satellite at each new propagation