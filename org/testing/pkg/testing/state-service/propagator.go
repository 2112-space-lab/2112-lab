package stateservice

import (
	"context"
	"fmt"

	"github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
)

type PropagatorState struct {
	Propagators map[models.ServiceName]*xtestcontainer.BaseContainer
}

func NewPropagatorState() PropagatorState {
	return PropagatorState{
		Propagators: map[models.ServiceName]*xtestcontainer.BaseContainer{},
	}
}

func (s *PropagatorState) GetPropagatorServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error) {
	if app, ok := s.Propagators[serviceName]; ok {
		return app, nil
	}
	return nil, fmt.Errorf("no Propagator service container registered for service [%s]", serviceName)
}

func (s *PropagatorState) GetScenarioPropagatorContainers() map[models.ServiceName]*xtestcontainer.BaseContainer {
	return s.Propagators
}

func (s *PropagatorState) RegisterPropagatorServiceContainer(ctx context.Context, serviceName models.ServiceName, container *xtestcontainer.BaseContainer) {
	s.Propagators[serviceName] = container
}
