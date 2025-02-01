package testservicecontainer

import (
	"context"

	"github.com/docker/go-connections/nat"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func (m *ServiceContainerManager) SpawnServicePropagatorService(
	ctx context.Context,
	scenarioState GatewwayContainerScenarioState,
	serviceName models_service.ServiceName,
	serviceEnvOverrides models_cont.EnvVarKeyValueMap,
) error {
	appContName, appSvcName := prepareServiceContainerNames(scenarioState, serviceName, models_service.PropagatorAppName)

	env := fx.MergeMapOverrides(
		GetAppDefaultEnv(),                         // defaults
		scenarioState.GetAppEnvScenarioOverrides(), // overrides for all services in scenario
		serviceEnvOverrides,                        // service specific overrides
		map[string]string{ // technical env vars
			"BIND_ADDRESS": PropagatorServiceHttpPort.ToBindAddress(),
		},
	)
	req := testcontainers.ContainerRequest{
		Image: string(m.PropagatorServiceDockerImage),
		Name:  string(appContName),
		ExposedPorts: []string{
			PropagatorServiceHttpPort.ToDockerTcpExposedPort(),
		},
		Networks: []string{string(m.networkName)},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(PropagatorServiceHttpPort.ToDockerTcpExposedPort())),
		),
		Env: env,
	}

	cont, err := xtestcontainer.SpawnContainerWithManagedLogConsumer(
		ctx,
		scenarioState.GetLogger(),
		appSvcName,
		req,
		scenarioState.GetScenarioFolder(),
		m.testRunningEnv,
	)
	if err != nil {
		return err
	}

	scenarioState.RegisterPropagatorServiceContainer(ctx, serviceName, cont)
	return nil
}
