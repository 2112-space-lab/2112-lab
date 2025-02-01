Feature: Propagator Service allows propagation of TLE
  To support App Service testing, Propgator Service should be able to receive requests

  @Propagator
  @WIP
  Scenario: Propagator service is running
    Given I wait "0s" and set now time as checkpoint "begin"
    And I register Propagator service default scenario environment variable overrides:
      | Key       | Value |
      | LOG_LEVEL | debug |
    And a Propagator service is created for service "TEST"
    And I wait "5s" and set now time as checkpoint "ready"

