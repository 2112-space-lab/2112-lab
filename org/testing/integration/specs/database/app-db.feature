Feature: Provision App database
  To run App properly, a database must be provisioned with migrations and seeds

  @DB
  Scenario: DB is migrated
    Given a App database is created for service "TEST"
    When I apply App database migrations for service "TEST"
    Then App database version for service "TEST" should be "2"

  @DB
  Scenario: DB is migrated and seeded
    Given a App database is created and migrated for service "TEST"
    When I apply App seed SQL "../assets/scripts/test_init.sql" for service "TEST":
      | Key | Value |
    Then App database table "config_schema.contexts" should not be empty for service "TEST"