package steps

import (
	"context"
	"log"

	"github.com/cucumber/godog"
	testservice "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service"
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
	testservice.AppClientScenarioState
}

func RegisterAppServiceSteps(ctx *godog.ScenarioContext, state appServiceStepsState, containerRsc appServiceContainerResource) {
	s := &AppServiceSteps{
		state:                       state,
		appServiceContainerResource: containerRsc,
	}
	ctx.Step(`^a App service is created for service "([^"]*)"$`, s.appServiceCreate)
	ctx.Step(`^a App service is created for service "([^"]*)" with env overrides:$`, s.appServiceCreateWithEnv)
	ctx.Step(`^I register App service default scenario environment variable overrides:$`, s.registerCommonEnvVars)
	ctx.Step(`^I subscribe as consumer "([^"]*)" with registered callbacks:$`, s.subscribeToEvents)
	ctx.Step(`^Events are expected for service "([^"]*)":$`, s.verifyEvents)
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

func (steps *AppServiceSteps) registerCommonEnvVars(envVars *godog.Table) error {
	env, err := GodogTableToKeyValueMap[string, string](envVars, true)
	if err != nil {
		return err
	}
	steps.state.RegisterAppEnvScenarioOverrides(env)
	return nil
}

func (steps *AppServiceSteps) subscribeToEvents(consumer string, eventTable *godog.Table) error {
	events, err := GodogTableToSlice[models_service.EventCallbackInfo](eventTable)
	if err != nil {
		log.Printf("Error parsing event subscription table: %v", err)
		return err
	}

	log.Printf("Subscribing as consumer '%s' with events: %+v", consumer, events)
	_, err = testservice.Subscribe(context.Background(), steps.state, models_service.ServiceName(consumer), events)
	return err
}

func (steps *AppServiceSteps) verifyEvents(serviceName string, eventTable *godog.Table) error {
	expectedEvents, err := GodogTableToSlice[models_service.ExpectedEvent](eventTable)
	if err != nil {
		log.Printf("Error parsing expected event table: %v", err)
		return err
	}

	log.Printf("Verifying expected events for service: %s", serviceName)
	return testservice.VerifyEvents(steps.state, serviceName, expectedEvents)
}
