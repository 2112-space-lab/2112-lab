package models

import (
	"context"
	"log/slog"
)

type ResourceTeardownLoggedFunc func(ctx context.Context, logger *slog.Logger) error
type ResourceTeardownSimpleFunc func(ctx context.Context) error

type TestRunningEnv string

const (
	TestRunningOnHost   TestRunningEnv = "host"
	TestRunningOnDocker TestRunningEnv = "docker"
)
