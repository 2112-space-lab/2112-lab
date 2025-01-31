package stateservice

import (
	"context"
	"fmt"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestdb "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db"
)

type AppDbState struct {
	AppDatabases map[models_service.ServiceName]xtestdb.DatabaseConnection
}

func NewAppDBState() AppDbState {
	return AppDbState{
		AppDatabases: map[models_service.ServiceName]xtestdb.DatabaseConnection{},
	}
}

func (s *AppDbState) RegisterAppDatabase(ctx context.Context, serviceName models_service.ServiceName, dbInfo xtestdb.DatabaseConnection) {
	s.AppDatabases[serviceName] = dbInfo
}

func (s *AppDbState) GetScenarioAppDatabases(ctx context.Context) []xtestdb.DatabaseConnection {
	return fx.Values(s.AppDatabases)
}

func (s *AppDbState) GetServiceAppDatabase(ctx context.Context, serviceName models_service.ServiceName) (xtestdb.DatabaseConnection, error) {
	if db, ok := s.AppDatabases[serviceName]; ok {
		return db, nil
	}
	return xtestdb.DatabaseConnection{}, fmt.Errorf("no App database registered for service [%s]", serviceName)
}
