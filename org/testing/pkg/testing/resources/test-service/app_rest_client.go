package testservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service/models"
	xtime "github.com/org/2112-space-lab/org/testing/pkg/x-time"
)

// AppRestClient manages API calls manually without app_client.
type AppRestClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAppRestClient initializes a new AppRestClient.
func NewAppRestClient(ctx context.Context, scenarioState AppClientScenarioState, serviceName models.ServiceName) (*AppRestClient, error) {
	baseURL, err := getServiceBaseURL(ctx, scenarioState, serviceName)
	if err != nil {
		return nil, err
	}
	return &AppRestClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}, nil
}

// getServiceBaseURL determines the service URL dynamically.
func getServiceBaseURL(ctx context.Context, scenarioState AppClientScenarioState, serviceName models.ServiceName) (string, error) {
	cont, err := scenarioState.GetAppServiceContainer(ctx, serviceName)
	if err != nil {
		return "", err
	}

	httpBoundPort, err := cont.GetBoundPort(ctx, testservicecontainer.AppServiceHttpPort)
	if err != nil {
		return "", err
	}
	hostname, err := cont.GetHostName(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%s", hostname, httpBoundPort.Port()), nil
}

// sendRequest is a helper function to make HTTP requests.
func (c *AppRestClient) sendRequest(ctx context.Context, method, endpoint string, body interface{}) ([]byte, int, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// CreateContext sends a request to create a new models.GameContext.
func (c *AppRestClient) CreateContext(ctx context.Context, gameContext models.GameContext) (*models.GameContext, error) {
	respBody, statusCode, err := c.sendRequest(ctx, "POST", "/contexts/", gameContext)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create context, status code: %d", statusCode)
	}

	var createdContext models.GameContext
	if err := json.Unmarshal(respBody, &createdContext); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &createdContext, nil
}

// UpdateContext sends a request to update an existing models.GameContext.
func (c *AppRestClient) UpdateContext(ctx context.Context, name string, gameContext models.GameContext) (*models.GameContext, error) {
	respBody, statusCode, err := c.sendRequest(ctx, "PUT", "/contexts/"+name, gameContext)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update context, status code: %d", statusCode)
	}

	var updatedContext models.GameContext
	if err := json.Unmarshal(respBody, &updatedContext); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &updatedContext, nil
}

// AssignSatellitesToContext assigns satellites to an existing GameContext.
func (c *AppRestClient) AssignSatellitesToContext(ctx context.Context, contextName string, satelliteNames []string) error {
	payload := map[string]interface{}{
		"name":           contextName,
		"satelliteNames": satelliteNames,
	}

	respBody, statusCode, err := c.sendRequest(ctx, "POST", fmt.Sprintf("/contexts/%s/assign/satellites", contextName), payload)
	if err != nil {
		return fmt.Errorf("failed to assign satellites to context: %w", err)
	}

	if statusCode != http.StatusNoContent {
		return fmt.Errorf("failed to assign satellites, status code: %d, response: %s", statusCode, string(respBody))
	}

	return nil
}

// ActivateContext activates a GameContext by its unique name.
func (c *AppRestClient) ActivateContext(ctx context.Context, contextName string) error {
	respBody, statusCode, err := c.sendRequest(ctx, "PUT", fmt.Sprintf("/contexts/%s/activate", contextName), nil)
	if err != nil {
		return fmt.Errorf("failed to activate context: %w", err)
	}
	if statusCode != http.StatusNoContent {
		return fmt.Errorf("failed to activate context, status code: %d, response: %s", statusCode, string(respBody))
	}
	return nil
}

// GetContextByName retrieves a models.GameContext by its unique name.
func (c *AppRestClient) GetContextByName(ctx context.Context, name string) (*models.GameContext, error) {
	respBody, statusCode, err := c.sendRequest(ctx, "GET", "/contexts/"+name, nil)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("context not found, status code: %d", statusCode)
	}

	var gameContext models.GameContext
	if err := json.Unmarshal(respBody, &gameContext); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &gameContext, nil
}

// DeleteContextByName deletes a models.GameContext by its unique name.
func (c *AppRestClient) DeleteContextByName(ctx context.Context, name string) error {
	_, statusCode, err := c.sendRequest(ctx, "DELETE", "/contexts/"+name, nil)
	if err != nil {
		return err
	}
	if statusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete context, status code: %d", statusCode)
	}

	return nil
}

// MapToOptionalTime converts a string to an optional time.Time pointer.
func MapToOptionalTime(value string) *time.Time {
	if value != "" {
		t, err := xtime.FromString(xtime.DateTimeFormat(value))
		if err != nil {
			return nil
		}
		ti := t.Inner()
		return &ti
	}
	return nil
}

// MapToOptionalString converts a string to an optional string pointer.
func MapToOptionalString(value string) *string {
	if value != "" {
		str := fmt.Sprintf("%v", value)
		return &str
	}
	return nil
}
