package testservice

import (
	"context"

	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
)

func PropagatorRequest(ctx context.Context, scenarioState PropagatorClientScenarioState, serviceName models_service.ServiceName, requests []models_service.SatellitePropagationRequest) error {
	client, err := GetPropagatorRestClient(ctx, scenarioState, serviceName)
	if err != nil {
		return err
	}

	for _, r := range requests {
		err = client.RequestSatellitePropagation(ctx, r)
		if err != nil {
			return err
		}
	}

	return nil
}
