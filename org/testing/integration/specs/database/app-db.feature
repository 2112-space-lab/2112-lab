Feature: Provision App database
  To run App properly, a database must be provisioned with migrations and seeds

  @DB
  Scenario: DB is migrated
    Given a App database is created for gateway "TEST"
    When I apply App database migrations for gateway "TEST"
    Then App database version for gateway "TEST" should be "1"

  @DB
  Scenario: DB is migrated and seeded
    Given a App database is created and migrated for gateway "TEST"
    When I apply App seed SQL "assets/scripts/init.sql" for gateway "TEST":
      | Key  | Value   |
      | Mode | Offline |
    Then App database table "config_schema.inventory" should not be empty for gateway "TEST"