package xtestlog

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
)

func PrepareSuiteLogger(
	ctx context.Context,
	runLogger *slog.Logger,
	logLevelOption string,
	logFormat string,
	runRandID models_common.RunRandID,
	suiteFolderPath string,
	suiteName string,
) (*slog.Logger, models_common.ResourceTeardownSimpleFunc) {
	isDefaultLogger := false

	logSuiteGroupAttr := slog.Group("testingCtx",
		slog.String("runRandID", string(runRandID)),
		slog.String("suiteName", suiteName),
	)
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

	var suiteLogger *slog.Logger
	testLogPath := filepath.Join(suiteFolderPath, "testing_suite.log")
	fd, err := os.Create(testLogPath)
	if err != nil {
		isDefaultLogger = true
		suiteLogger = defaultLogger
		suiteLogger.Error("cannot create suite log file - fallback to run logger",
			slog.Any(AttrErrorKey, err),
		)
	} else if logFormat == "json" {
		suiteLogger = slog.New(slog.NewJSONHandler(fd, &logOptions).WithAttrs([]slog.Attr{logSuiteGroupAttr}))
	} else {
		suiteLogger = slog.New(slog.NewTextHandler(fd, &logOptions).WithAttrs([]slog.Attr{logSuiteGroupAttr}))
	}
	return suiteLogger, func(ctx context.Context) error {
		suiteLogger.Info("tearing down suite logger")
		if isDefaultLogger {
			return nil
		}
		return fd.Close()
	}
}
