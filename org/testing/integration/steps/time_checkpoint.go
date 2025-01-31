package steps

import (
	"context"
	"log/slog"
	"time"

	"github.com/cucumber/godog"
	xtesttime "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time"
	models_time "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
)

type TimeCheckpointSteps struct {
	state timeCheckpointStepsState
}

type timeCheckpointStepsState interface {
	GetLogger() *slog.Logger
	xtesttime.ScenarioState
}

func RegisterTimeCheckpointSteps(ctx *godog.ScenarioContext, state timeCheckpointStepsState) {
	steps := &TimeCheckpointSteps{
		state: state,
	}
	ctx.Step(`^I wait "([^"]*)" and set now time as checkpoint "([^"]*)"$`, steps.timeWaitAndSetCheckpoint)
	ctx.Step(`^I wait until "([^"]*)"`, steps.sleepUntilCheckpoint)
}

func (steps *TimeCheckpointSteps) timeWaitAndSetCheckpoint(ctx context.Context, waitDur string, checkpointName string) error {
	d, err := time.ParseDuration(waitDur)
	if err != nil {
		return err
	}
	if d > 0 {
		time.Sleep(d)
	}
	err = steps.state.RegisterCheckpoint(models_time.TimeCheckpointName(checkpointName), models_time.TimeCheckpointValue(time.Now().UTC()))
	return err
}

func (steps *TimeCheckpointSteps) sleepUntilCheckpoint(ctx context.Context, checkpointExpression string) error {
	err := xtesttime.SleepUntilCheckpoint(steps.state, steps.state.GetLogger(), models_time.TimeCheckpointExpression(checkpointExpression))
	return err
}
