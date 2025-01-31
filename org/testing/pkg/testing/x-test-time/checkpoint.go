package xtesttime

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	models_time "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-time/models"
)

func EvaluateOptionalCheckpoint(s ScenarioState, expr models_time.TimeCheckpointExpression) (fx.Option[models_time.TimeCheckpointValue], error) {
	if expr == "" {
		return fx.NewEmptyOption[models_time.TimeCheckpointValue](), nil
	}
	t, err := EvaluateCheckpoint(s, expr)
	return fx.NewValueOption(t), err
}

func SleepUntilCheckpoint(s ScenarioState, logger *slog.Logger, expr models_time.TimeCheckpointExpression) error {
	if expr == "" {
		logger.Info("empty checkpoint expression for SleepUntilCheckpoint - no wait needed")
		return nil
	}
	cp, err := EvaluateCheckpoint(s, expr)
	if err != nil {
		return err
	}
	sleepDur := time.Until(time.Time(cp))
	if sleepDur < 0 {
		logger.Warn("checkpoint already past - sleep not required",
			slog.String("checkpointExpression", string(expr)),
			slog.Duration("durationUntilCheckpoint", sleepDur),
		)
		return nil
	}
	time.Sleep(sleepDur)
	return nil
}

func EvaluateCheckpoint(s ScenarioState, expr models_time.TimeCheckpointExpression) (models_time.TimeCheckpointValue, error) {
	baseCheckpoint := models_time.TimeCheckpointName("")
	offsetStr := ""
	outCheckpoint := models_time.TimeCheckpointName("")
	parts := strings.Split(string(expr), ">")
	if len(parts) < 2 || len(parts) > 3 {
		return models_time.TimeCheckpointValue(time.Now()),
			fmt.Errorf("invalid checkpoint expression [%s] - expected format [baseCheckpoint>+-offsetDuration] or [baseCheckpoint>+-offsetDuration>outputCheckpoint]", expr)
	}
	baseCheckpoint = models_time.TimeCheckpointName(parts[0])
	offsetStr = parts[1]
	if len(parts) == 3 {
		outCheckpoint = models_time.TimeCheckpointName(parts[2])
	}
	offset, err := time.ParseDuration(offsetStr)
	if err != nil {
		return models_time.TimeCheckpointValue(time.Now()), fmt.Errorf("invalid checkpoint expression [%s] - invalid offset duration format [%s] - [%w]", expr, offsetStr, err)
	}
	baseT, err := s.GetCheckpointValue(baseCheckpoint)
	if err != nil {
		return models_time.TimeCheckpointValue(time.Now()), err
	}
	outT := time.Time(baseT).Add(offset)
	if outCheckpoint != "" {
		err = s.RegisterCheckpoint(outCheckpoint, models_time.TimeCheckpointValue(outT))
		if err != nil {
			return models_time.TimeCheckpointValue(outT), err
		}
	}
	return models_time.TimeCheckpointValue(outT), nil
}
