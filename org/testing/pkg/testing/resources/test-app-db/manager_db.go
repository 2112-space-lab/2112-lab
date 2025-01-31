package testappdb

import (
	"context"
	"fmt"
	"path/filepath"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	xtestdb "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db"
	models_db "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"

	// import postgres driver
	_ "github.com/lib/pq"

	// used in testing for adding migrate file
	_ "github.com/golang-migrate/migrate/source/file"
)

func (m *AppDatabaseManager) CreateAppDatabase(
	ctx context.Context,
	scenarioState ScenarioState,
	serviceName models_service.ServiceName,
) (xtestdb.DatabaseConnection, error) {
	sc := scenarioState.GetScenarioInfo()
	pDbName := generateAppDatabaseName(sc.RunRandID, sc.ScenarioRandID, serviceName)
	dbConn, err := m.rootDatabaseManager.CreateDatabase(ctx, scenarioState.GetLogger(), models_db.PotentialDatabaseName(pDbName))
	if err != nil {
		return dbConn, err
	}
	scenarioState.RegisterAppDatabase(ctx, serviceName, dbConn)
	return dbConn, nil
}

func (m *AppDatabaseManager) MigrateAppUp(ctx context.Context, scenarioState ScenarioState, serviceName models_service.ServiceName) error {
	serviceDb, err := scenarioState.GetServiceAppDatabase(ctx, serviceName)
	if err != nil {
		return err
	}
	err = serviceDb.MigrateUp(ctx, scenarioState.GetLogger(), m.migrationsPath)
	return err
}

func (m *AppDatabaseManager) CreateAppDatabaseAndMigrateUp(ctx context.Context, scenarioState ScenarioState, serviceName models_service.ServiceName) (xtestdb.DatabaseConnection, error) {
	db, err := m.CreateAppDatabase(ctx, scenarioState, serviceName)
	if err != nil {
		return db, err
	}
	err = db.MigrateUp(ctx, scenarioState.GetLogger(), m.migrationsPath)
	return db, err
}

func (m *AppDatabaseManager) ExecuteAppSqlScript(ctx context.Context, scenarioState ScenarioState, serviceName models_service.ServiceName, sqlFile string, placeholders map[string]string) error {
	sqlFilePath := filepath.Join(m.appWorkspaceRootPath, sqlFile)
	serviceDb, err := scenarioState.GetServiceAppDatabase(ctx, serviceName)
	if err != nil {
		return err
	}
	err = serviceDb.ExecuteSqlScript(ctx, scenarioState.GetLogger(), sqlFilePath, placeholders)
	return err
}

func generateAppDatabaseName(runID models_common.RunRandID, scenarioRandID models_common.ScenarioRandID, serviceName models_service.ServiceName) models_db.DatabaseName {
	dbName := fmt.Sprintf("%s_%s_%s_%s",
		runID,
		scenarioRandID,
		models_service.AppName,
		serviceName,
	)
	return models_db.DatabaseName(dbName)
}
