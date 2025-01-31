package xtestlog

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
)

func PrepareScenarioLogger(
	ctx context.Context,
	suiteLogger *slog.Logger,
	logLevelOption string,
	logFormat string,
	scenarioFolderPath string,
	scenarioInfo models_common.ScenarioInfo,
) (*slog.Logger, models_common.ResourceTeardownSimpleFunc) {
	isDefaultLogger := false

	logSuiteGroupAttr := slog.Any("scenarioInfo", scenarioInfo)
	defaultLogger := runLogger.With(logSuiteGroupAttr)

	level := slog.LevelDebug
	err := level.UnmarshalText([]byte(logLevelOption))
	if err != nil {
		level = slog.LevelDebug
		defaultLogger.Error("invalid log level - defaulting to debug", slog.Any("error", err))
	}
	logOptions := slog.HandlerOptions{
		Level: level,
	}

	var scenarioLogger *slog.Logger
	testLogPath := filepath.Join(scenarioFolderPath, "testing_scenario.log")
	fd, err := os.Create(testLogPath)
	if err != nil {
		isDefaultLogger = true
		scenarioLogger = defaultLogger
		scenarioLogger.Error("cannot create run log file - fallback to run logger",
			slog.Any(AttrErrorKey, err),
		)
	} else if logFormat == "json" {
		scenarioLogger = slog.New(slog.NewJSONHandler(fd, &logOptions).WithAttrs([]slog.Attr{logSuiteGroupAttr}))
	} else {
		scenarioLogger = slog.New(slog.NewTextHandler(fd, &logOptions).WithAttrs([]slog.Attr{logSuiteGroupAttr}))
	}
	return scenarioLogger, func(ctx context.Context) error {
		scenarioLogger.Info("tearing down scenario logger")
		if isDefaultLogger {
			return nil
		}
		return fd.Close()
	}
}
