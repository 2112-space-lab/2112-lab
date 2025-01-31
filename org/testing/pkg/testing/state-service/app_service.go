package stateservice

import (
	"context"
	"fmt"
	"sync"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
)

type AppState struct {
	isInit                          bool
	AppEnvScenarioOverrides         models_cont.EnvVarKeyValueMap
	AppGlobalPropsScenarioOverrides models_service.GlobalPropKeyValueMap
	AppServices                     map[models_service.ServiceName]*xtestcontainer.BaseContainer
	eventLock                       sync.Mutex
	NamedEventReferences            map[models_service.NamedAppEventReference]models_service.AppEventRawJSON
	CallbackResults                 map[string]interface{}
	V2StreamSubscribersCancel       map[string]context.CancelFunc
}

func NewAppState() AppState {
	return AppState{
		isInit:                          true,
		AppEnvScenarioOverrides:         models_cont.EnvVarKeyValueMap{},
		AppGlobalPropsScenarioOverrides: models_service.GlobalPropKeyValueMap{},
		AppServices:                     map[models_service.ServiceName]*xtestcontainer.BaseContainer{},
		CallbackResults:                 map[string]interface{}{},
		V2StreamSubscribersCancel:       map[string]context.CancelFunc{},
		NamedEventReferences:            map[models_service.NamedAppEventReference]models_service.AppEventRawJSON{},
	}
}

func (s *AppState) init() {
	if s.isInit {
		return
	}
	*s = NewAppState()
}

func (s *AppState) RegisterAppEnvScenarioOverrides(envVars models_cont.EnvVarKeyValueMap) {
	s.init()
	s.AppEnvScenarioOverrides = envVars
}

func (s *AppState) RegisterAppGlobalPropsScenarioOverrides(globalProps models_service.GlobalPropKeyValueMap) {
	s.init()
	s.AppGlobalPropsScenarioOverrides = globalProps
}

func (s *AppState) GetAppEnvScenarioOverrides() models_cont.EnvVarKeyValueMap {
	s.init()
	return s.AppEnvScenarioOverrides
}

func (s *AppState) GetScenarioAppServiceContainers() map[models_service.ServiceName]*xtestcontainer.BaseContainer {
	s.init()
	return s.AppServices
}

func (s *AppState) RegisterAppServiceContainer(ctx context.Context, serviceName models_service.ServiceName, container *xtestcontainer.BaseContainer) {
	s.init()
	s.AppServices[serviceName] = container
}

func (s *AppState) GetAppServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error) {
	s.init()
	if app, ok := s.AppServices[serviceName]; ok {
		return app, nil
	}
	return nil, fmt.Errorf("no App service container registered for service [%s]", serviceName)
}

func (s *AppState) RegisterNamedEventReference(ref models_service.NamedAppEventReference, jsonData models_service.AppEventRawJSON) {
	s.NamedEventReferences[ref] = jsonData
}

func (s *AppState) GetNamedEventByReference(ref models_service.NamedAppEventReference) (models_service.AppEventRawJSON, bool) {
	data, found := s.NamedEventReferences[ref]
	return data, found
}

func (s *AppState) RegisterCancelV2Stream(serviceName models_service.ServiceName, subscriber string, cancel context.CancelFunc) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	s.V2StreamSubscribersCancel[subscriber] = cancel
}

func (s *AppState) CancelAllV2Stream() {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	for _, c := range s.V2StreamSubscribersCancel {
		c()
	}
	s.V2StreamSubscribersCancel = map[string]context.CancelFunc{}
}
