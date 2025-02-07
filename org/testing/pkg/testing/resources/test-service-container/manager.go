package testservicecontainer

import (
	"fmt"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	models_db "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

const AppServiceHttpPort models_cont.ContainerPort = 8090
const PropagatorServiceHttpPort models_cont.ContainerPort = 5000
const RabbitMQURL = "amqp://2112:2112@localhost:5672"

type ServiceContainerManager struct {
	AppServiceDockerImage        models_cont.DockerContainerImage
	PropagatorServiceDockerImage models_cont.DockerContainerImage
	dbConf                       models_db.DatabaseConnectionInfo
	networkName                  models_cont.NetworkName
	testRunningEnv               models_common.TestRunningEnv
}

func NewServiceContainerManager(
	appServiceDockerImage models_cont.DockerContainerImage,
	propagatorServiceDockerImage models_cont.DockerContainerImage,
	dbConf models_db.DatabaseConnectionInfo,
	networkName models_cont.NetworkName,
	testRunningEnv models_common.TestRunningEnv,
) *ServiceContainerManager {
	m := &ServiceContainerManager{
		AppServiceDockerImage:        appServiceDockerImage,
		PropagatorServiceDockerImage: propagatorServiceDockerImage,
		dbConf:                       dbConf,
		networkName:                  networkName,
		testRunningEnv:               testRunningEnv,
	}
	return m
}

func prepareServiceContainerNames(scenarioState GatewwayContainerScenarioState, serviceName models_service.ServiceName, appName models_service.ServiceAppName) (models_cont.ContainerName, models_cont.ServiceName) {
	i := scenarioState.GetScenarioInfo()
	containerName := fmt.Sprintf("%s-%s-%s", i.ScenarioRandID, serviceName, appName)
	serviceNameApp := fmt.Sprintf("%s-%s", serviceName, appName)
	return models_cont.ContainerName(containerName), models_cont.ServiceName(serviceNameApp)
}
