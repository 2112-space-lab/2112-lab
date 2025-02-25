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

	GetAppServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error)

	RabbitMqClientScenarioState
}

type PropagatorClientScenarioState interface {
	GetScenarioInfo() models_common.ScenarioInfo
	GetScenarioFolder() string
	GetPropagatorServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error)

	RabbitMqClientScenarioState
}

type RabbitMqClientScenarioState interface {
	xtesttime.ScenarioState
	GetLogger() *slog.Logger
	RegisterNamedEventReference(ref models_service.NamedEventReference, jsonData models_service.EventRawJSON)
	GetNamedEventByReference(ref models_service.NamedEventReference) (models_service.EventRawJSON, bool)
	GetReceivedEvents(serviceName models_service.ServiceName, from time.Time, to time.Time) []models_service.EventRoot
	SaveReceivedEvent(event *models_service.EventRoot, serviceName models_service.ServiceName)
}
