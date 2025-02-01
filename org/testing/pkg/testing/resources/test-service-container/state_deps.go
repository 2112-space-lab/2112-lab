package testservicecontainer

import (
	"context"
	"log/slog"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	xtestdb "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db"
	// "github.com/org/2112-space-lab/org/testing/testing/integration/models"
)

type GatewwayContainerScenarioState interface {
	GetScenarioInfo() models_common.ScenarioInfo
	GetScenarioFolder() string
	GetLogger() *slog.Logger

	GetServiceAppDatabase(ctx context.Context, serviceName models_service.ServiceName) (xtestdb.DatabaseConnection, error)
	RegisterAppServiceContainer(ctx context.Context, serviceName models_service.ServiceName, container *xtestcontainer.BaseContainer)
	GetAppServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error)

	GetAppEnvScenarioOverrides() models_cont.EnvVarKeyValueMap
	RegisterAppEnvScenarioOverrides(models_cont.EnvVarKeyValueMap)

	RegisterPropagatorServiceContainer(ctx context.Context, serviceName models_service.ServiceName, container *xtestcontainer.BaseContainer)
	GetPropagatorServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error)
}
