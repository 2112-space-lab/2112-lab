package testservice

import (
	"context"
	"log/slog"
	"time"

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

type PropagatorClientScenarioState interface {
	GetScenarioInfo() models_common.ScenarioInfo
	GetScenarioFolder() string
	GetLogger() *slog.Logger

	xtesttime.ScenarioState

	GetPropagatorServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error)
	RegisterNamedEventReference(ref models_service.NamedAppEventReference, jsonData models_service.AppEventRawJSON)
	GetNamedEventByReference(ref models_service.NamedAppEventReference) (models_service.AppEventRawJSON, bool)
	GetReceivedEvents(serviceName models_service.ServiceName, from time.Time, to time.Time) []models_service.EventRoot
	SaveReceivedEvent(event *models_service.EventRoot, serviceName models_service.ServiceName)
}
