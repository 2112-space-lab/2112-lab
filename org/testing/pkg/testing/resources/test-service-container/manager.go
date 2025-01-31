package testservicecontainer

import (
	"fmt"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	models_db "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

const AppServiceHttpPort models_cont.ContainerPort = 8090
const AppServiceGrpcPort models_cont.ContainerPort = 61007

const GemMockHttpPort models_cont.ContainerPort = 8087
const GemMockWssPort models_cont.ContainerPort = 8443
const GemMockGrpcPort models_cont.ContainerPort = 10001

type ServiceContainerManager struct {
	AppServiceDockerImage models_cont.DockerContainerImage
	GemMockDockerImage    models_cont.DockerContainerImage
	dbConf                models_db.DatabaseConnectionInfo
	networkName           models_cont.NetworkName
	testRunningEnv        models_common.TestRunningEnv
}

func NewServiceContainerManager(
	appServiceDockerImage models_cont.DockerContainerImage,
	gemMockDockerImage models_cont.DockerContainerImage,
	dbConf models_db.DatabaseConnectionInfo,
	networkName models_cont.NetworkName,
	testRunningEnv models_common.TestRunningEnv,
) *ServiceContainerManager {
	m := &ServiceContainerManager{
		AppServiceDockerImage: appServiceDockerImage,
		GemMockDockerImage:    gemMockDockerImage,
		dbConf:                dbConf,
		networkName:           networkName,
		testRunningEnv:        testRunningEnv,
	}
	return m
}

func prepareServiceContainerNames(scenarioState GatewwayContainerScenarioState, serviceName models_service.ServiceName, appName models_service.ServiceAppName) (models_cont.ContainerName, models_cont.ServiceName) {
	i := scenarioState.GetScenarioInfo()
	containerName := fmt.Sprintf("%s-%s-%s", i.ScenarioRandID, serviceName, appName)
	serviceNameApp := fmt.Sprintf("%s-%s", serviceName, appName)
	return models_cont.ContainerName(containerName), models_cont.ServiceName(serviceNameApp)
}
