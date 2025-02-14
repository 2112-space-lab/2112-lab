Feature: Propagator Service allows propagation of TLE from RabbitMQ event
  To support App Service testing, Propagator Service should be able to receive requests, propagate positions, store to Redis and emits events when computaiton is done

  Background: Propagator and Redis services are setup and running
    Given I wait "0s" and set now time as checkpoint "begin"
    And I register Propagator service default scenario environment variable overrides:
      | Key                   | Value                            |
      | LOG_LEVEL             | debug                            |
      | RABBITMQ_INPUT_QUEUE  | test-subscribe-propagator-input  |
      | RABBITMQ_OUTPUT_QUEUE | test-subscribe-propagator-output |
    And a Propagator service is created for service "TEST"
    And I register queues to RabbitMQ
      | Key                         | Value                       |
      | test-subscribe-propagator-input  | test-subscribe-propagator-input  |
      | test-subscribe-propagator-output | test-subscribe-propagator-output |
    And Propagator subscribes as consumer "TEST" for "test-subscribe-propagator-output" with registered callbacks:
      | EventType | Action | ActionHandlerArgs | ActionDelay |
    And I wait "5s" and set now time as checkpoint "ready"

  @Propagator
  Scenario: Propagator subscribes to propagation requests and emits propagation events
    Then I publish propagator events for service "TEST" on "test-subscribe-propagator-input" from file "./resources/tle_propagation_request.json"
    Then Propagator events are expected for service "TEST":
      | EventType                | Occurrence | IsReject | FromTime  | ToTimeWarn | ToTimeErr | ProduceCheckpointEventTime | AssignRef | XPathQuery | XPathValue |
      | SATELLITE_TLE_PROPAGATED | 1          |          | ready>+0s |            | ready>+2s |                            |           |            |            |
