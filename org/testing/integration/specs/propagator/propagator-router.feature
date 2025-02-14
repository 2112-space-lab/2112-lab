Feature: Propagator Service allows propagation of TLE from HTTP request
  To support App Service testing, Propagator Service should be able to receive requests, propagate positions, store to Redis and emits events when computaiton is done

  Background: Propagator and Redis services are setup and running
    Given I wait "0s" and set now time as checkpoint "begin"
    And I register Propagator service default scenario environment variable overrides:
      | Key                   | Value                       |
      | LOG_LEVEL             | debug                       |
      | RABBITMQ_INPUT_QUEUE  | test-http-propagator-input  |
      | RABBITMQ_OUTPUT_QUEUE | test-http-propagator-output |
    And a Propagator service is created for service "TEST_PROPAGATOR"
    And Propagator subscribes as consumer "TEST_PROPAGATOR" for "test-http-propagator-output" with registered callbacks:
      | EventType | Action | ActionHandlerArgs | ActionDelay |
    And I wait "15s" and set now time as checkpoint "ready"

  @Propagator
  Scenario: Propagator generates satellite positions and store it into Redis
    Then I request satellite propagation on propagator for service "TEST_PROPAGATOR"
      | NoradId | TleLine1                                                              | TleLine2                                                              | StartTime            | DurationMinutes | IntervalSeconds |
      | 25544   | 1 25544U 98067A   24037.49236111  .00016717  00000-0  30709-3 0  9993 | 2 25544  51.6456 264.6625 0007377  24.9892 335.1377 15.50051545369613 | 2024-02-07T12:00:00Z | 60              | 10              |
    Then Propagator events are expected for service "TEST_PROPAGATOR":
      | EventType                | Occurrence | IsReject | FromTime  | ToTimeWarn | ToTimeErr  | ProduceCheckpointEventTime | AssignRef | XPathQuery | XPathValue |
      | SATELLITE_TLE_PROPAGATED | 1          |          | begin>+0s |            | ready>+15s |                            |           |            |            |