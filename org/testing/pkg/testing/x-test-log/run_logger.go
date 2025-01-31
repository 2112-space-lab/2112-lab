package xtestlog

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
)

var runLogger *slog.Logger
var runLoggerLock sync.Mutex

func GetRunLogger() *slog.Logger {
	runLoggerLock.Lock()
	defer runLoggerLock.Unlock()

	if runLogger == nil {
		panic("unexpected nil runLogger - must be initialized before starting tests")
	}
	return runLogger
}

func InitRunLoggerOrPanic(
	ctx context.Context,
	defaultLogger *slog.Logger,
	logLevelOption string,
	logFormat string,
	runRandID models_common.RunRandID,
	runFolderPath string,
) (*slog.Logger, models_common.ResourceTeardownSimpleFunc) {
	runLoggerLock.Lock()
	defer runLoggerLock.Unlock()

	if runLogger != nil {
		LogAndPanicIfError(runLogger, "invalid init sequence", ErrRunLoggerAlreadyInit)
	}

	level := slog.LevelDebug
	err := level.UnmarshalText([]byte(logLevelOption))
	if err != nil {
		level = slog.LevelDebug
		defaultLogger.Error("invalid log level - defaulting to debug", slog.Any("error", err))
	}
	logOptions := slog.HandlerOptions{
		Level: level,
	}
	isDefaultLogger := false
	runLogFile := filepath.Join(runFolderPath, "testing_run.log")
	fd, err := os.Create(runLogFile)
	if err != nil {
		isDefaultLogger = true
		runLogger = defaultLogger
		runLogger.Error("cannot create run log file - fallback to default logger",
			slog.Any(AttrErrorKey, err),
		)
	} else if logFormat == "json" {
		runLogger = slog.New(slog.NewJSONHandler(fd, &logOptions))
	} else {
		runLogger = slog.New(slog.NewTextHandler(fd, &logOptions))
	}
	runLogger = runLogger.With(slog.Group("testingCtx",
		slog.String("runRandID", string(runRandID)),
	))

	return runLogger, func(ctx context.Context) error {
		runLogger.Info("tearing down run logger")
		if isDefaultLogger {
			return nil
		}
		return fd.Close()
	}
}
