package testservice

import (
	"context"
	"log/slog"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	xtesttime "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time"
)

type AppClientScenarioState interface {
	GetScenarioInfo() models_common.ScenarioInfo
	GetScenarioFolder() string
	GetLogger() *slog.Logger

	xtesttime.ScenarioState

	GetAppServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error)
	GetScenarioAppServiceContainers() map[models_service.ServiceName]*xtestcontainer.BaseContainer
}

type GemMockClientScenarioState interface {
	GetScenarioInfo() models_common.ScenarioInfo
	GetScenarioFolder() string
	GetLogger() *slog.Logger

	GetGemMockContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error)
}
