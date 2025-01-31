package xteststate

import (
	"fmt"
	"sync"
	"time"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
)

type TimeCheckpointState struct {
	isInit                   bool
	checkpoints              map[models.TimeCheckpointName]models.TimeCheckpointValue
	backgroundOperations     map[models.BackgroundOperationCompoundKey]models.BackgroundOperation
	backgroundErrors         []error
	backgroundOperationsLock sync.Mutex
	backgroundOperationsWg   sync.WaitGroup
}

func NewTimeCheckpointState() TimeCheckpointState {
	return TimeCheckpointState{
		isInit:               true,
		checkpoints:          map[models.TimeCheckpointName]models.TimeCheckpointValue{},
		backgroundOperations: map[models.BackgroundOperationCompoundKey]models.BackgroundOperation{},
	}
}

func (s *TimeCheckpointState) RegisterCheckpoint(checkpoint models.TimeCheckpointName, t models.TimeCheckpointValue) error {
	if tt, ok := s.checkpoints[checkpoint]; ok {
		return fmt.Errorf("time checkpoint [%s] cannot be registered at [%s] - already registered at [%s]",
			checkpoint,
			time.Time(t).Format(time.RFC3339),
			time.Time(tt).Format(time.RFC3339),
		)
	}
	s.checkpoints[checkpoint] = t
	return nil
}

func (s *TimeCheckpointState) GetCheckpointValue(checkpoint models.TimeCheckpointName) (models.TimeCheckpointValue, error) {
	if tt, ok := s.checkpoints[checkpoint]; ok {
		return models.TimeCheckpointValue(tt), nil
	}
	return models.TimeCheckpointValue(time.Now()), fmt.Errorf("unknown checkpoint [%s]", checkpoint)
}

func (s *TimeCheckpointState) RegisterBackgroundOperation(operation models.BackgroundOperation) error {
	s.backgroundOperationsLock.Lock()
	defer s.backgroundOperationsLock.Unlock()
	if op, found := s.backgroundOperations[operation.GetKey()]; found {
		return fmt.Errorf("cannot register operation [%+v] - same key already used for another operation [%+v]", operation, op)
	}
	s.backgroundOperationsWg.Add(1)
	s.backgroundOperations[operation.GetKey()] = operation
	return nil
}

func (s *TimeCheckpointState) RegisterBackgroundOperationComplete(operation models.BackgroundOperation) error {
	s.backgroundOperationsLock.Lock()
	defer s.backgroundOperationsLock.Unlock()
	if _, found := s.backgroundOperations[operation.GetKey()]; !found {
		return fmt.Errorf("cannot up operation [%+v] - key not found", operation)
	}
	s.backgroundOperations[operation.GetKey()] = operation
	s.backgroundOperationsWg.Done()
	return nil
}

func (s *TimeCheckpointState) GetBackgroundOperations() map[models.BackgroundOperationCompoundKey]models.BackgroundOperation {
	s.backgroundOperationsLock.Lock()
	defer s.backgroundOperationsLock.Unlock()
	v := fx.Values(s.backgroundOperations)
	snapshot := fx.ToMapByKeySelector(v, models.BackgroundOperation.GetKey)
	return snapshot
}

func (s *TimeCheckpointState) ReportBackgroundError(err error) {
	s.backgroundOperationsLock.Lock()
	defer s.backgroundOperationsLock.Unlock()
	s.backgroundErrors = append(s.backgroundErrors, err)
}

func (s *TimeCheckpointState) GetBackgroundErrors() []error {
	s.backgroundOperationsLock.Lock()
	defer s.backgroundOperationsLock.Unlock()
	return s.backgroundErrors
}

func (s *TimeCheckpointState) WaitBackgroundsCompletionWithTimeout(timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		s.backgroundOperationsWg.Wait()
	}()
	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout while waiting for background operations to complete")
	}
}
