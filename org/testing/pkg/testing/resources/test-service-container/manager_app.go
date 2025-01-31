package testservicecontainer

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestcontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func (m *ServiceContainerManager) SpawnServiceAppService(
	ctx context.Context,
	scenarioState GatewwayContainerScenarioState,
	serviceName models_service.ServiceName,
	serviceEnvOverrides models_cont.EnvVarKeyValueMap,
) error {
	appDB, err := scenarioState.GetServiceAppDatabase(ctx, serviceName)
	if err != nil {
		return err
	}
	appContName, appSvcName := prepareServiceContainerNames(scenarioState, serviceName, models_service.AppName)

	env := fx.MergeMapOverrides(
		GetAppDefaultEnv(),                         // defaults
		scenarioState.GetAppEnvScenarioOverrides(), // overrides for all services in scenario
		serviceEnvOverrides,                        // service specific overrides
		map[string]string{ // technical env vars
			"BIND_ADDRESS": AppServiceHttpPort.ToBindAddress(),
			"SERVER_TLS":   "false",

			"DB_NAME": string(appDB.Info.DatabaseName),
			"DB_HOST": string(m.dbConf.HostNameDocker),
			"DB_PORT": fmt.Sprintf("%d", m.dbConf.PortDocker),
			"DB_USER": m.dbConf.OwnerUser,
			"DB_PASS": m.dbConf.OwnerPassword,
		},
	)
	req := testcontainers.ContainerRequest{
		Image: string(m.AppServiceDockerImage),
		Name:  string(appContName),
		ExposedPorts: []string{
			AppServiceHttpPort.ToDockerTcpExposedPort(),
			AppServiceGrpcPort.ToDockerTcpExposedPort(),
		},
		Networks: []string{string(m.networkName)},
		WaitingFor: wait.ForAll(
			wait.ForLog("REHYDRATE"),
			wait.ForListeningPort(nat.Port(AppServiceHttpPort.ToDockerTcpExposedPort())),
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

	scenarioState.RegisterAppServiceContainer(ctx, serviceName, cont)
	return nil
}
