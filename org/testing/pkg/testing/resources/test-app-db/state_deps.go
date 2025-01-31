package testappdb

import (
	"context"
	"log/slog"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	xtestdb "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db"
)

type ScenarioState interface {
	GetScenarioInfo() models_common.ScenarioInfo
	GetLogger() *slog.Logger
	RegisterAppDatabase(ctx context.Context, serviceName models_service.ServiceName, dbInfo xtestdb.DatabaseConnection)
	GetScenarioAppDatabases(ctx context.Context) []xtestdb.DatabaseConnection
	GetServiceAppDatabase(ctx context.Context, serviceName models_service.ServiceName) (xtestdb.DatabaseConnection, error)
}
