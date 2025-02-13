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

type PropagatorServiceSteps struct {
	state                              propagatorServiceStepsState
	propagatorServiceContainerResource propagatorServiceContainerResource
}

type propagatorServiceContainerResource interface {
	SpawnServicePropagatorService(ctx context.Context, scenarioState testservicecontainer.GatewwayContainerScenarioState, serviceName models_service.ServiceName, serviceEnvOverrides models_cont.EnvVarKeyValueMap) error
}

type propagatorServiceStepsState interface {
	testservicecontainer.GatewwayContainerScenarioState
	testservice.PropagatorClientScenarioState
}

func RegisterPropagatorServiceSteps(ctx *godog.ScenarioContext, state propagatorServiceStepsState, containerRsc propagatorServiceContainerResource) {
	s := &PropagatorServiceSteps{
		state:                              state,
		propagatorServiceContainerResource: containerRsc,
	}

	ctx.Step(`^a Propagator service is created for service "([^"]*)"$`, s.propagatorServiceCreate)
	ctx.Step(`^a Propagator service is created for service "([^"]*)" with env overrides:$`, s.propagatorServiceCreateWithEnv)
	ctx.Step(`^I register Propagator service default scenario environment variable overrides:$`, s.propagatorRegisterCommonEnvVars)
	ctx.Step(`^I request satellite propagation on propagator for service "([^"]*)"$`, s.requestPropagation)
	ctx.Step(`^I subscribe as consumer "([^"]*)" with registered callbacks:$`, s.subscribeToEvents)
	ctx.Step(`^Events are expected for service "([^"]*)":$`, s.verifyEvents)
}

func (steps *PropagatorServiceSteps) propagatorRegisterCommonEnvVars(envVars *godog.Table) error {
	env, err := GodogTableToKeyValueMap[string, string](envVars, true)
	if err != nil {
		return err
	}
	steps.state.RegisterAppEnvScenarioOverrides(env)
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

func (steps *PropagatorServiceSteps) requestPropagation(ctx context.Context, serviceName string, propagationRequest *godog.Table) error {
	propSettings, err := GodogTableToSlice[models_service.SatellitePropagationRequest](propagationRequest)
	if err != nil {
		log.Printf("Error parsing propagation request: %v", err)
		return err
	}

	log.Printf("Requesting propagation for service: %s with settings: %+v", serviceName, propSettings)
	return testservice.PropagatorRequest(ctx, steps.state, models_service.ServiceName(serviceName), propSettings)
}

func (steps *PropagatorServiceSteps) subscribeToEvents(consumer string, eventTable *godog.Table) error {
	events, err := GodogTableToSlice[models_service.EventCallbackInfo](eventTable)
	if err != nil {
		log.Printf("Error parsing event subscription table: %v", err)
		return err
	}

	log.Printf("Subscribing as consumer '%s' with events: %+v", consumer, events)
	_, err = testservice.Subscribe(context.Background(), steps.state, consumer, events)
	return err
}

func (steps *PropagatorServiceSteps) verifyEvents(serviceName string, eventTable *godog.Table) error {
	expectedEvents, err := GodogTableToSlice[models_service.ExpectedEvent](eventTable)
	if err != nil {
		log.Printf("Error parsing expected event table: %v", err)
		return err
	}

	log.Printf("Verifying expected events for service: %s", serviceName)
	return testservice.VerifyEvents(steps.state, serviceName, expectedEvents)
}
