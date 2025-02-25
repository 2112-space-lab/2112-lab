package steps

import (
	"context"
	"fmt"
	"log"

	"github.com/cucumber/godog"
	testservice "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service"
	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	xtime "github.com/org/2112-space-lab/org/testing/pkg/x-time"
)

type AppServiceSteps struct {
	state                       appServiceStepsState
	appServiceContainerResource appServiceContainerResource
	restClient                  *testservice.AppRestClient
}

type appServiceContainerResource interface {
	SpawnServiceAppService(ctx context.Context, scenarioState testservicecontainer.GatewwayContainerScenarioState, serviceName models_service.ServiceName, serviceEnvOverrides models_cont.EnvVarKeyValueMap) error
}

type appServiceStepsState interface {
	testservicecontainer.GatewwayContainerScenarioState
	testservice.AppClientScenarioState
}

func RegisterAppServiceSteps(ctx *godog.ScenarioContext, state appServiceStepsState, containerRsc appServiceContainerResource) error {
	restClient, err := testservice.NewAppRestClient(context.Background(), state, models_service.ServiceName("localhost"))
	if err != nil {
		return err
	}
	s := &AppServiceSteps{
		state:                       state,
		appServiceContainerResource: containerRsc,
		restClient:                  restClient,
	}
	ctx.Step(`^a App service is created for service "([^"]*)"$`, s.appServiceCreate)
	ctx.Step(`^a App service is created for service "([^"]*)" with env overrides:$`, s.appServiceCreateWithEnv)
	ctx.Step(`^I register App service default scenario environment variable overrides:$`, s.registerCommonEnvVars)
	ctx.Step(`^I subscribes as consumer "([^"]*)" for "([^"]*)" with registered callbacks:$`, s.subscribeToEvents)
	ctx.Step(`^Events are expected for service "([^"]*)":$`, s.verifyEvents)
	ctx.Step(`^I create a game context "([^"]*)" with the following satellites:$`, s.createGameContextWithSatellites)
	ctx.Step(`^I activate the game context "([^"]*)"$`, s.activateGameContext)

	return nil
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

func (steps *AppServiceSteps) subscribeToEvents(consumer string, queueName string, eventTable *godog.Table) error {
	events, err := GodogTableToSlice[models_service.EventCallbackInfo](eventTable)
	if err != nil {
		log.Printf("Error parsing event subscription table: %v", err)
		return err
	}

	log.Printf("Subscribing as consumer '%s' with events: %+v", consumer, events)
	_, err = testservice.Subscribe(context.Background(), steps.state, models_service.ServiceName(consumer), queueName, events)
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

// createGameContextWithSatellites creates a new GameContext and assigns satellites.
func (steps *AppServiceSteps) createGameContextWithSatellites(ctx context.Context, contextName string, satelliteTable *godog.Table) error {
	satellites, err := GodogTableToSlice[models_service.SatelliteDefinition](satelliteTable)
	if err != nil {
		log.Printf("Error parsing satellite table: %v", err)
		return err
	}

	// Extract satellite names
	satelliteNames := make([]string, len(satellites))
	for i, satellite := range satellites {
		satelliteNames[i] = satellite.Name
	}

	// Create the GameContext
	gameContext := models_service.GameContext{
		Name:        contextName,
		Description: fmt.Sprintf("Context created at: %s", xtime.UtcNow().Inner()),
		IsActive:    false,
	}

	createdContext, err := steps.restClient.CreateContext(ctx, gameContext)
	if err != nil {
		return fmt.Errorf("failed to create game context: %w", err)
	}

	log.Printf("Created GameContext: %+v", createdContext)

	// Assign satellites
	err = steps.restClient.AssignSatellitesToContext(ctx, contextName, satelliteNames)
	if err != nil {
		return fmt.Errorf("failed to assign satellites: %w", err)
	}

	log.Printf("Satellites successfully assigned to context %s", contextName)
	return nil
}

// activateGameContext activates a GameContext.
func (steps *AppServiceSteps) activateGameContext(ctx context.Context, contextName string) error {
	err := steps.restClient.ActivateContext(ctx, contextName)
	if err != nil {
		return fmt.Errorf("failed to activate game context: %w", err)
	}

	log.Printf("GameContext '%s' activated successfully", contextName)
	return nil
}
