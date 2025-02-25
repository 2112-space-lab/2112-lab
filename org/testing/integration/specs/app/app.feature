Feature: App Service registers mapping of satellite at each new propagation
  App Service registers mapping of satellite at each new propagation

  Background: App and Redis services are setup and running
    Given I wait "0s" and set now time as checkpoint "begin"
    And a App database is created and migrated for service "TEST"
    And I register App service default scenario environment variable overrides:
      | Key       | Value |
      | LOG_LEVEL | debug |
    And a App service is created for service "TEST"
    And I register queues to RabbitMQ
      | Key             | Value           |
      | test-app-input  | test-app-input  |
      | test-app-output | test-app-output |
    And I subscribe as consumer "TEST" for "test-app-output" with registered callbacks:
      | EventType | Action | ActionHandlerArgs | ActionDelay |
    And I wait "5s" and set now time as checkpoint "ready"
    And I apply App seed SQL "./assets/scripts/test_init.sql" for service "TEST":
      | Key | Value |

  @App
  @WIP
  Scenario: App registers mapping of satellite at each new propagation
    Given I create for service "TEST" a game context "TEST_CONTEXT" with the following satellites:
      | SatelliteName |
      | SAT1          |
      | SAT2          |
    When I rehydrate for service "TEST" the game context "TEST_CONTEXT"
    Then App events are expected for service "TEST":
      | EventType                | Occurrence | IsReject | FromTime  | ToTimeWarn | ToTimeErr | ProduceCheckpointEventTime | AssignRef | XPathQuery | XPathValue |
      | SATELLITE_TLE_PROPAGATED | 1          |          | ready>+0s |            | ready>+2s |                            |           |            |            |