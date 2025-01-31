package xtesttime

import (
	"time"

	models_time "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
)

type ScenarioState interface {
	RegisterCheckpoint(checkpoint models_time.TimeCheckpointName, t models_time.TimeCheckpointValue) error
	GetCheckpointValue(checkpoint models_time.TimeCheckpointName) (models_time.TimeCheckpointValue, error)
	RegisterBackgroundOperation(operation models_time.BackgroundOperation) error
	RegisterBackgroundOperationComplete(operation models_time.BackgroundOperation) error
	GetBackgroundOperations() map[models_time.BackgroundOperationCompoundKey]models_time.BackgroundOperation
	ReportBackgroundError(err error)
	GetBackgroundErrors() []error
	WaitBackgroundsCompletionWithTimeout(timeout time.Duration) error
}
