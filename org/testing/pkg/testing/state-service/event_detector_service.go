package stateservice

import (
	"context"
	"fmt"
	"sync"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
)

type EventDetectorState struct {
	isInit                                    bool
	EventDetectorEnvScenarioOverrides         models_cont.EnvVarKeyValueMap
	EventDetectorGlobalPropsScenarioOverrides models_service.GlobalPropKeyValueMap
	EventDetectorServices                     map[models_service.ServiceName]*xtestcontainer.BaseContainer
	eventLock                                 sync.Mutex
	NamedEventReferences                      map[models_service.NamedEventReference]models_service.EventRawJSON
	CallbackResults                           map[string]interface{}
	V2StreamSubscribersCancel                 map[string]context.CancelFunc
}

func NewEventDetectorState() EventDetectorState {
	return EventDetectorState{
		isInit:                            true,
		EventDetectorEnvScenarioOverrides: models_cont.EnvVarKeyValueMap{},
		EventDetectorGlobalPropsScenarioOverrides: models_service.GlobalPropKeyValueMap{},
		EventDetectorServices:                     map[models_service.ServiceName]*xtestcontainer.BaseContainer{},
		CallbackResults:                           map[string]interface{}{},
		V2StreamSubscribersCancel:                 map[string]context.CancelFunc{},
		NamedEventReferences:                      map[models_service.NamedEventReference]models_service.EventRawJSON{},
	}
}

func (s *EventDetectorState) init() {
	if s.isInit {
		return
	}
	*s = NewEventDetectorState()
}

func (s *EventDetectorState) RegisterEventDetectorEnvScenarioOverrides(envVars models_cont.EnvVarKeyValueMap) {
	s.init()
	s.EventDetectorEnvScenarioOverrides = envVars
}

func (s *EventDetectorState) RegisterEventDetectorGlobalPropsScenarioOverrides(globalProps models_service.GlobalPropKeyValueMap) {
	s.init()
	s.EventDetectorGlobalPropsScenarioOverrides = globalProps
}

func (s *EventDetectorState) GetEventDetectorEnvScenarioOverrides() models_cont.EnvVarKeyValueMap {
	s.init()
	return s.EventDetectorEnvScenarioOverrides
}

func (s *EventDetectorState) GetScenarioEventDetectorServiceContainers() map[models_service.ServiceName]*xtestcontainer.BaseContainer {
	s.init()
	return s.EventDetectorServices
}

func (s *EventDetectorState) RegisterEventDetectorServiceContainer(ctx context.Context, serviceName models_service.ServiceName, container *xtestcontainer.BaseContainer) {
	s.init()
	s.EventDetectorServices[serviceName] = container
}

func (s *EventDetectorState) GetEventDetectorServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error) {
	s.init()
	if app, ok := s.EventDetectorServices[serviceName]; ok {
		return app, nil
	}
	return nil, fmt.Errorf("no EventDetector service container registered for service [%s]", serviceName)
}

func (s *EventDetectorState) RegisterEventDetectorNamedEventReference(ref models_service.NamedEventReference, jsonData models_service.EventRawJSON) {
	s.NamedEventReferences[ref] = jsonData
}

func (s *EventDetectorState) GetEventDetectorNamedEventByReference(ref models_service.NamedEventReference) (models_service.EventRawJSON, bool) {
	data, found := s.NamedEventReferences[ref]
	return data, found
}

func (s *EventDetectorState) RegisterCancelV2Stream(serviceName models_service.ServiceName, subscriber string, cancel context.CancelFunc) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	s.V2StreamSubscribersCancel[subscriber] = cancel
}

func (s *EventDetectorState) CancelAllV2Stream() {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	for _, c := range s.V2StreamSubscribersCancel {
		c()
	}
	s.V2StreamSubscribersCancel = map[string]context.CancelFunc{}
}
