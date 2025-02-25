package stateservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
)

// PropagatorState stores the state of propagators
type PropagatorState struct {
	Propagators                  map[models.ServiceName]*xtestcontainer.BaseContainer
	isInit                       bool
	EnvScenarioOverrides         models_cont.EnvVarKeyValueMap
	GlobalPropsScenarioOverrides models_service.GlobalPropKeyValueMap
	eventLock                    sync.Mutex
	ReceivedEvents               map[models_service.ServiceName][]*models_service.EventRoot
	NamedEventReferences         map[models_service.NamedEventReference]models_service.EventRawJSON
	CallbackResults              map[string]interface{}
	StreamSubscribersCancel      map[string]context.CancelFunc
}

// NewPropagatorState initializes a new propagator state
func NewPropagatorState() PropagatorState {
	return PropagatorState{
		Propagators:                  make(map[models.ServiceName]*xtestcontainer.BaseContainer),
		isInit:                       true,
		EnvScenarioOverrides:         make(models_cont.EnvVarKeyValueMap),
		GlobalPropsScenarioOverrides: make(models_service.GlobalPropKeyValueMap),
		ReceivedEvents:               make(map[models_service.ServiceName][]*models_service.EventRoot),
		CallbackResults:              make(map[string]interface{}),
		StreamSubscribersCancel:      make(map[string]context.CancelFunc),
		NamedEventReferences:         make(map[models_service.NamedEventReference]models_service.EventRawJSON),
	}
}

// GetPropagatorServiceContainer retrieves a propagator service container
func (s *PropagatorState) GetPropagatorServiceContainer(ctx context.Context, serviceName models_service.ServiceName) (*xtestcontainer.BaseContainer, error) {
	if app, ok := s.Propagators[serviceName]; ok {
		return app, nil
	}
	return nil, fmt.Errorf("no Propagator service container registered for service [%s]", serviceName)
}

// RegisterPropagatorServiceContainer registers a propagator service container
func (s *PropagatorState) RegisterPropagatorServiceContainer(ctx context.Context, serviceName models.ServiceName, container *xtestcontainer.BaseContainer) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	s.Propagators[serviceName] = container
}

// SaveReceivedEventV2 stores received events in the state
func (s *PropagatorState) SaveReceivedEvent(event *models_service.EventRoot, serviceName models_service.ServiceName) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	s.ReceivedEvents[serviceName] = append(s.ReceivedEvents[serviceName], event)
}

// RegisterCallbackResult stores the callback result for a received event
func (s *PropagatorState) RegisterCallbackResult(event *models_service.EventRoot, actionName string, response interface{}, callbackErr error) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	resultKey := fmt.Sprintf("%s_%s", event.EventUid, actionName)

	if callbackErr != nil {
		s.CallbackResults[resultKey] = callbackErr
	} else {
		s.CallbackResults[resultKey] = response
	}
}

// GetCallbackResult retrieves callback results based on event occurrence
func (s *PropagatorState) GetCallbackResult(serviceName models_service.ServiceName, eventType string, occurrence int, actionName string) (*models_service.EventRoot, interface{}, error) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()

	occ := 0
	events, ok := s.ReceivedEvents[serviceName]
	if !ok {
		return nil, nil, fmt.Errorf("no events received for service [%s]", serviceName)
	}

	for _, evt := range events {
		if evt.EventType == eventType {
			occ++
			if occ == occurrence {
				resultKey := fmt.Sprintf("%s_%s", evt.EventUid, actionName)
				resCallback := s.CallbackResults[resultKey]
				if err, ok := resCallback.(error); ok {
					return evt, nil, err
				}
				return evt, resCallback, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("no event matching type [%s] with occurrence [%d] found for action [%s]", eventType, occurrence, actionName)
}

// GetGwReceivedEvents fetches received events in a given time range
func (s *PropagatorState) GetReceivedEvents(serviceName models_service.ServiceName, from time.Time, to time.Time) []models_service.EventRoot {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	res := []models_service.EventRoot{}
	events, ok := s.ReceivedEvents[serviceName]
	if !ok {
		return res
	}

	for _, evt := range events {
		if evt.GetEventTimeUtc().Inner().After(from) && evt.GetEventTimeUtc().Inner().Before(to) {
			res = append(res, *evt)
		}
	}
	return res
}

// RegisterNamedEventReference saves an event reference for later retrieval
func (s *PropagatorState) RegisterNamedEventReference(ref models_service.NamedEventReference, jsonData models_service.EventRawJSON) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	s.NamedEventReferences[ref] = jsonData
}

// GetNamedEventByReference retrieves a stored event by reference
func (s *PropagatorState) GetNamedEventByReference(ref models_service.NamedEventReference) (models_service.EventRawJSON, bool) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	data, found := s.NamedEventReferences[ref]
	return data, found
}

// RegisterCancelV2Stream registers a cancel function for event stream subscribers
func (s *PropagatorState) RegisterCancelV2Stream(serviceName models_service.ServiceName, subscriber string, cancel context.CancelFunc) {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	s.StreamSubscribersCancel[subscriber] = cancel
}

// CancelAllV2Stream cancels all event stream subscribers
func (s *PropagatorState) CancelAllV2Stream() {
	s.eventLock.Lock()
	defer s.eventLock.Unlock()
	for _, cancelFunc := range s.StreamSubscribersCancel {
		cancelFunc()
	}
	s.StreamSubscribersCancel = make(map[string]context.CancelFunc)
}
