package xtesttime

import (
	"context"
	"time"

	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
)

func DeferredExecution(
	ctx context.Context,
	scenarioState ScenarioState,
	checkpoint models.TimeCheckpointExpression,
	operationKey string,
	operation func(context.Context) error,
) error {
	if checkpoint == "" {
		err := operation(ctx)
		return err
	}
	cpValue, err := EvaluateCheckpoint(scenarioState, checkpoint)
	if err != nil {
		return err
	}
	op := models.NewBackgroundOperation(operationKey, checkpoint)
	err = scenarioState.RegisterBackgroundOperation(*op)
	if err != nil {
		return err
	}
	go func() {
		wait := time.Until(time.Time(cpValue))
		time.Sleep(wait)
		op.Start()
		err := operation(context.Background())
		op.Complete(err)
		if err != nil {
			scenarioState.ReportBackgroundError(err)
		}
		err = scenarioState.RegisterBackgroundOperationComplete(*op)
		if err != nil {
			scenarioState.ReportBackgroundError(err)
		}
	}()
	return nil
}
