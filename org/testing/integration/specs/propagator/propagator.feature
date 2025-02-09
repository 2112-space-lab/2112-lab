Feature: Propagator Service allows propagation of TLE
  To support App Service testing, Propgator Service should be able to receive requests, propagate positions, store to Redis and emits events when computaiton is done

  Background: Propagator and Redis services are setup and running
    Given I wait "0s" and set now time as checkpoint "begin"
    And I register Propagator service default scenario environment variable overrides:
      | Key       | Value |
      | LOG_LEVEL | debug |
    And a Propagator service is created for service "TEST"
    And I subscribe as consumer "TEST" with registered callbacks:
      | EventType | Action | ActionHandlerArgs | ActionDelay |
    And I wait "5s" and set now time as checkpoint "ready"

  @Propagator
  Scenario: Propagator generates satellite positions and store it into Redis
    Then I request satellite propagation on propagator for service "TEST"
      | NoradId | TleLine1                                                              | TleLine2                                                              | StartTime            | DurationMinutes | IntervalSeconds |
      | 25544   | 1 25544U 98067A   24037.49236111  .00016717  00000-0  30709-3 0  9993 | 2 25544  51.6456 264.6625 0007377  24.9892 335.1377 15.50051545369613 | 2024-02-07T12:00:00Z | 60              | 10              |
    Then Propagator events are expected for service "TEST":
      | EventType                | Occurence | IsReject | FromTime  | ToTimeWarn | ToTimeErr | ProduceCheckpointEventTime | AssignRef | XPathQuery | XPathValue |
      | SATELLITE_TLE_PROPAGATED | 1         |          | ready>+0s |            | ready>+5s |                            |           |            |            |

