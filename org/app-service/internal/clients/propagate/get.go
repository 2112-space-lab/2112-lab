package propagator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Elbujito/2112/src/app-service/internal/config"
)

// PropagatorClient definition
type PropagatorClient struct {
	env *config.SEnv
}

// SatellitePropagationRequest represents the payload for the propagation request.
type SatellitePropagationRequest struct {
	TLELine1        string `json:"tle_line1"`
	TLELine2        string `json:"tle_line2"`
	StartTime       string `json:"start_time"`
	DurationMinutes int    `json:"duration_minutes"`
	IntervalSeconds int    `json:"interval_seconds"`
	NoradID         string `json:"norad_id"`
}

// SatellitePosition represents the propagated position of a satellite.
type SatellitePosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Time      string  `json:"timestamp"`
}

// PropagationResponse definition
type PropagationResponse struct {
	Positions []*SatellitePosition `json:"positions"`
}

// NewPropagatorClient creates a new PropagatorClient with the given base URL.
func NewPropagatorClient(env *config.SEnv) *PropagatorClient {
	return &PropagatorClient{
		env: env,
	}
}

// FetchPropagation fetches propagated positions for a given TLE and parameters without waiting for the response.
func (client *PropagatorClient) FetchPropagation(
	ctx context.Context,
	tle1, tle2, startTime string,
	durationMinutes, intervalSeconds int, noradID string,
) (<-chan []*SatellitePosition, <-chan error) {
	// Create channels for results and errors
	resultChan := make(chan []*SatellitePosition, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		// Validate input
		if tle1 == "" || tle2 == "" || startTime == "" {
			errorChan <- fmt.Errorf("TLE lines and start time are required")
			return
		}
		if durationMinutes <= 0 || intervalSeconds <= 0 {
			errorChan <- fmt.Errorf("duration and interval must be greater than zero")
			return
		}

		// Prepare request payload
		requestPayload := SatellitePropagationRequest{
			TLELine1:        tle1,
			TLELine2:        tle2,
			StartTime:       startTime,
			DurationMinutes: durationMinutes,
			IntervalSeconds: intervalSeconds,
			NoradID:         noradID,
		}
		payloadBytes, err := json.Marshal(requestPayload)
		if err != nil {
			errorChan <- fmt.Errorf("failed to marshal request payload: %v", err)
			return
		}

		config := client.env.EnvVars.Propagator

		// Send the HTTP request
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.BaseUrl, bytes.NewBuffer(payloadBytes))
		if err != nil {
			errorChan <- fmt.Errorf("failed to create HTTP request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("failed to send HTTP request: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("failed to fetch propagation data: HTTP status %d, response: %s", resp.StatusCode, string(body))
			return
		}

		// Parse the response
		var response PropagationResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("failed to parse response body: %v, response body: %s", err, string(body))
			return
		}

		resultChan <- response.Positions
	}()

	return resultChan, errorChan
}
