package steps

import (
	"context"

	"github.com/cucumber/godog"
	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
)

type PropagatorServiceSteps struct {
	state                              propagatorServiceStepsState
	propagatorServiceContainerResource propagatorServiceContainerResource
}

type propagatorServiceContainerResource interface {
	SpawnServicePropagatorService(ctx context.Context, scenarioState testservicecontainer.GatewwayContainerScenarioState, serviceName models_service.ServiceName, serviceEnvOverrides models_cont.EnvVarKeyValueMap) error
}

type propagatorServiceStepsState interface {
	testservicecontainer.GatewwayContainerScenarioState
}

func RegisterPropagatorServiceSteps(ctx *godog.ScenarioContext, state propagatorServiceStepsState, containerRsc propagatorServiceContainerResource) {
	s := &PropagatorServiceSteps{
		state:                              state,
		propagatorServiceContainerResource: containerRsc,
	}
	ctx.Step(`^a Propagator service is created for service "([^"]*)"$`, s.propagatorServiceCreate)
	ctx.Step(`^a Propagator service is created for service "([^"]*)" with env overrides:$`, s.propagatorServiceCreateWithEnv)
	ctx.Step(`^I register Propagator service default scenario environment variable overrides:$`, s.propagatorRegisterCommonEnvVars)
}

func (steps *PropagatorServiceSteps) propagatorRegisterCommonEnvVars(envVars *godog.Table) error {
	gsaEnvVars, err := GodogTableToKeyValueMap[string, string](envVars, true)
	if err != nil {
		return err
	}
	steps.state.RegisterAppEnvScenarioOverrides(gsaEnvVars)
	return nil
}

func (steps *PropagatorServiceSteps) propagatorServiceCreate(ctx context.Context, serviceName string) error {
	err := steps.propagatorServiceContainerResource.SpawnServicePropagatorService(ctx, steps.state, models_service.ServiceName(serviceName), models_cont.EnvVarKeyValueMap{})
	return err
}

func (steps *PropagatorServiceSteps) propagatorServiceCreateWithEnv(ctx context.Context, serviceName string, envVars *godog.Table) error {
	propagatorEnvVars, err := GodogTableToKeyValueMap[string, string](envVars, true)
	if err != nil {
		return err
	}
	err = steps.propagatorServiceContainerResource.SpawnServicePropagatorService(ctx, steps.state, models_service.ServiceName(serviceName), propagatorEnvVars)
	return err
}
