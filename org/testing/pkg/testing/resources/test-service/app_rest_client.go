package testservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	models_service "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtestlog "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-log"
)

type AppRestClientRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type AppRestClient struct {
	hostName   string
	httpClient AppRestClientRequestDoer
}

func GetAppRestClient(ctx context.Context, scenarioState AppClientScenarioState, serviceName models_service.ServiceName) (*AppRestClient, error) {
	logger := scenarioState.GetLogger()
	cont, err := scenarioState.GetAppServiceContainer(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	httpBoundPort, err := cont.GetBoundPort(ctx, testservicecontainer.AppServiceHttpPort)
	if err != nil {
		return nil, err
	}
	hostname, err := cont.GetHostName(ctx)
	if err != nil {
		return nil, err
	}

	restClient := &AppRestClient{
		hostName:   fmt.Sprintf("http://%s:%s", hostname, httpBoundPort.Port()),
		httpClient: xtestlog.NewRequestClientLogger(&http.Client{}, logger),
	}
	return restClient, err
}

// CreateContext sends a request to create a new GameContext.
func (c *AppRestClient) CreateContext(ctx context.Context, gameContext models_service.GameContext) (*models_service.GameContext, error) {
	url := fmt.Sprintf("%s/contexts/", c.hostName)

	reqBody, err := json.Marshal(gameContext)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create context, status code: %d", resp.StatusCode)
	}

	var createdContext models_service.GameContext
	if err := json.NewDecoder(resp.Body).Decode(&createdContext); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &createdContext, nil
}

// UpdateContext sends a request to update an existing GameContext.
func (c *AppRestClient) UpdateContext(ctx context.Context, name string, gameContext models_service.GameContext) (*models_service.GameContext, error) {
	url := fmt.Sprintf("%s/contexts/%s", c.hostName, name)

	reqBody, err := json.Marshal(gameContext)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update context, status code: %d", resp.StatusCode)
	}

	var updatedContext models_service.GameContext
	if err := json.NewDecoder(resp.Body).Decode(&updatedContext); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &updatedContext, nil
}

// AssignSatellitesToContext assigns satellites to an existing GameContext.
func (c *AppRestClient) AssignSatellitesToContext(ctx context.Context, contextName string, satelliteNames []string) error {
	url := fmt.Sprintf("%s/contexts/%s/assign/satellites", c.hostName, contextName)

	payload := map[string]interface{}{
		"name":           contextName,
		"satelliteNames": satelliteNames,
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to assign satellites, status code: %d", resp.StatusCode)
	}

	return nil
}

// ActivateContext activates a GameContext by its unique name.
func (c *AppRestClient) ActivateContext(ctx context.Context, contextName string) error {
	url := fmt.Sprintf("%s/contexts/%s/activate", c.hostName, contextName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to activate context, status code: %d", resp.StatusCode)
	}

	return nil
}

// RehydrateContext triggers a rehydration process for a GameContext.
func (c *AppRestClient) RehydrateContext(ctx context.Context, contextName string) error {
	url := fmt.Sprintf("%s/contexts/%s/rehydrate", c.hostName, contextName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to rehydrate context, status code: %d", resp.StatusCode)
	}

	return nil
}

// GetContextByName retrieves a GameContext by its unique name.
func (c *AppRestClient) GetContextByName(ctx context.Context, name string) (*models_service.GameContext, error) {
	url := fmt.Sprintf("%s/contexts/%s", c.hostName, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("context not found, status code: %d", resp.StatusCode)
	}

	var gameContext models_service.GameContext
	if err := json.NewDecoder(resp.Body).Decode(&gameContext); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &gameContext, nil
}

// DeleteContextByName deletes a GameContext by its unique name.
func (c *AppRestClient) DeleteContextByName(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/contexts/%s", c.hostName, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete context, status code: %d", resp.StatusCode)
	}

	return nil
}
