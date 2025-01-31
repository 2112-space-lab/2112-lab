package xtestcontainer

import (
	"context"
	"strings"

	"github.com/docker/go-connections/nat"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	"github.com/testcontainers/testcontainers-go"
)

type BaseContainer struct {
	container      testcontainers.Container
	logConsumer    *ContainerLogConsumer
	testRunningEnv models_common.TestRunningEnv
}

func NewBaseContainer(
	container testcontainers.Container,
	logConsumer *ContainerLogConsumer,
	testRunningEnv models_common.TestRunningEnv,
) *BaseContainer {
	return &BaseContainer{
		container:      container,
		logConsumer:    logConsumer,
		testRunningEnv: testRunningEnv,
	}
}

func (c *BaseContainer) GetBoundPort(ctx context.Context, p models_cont.ContainerPort) (nat.Port, error) {
	if c.testRunningEnv != models_common.TestRunningOnHost {
		return nat.Port(p.ToString()), nil
	}
	boundPort, err := c.container.MappedPort(ctx, nat.Port(p.ToString()))
	return boundPort, err
}

func (c *BaseContainer) GetHostName(ctx context.Context) (string, error) {
	if c.testRunningEnv == models_common.TestRunningOnHost {
		return "localhost", nil
	}
	i, err := c.container.Inspect(ctx)
	if err != nil {
		return "", err
	}
	hostname := strings.TrimPrefix(i.Name, "/")
	return hostname, nil
}
