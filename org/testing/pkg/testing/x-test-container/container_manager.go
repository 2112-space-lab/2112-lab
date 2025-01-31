package xtestcontainer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	xtestlog "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-log"
	"github.com/testcontainers/testcontainers-go"
)

func SpawnContainerWithoutManagedLogConsumer(
	ctx context.Context,
	logger *slog.Logger,
	containerRequest testcontainers.ContainerRequest,
	testRunningEnv models_common.TestRunningEnv,
) (*BaseContainer, error) {
	dockerContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
		Logger:           &testContainerLogger{logger: logger},
	})
	if err != nil {
		return nil, err
	}
	cont := NewBaseContainer(dockerContainer, nil, testRunningEnv)
	return cont, nil

}

func SpawnContainerWithManagedLogConsumer(
	ctx context.Context,
	logger *slog.Logger,
	serviceName models_cont.ServiceName,
	containerRequest testcontainers.ContainerRequest,
	containerLogFolder string,
	testRunningEnv models_common.TestRunningEnv,
) (*BaseContainer, error) {
	var logConsumer *ContainerLogConsumer
	var err error

	logConsumer, err = NewContainerLogConsumer(containerLogFolder, serviceName)
	if err != nil {
		return nil, err
	}
	containerRequest.LogConsumerCfg = &testcontainers.LogConsumerConfig{
		Opts: []testcontainers.LogProductionOption{},
		Consumers: []testcontainers.LogConsumer{
			logConsumer,
		},
	}
	dockerContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
		Logger:           &testContainerLogger{logger: logger},
	})
	if err != nil {
		return nil, err
	}
	cont := NewBaseContainer(dockerContainer, logConsumer, testRunningEnv)
	return cont, nil
}

func TeardownContainer(ctx context.Context, logger *slog.Logger, container *BaseContainer) (err error) {
	containerName, errName := container.container.Name(ctx)
	if errName != nil {
		return errName
	}
	timeout := 5 * time.Second
	errStop := container.container.Stop(ctx, &timeout)
	errTeardownLog := container.logConsumer.Teardown(ctx)
	err = fx.FlattenErrorsIfAny(errStop, errTeardownLog)

	if err != nil {
		logger.Error("failure while teardown container",
			slog.String("containerName", containerName),
			slog.Any(xtestlog.AttrErrorKey, err),
		)
	}
	logger.Info("teardown container complete",
		slog.String("containerName", containerName),
		slog.Any(xtestlog.AttrErrorKey, err),
	)
	return fx.FlattenErrorsIfAny(errName, errStop)
}

func TeardownAllContainers(ctx context.Context, logger *slog.Logger, containers ...*BaseContainer) (err error) {
	errs := []error{}
	for _, container := range containers {
		err := TeardownContainer(ctx, logger, container)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return fx.FlattenErrorsIfAny(errs...)
}

type testContainerLogger struct {
	logger *slog.Logger
}

func (l *testContainerLogger) Printf(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}
