package steps

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
	testappdb "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-app-db"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestdb "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db"
)

type AppDatabaseSteps struct {
	state             databaseStepsState
	databaseResources databaseStepsDatabaseResource
}

type databaseStepsDatabaseResource interface {
	CreateAppDatabase(ctx context.Context, scenarioState testappdb.ScenarioState, serviceName models_service.ServiceName) (xtestdb.DatabaseConnection, error)
	CreateAppDatabaseAndMigrateUp(ctx context.Context, scenarioState testappdb.ScenarioState, serviceName models_service.ServiceName) (xtestdb.DatabaseConnection, error)
	MigrateAppUp(ctx context.Context, scenarioState testappdb.ScenarioState, serviceName models_service.ServiceName) error
	ExecuteAppSqlScript(ctx context.Context, scenarioState testappdb.ScenarioState, serviceName models_service.ServiceName, sqlFile string, placeholders map[string]string) error
}

type databaseStepsState interface {
	testappdb.ScenarioState
}

func RegisterAppDatabaseSteps(ctx *godog.ScenarioContext, state databaseStepsState, databaseResources databaseStepsDatabaseResource) {
	dbs := &AppDatabaseSteps{
		databaseResources: databaseResources,
		state:             state,
	}
	ctx.Step(`^a App database is created for service "([^"]*)"$`, dbs.dbAppCreate)
	ctx.Step(`^I apply App database migrations for service "([^"]*)"$`, dbs.dbAppMigrateUp)
	ctx.Step(`^App database version for service "([^"]*)" should be "(\d+)"$`, dbs.dbAppCheckMigrationVersion)

	ctx.Step(`^a App database is created and migrated for service "([^"]*)"$`, dbs.dbAppCreateAndMigrateUp)
	ctx.Step(`^App database table "([^"]*)" should not be empty for service "([^"]*)"$`, dbs.dbAppCheckTableNotEmpty)
	ctx.Step(`^I apply App seed SQL "([^"]*)" for service "([^"]*)":$`, dbs.dbAppApplySeed)
}

func (dbs *AppDatabaseSteps) dbAppCreate(ctx context.Context, serviceName string) error {
	_, err := dbs.databaseResources.CreateAppDatabase(ctx, dbs.state, models_service.ServiceName(serviceName))
	return err
}

func (dbs *AppDatabaseSteps) dbAppMigrateUp(ctx context.Context, serviceName string) error {
	err := dbs.databaseResources.MigrateAppUp(ctx, dbs.state, models_service.ServiceName(serviceName))
	return err
}

func (dbs *AppDatabaseSteps) dbAppCheckMigrationVersion(ctx context.Context, serviceName string, dbVersion int) error {
	dbs.state.GetLogger().Info("AppDatabaseCheckVersion")
	db, err := dbs.state.GetServiceAppDatabase(ctx, models_service.ServiceName(serviceName))
	if err != nil {
		return err
	}
	conn, err := db.ConnPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	r := conn.QueryRow(ctx, "SELECT version, dirty FROM public.schema_migrations;")
	v := 0
	isDirty := false
	err = r.Scan(&v, &isDirty)
	if err != nil {
		return err
	}
	if isDirty {
		return fmt.Errorf("migration state is dirty")
	}
	if v != dbVersion {
		return fmt.Errorf("invalid DB migration version [%d] - expected [%d]", v, dbVersion)
	}
	return nil
}

func (dbs *AppDatabaseSteps) dbAppCreateAndMigrateUp(ctx context.Context, serviceName string) error {
	_, err := dbs.databaseResources.CreateAppDatabaseAndMigrateUp(ctx, dbs.state, models_service.ServiceName(serviceName))
	return err
}

func (dbs *AppDatabaseSteps) dbAppCheckTableNotEmpty(ctx context.Context, tableName string, serviceName string) error {
	db, err := dbs.state.GetServiceAppDatabase(ctx, models_service.ServiceName(serviceName))
	if err != nil {
		return err
	}
	conn, err := db.ConnPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	r := conn.QueryRow(ctx, "SELECT count(*) FROM "+tableName+" LIMIT 1;")
	var val int
	err = r.Scan(&val)
	return err
}

func (dbs *AppDatabaseSteps) dbAppApplySeed(ctx context.Context, sqlFile string, serviceName string, keyValue *godog.Table) error {
	placeholders, err := GodogTableToKeyValueMap[string, string](keyValue, true)
	if err != nil {
		return err
	}
	err = dbs.databaseResources.ExecuteAppSqlScript(ctx, dbs.state, models_service.ServiceName(serviceName), sqlFile, placeholders)
	return err
}
