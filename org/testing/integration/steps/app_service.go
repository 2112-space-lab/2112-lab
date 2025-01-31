package steps

import (
	"context"

	"github.com/cucumber/godog"
	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
)

type AppServiceSteps struct {
	state                       appServiceStepsState
	appServiceContainerResource appServiceContainerResource
}

type appServiceContainerResource interface {
	SpawnServiceAppService(ctx context.Context, scenarioState testservicecontainer.GatewwayContainerScenarioState, serviceName models_service.ServiceName, serviceEnvOverrides models_cont.EnvVarKeyValueMap) error
}

type appServiceStepsState interface {
	testservicecontainer.GatewwayContainerScenarioState
}

func RegisterAppServiceSteps(ctx *godog.ScenarioContext, state appServiceStepsState, containerRsc appServiceContainerResource) {
	s := &AppServiceSteps{
		state:                       state,
		appServiceContainerResource: containerRsc,
	}
	ctx.Step(`^a App service is created for service "([^"]*)"$`, s.appServiceCreate)
	ctx.Step(`^a App service is created for service "([^"]*)" with env overrides:$`, s.appServiceCreateWithEnv)
}

func (steps *AppServiceSteps) appServiceCreate(ctx context.Context, serviceName string) error {
	err := steps.appServiceContainerResource.SpawnServiceAppService(ctx, steps.state, models_service.ServiceName(serviceName), models_cont.EnvVarKeyValueMap{})
	return err
}

func (steps *AppServiceSteps) appServiceCreateWithEnv(ctx context.Context, serviceName string, envVars *godog.Table) error {
	appEnvVars, err := GodogTableToKeyValueMap[string, string](envVars, true)
	if err != nil {
		return err
	}
	err = steps.appServiceContainerResource.SpawnServiceAppService(ctx, steps.state, models_service.ServiceName(serviceName), appEnvVars)
	return err
}
