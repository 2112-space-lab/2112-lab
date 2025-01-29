package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Elbujito/2112/src/graphql-api/go/graph/model"
)

type mutationResolver struct {
	*CustomResolver
}

// RequestSatelliteVisibilities handles the mutation to request satellite visibility calculations.
func (r *mutationResolver) RequestSatelliteVisibilities(
	ctx context.Context,
	uid string,
	userLocation model.UserLocationInput,
	startTime string,
	endTime string,
) (bool, error) {
	// Redis channel for visibility requests
	const visibilityRequestChannel = "visibility_requests"

	log.Printf(
		"RequestSatelliteVisibilities called: uid=%s, latitude=%.6f, longitude=%.6f, radius=%.2f, horizon=%.2f, startTime=%s, endTime=%s",
		uid, userLocation.Latitude, userLocation.Longitude, userLocation.Radius, userLocation.Horizon, startTime, endTime,
	)

	// Validate input parameters
	if err := validateRequestParameters(uid, userLocation, startTime, endTime); err != nil {
		log.Printf("Invalid request parameters: %v", err)
		return false, fmt.Errorf("invalid request parameters: %w", err)
	}

	// Prepare the request payload
	requestPayload := map[string]interface{}{
		"uid":       uid,
		"latitude":  userLocation.Latitude,
		"longitude": userLocation.Longitude,
		"radius":    userLocation.Radius,
		"horizon":   userLocation.Horizon,
		"startTime": startTime,
		"endTime":   endTime,
	}

	// Serialize the payload into JSON
	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		log.Printf("Error marshalling request payload: %v", err)
		return false, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Publish the request to Redis
	if err := r.rdb.Publish(ctx, visibilityRequestChannel, requestJSON).Err(); err != nil {
		log.Printf("Error publishing visibility request to Redis channel '%s': %v", visibilityRequestChannel, err)
		return false, fmt.Errorf("failed to publish request to channel '%s': %w", visibilityRequestChannel, err)
	}

	log.Printf("Successfully published visibility request to Redis channel '%s': %s", visibilityRequestChannel, string(requestJSON))
	return true, nil
}

// validateRequestParameters validates the input parameters for the visibility request.
func validateRequestParameters(
	uid string,
	userLocation model.UserLocationInput,
	startTime string,
	endTime string,
) error {
	if uid == "" {
		return fmt.Errorf("uid cannot be empty")
	}

	if userLocation.Latitude < -90 || userLocation.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90 degrees")
	}

	if userLocation.Longitude < -180 || userLocation.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180 degrees")
	}

	if userLocation.Radius <= 0 {
		return fmt.Errorf("radius must be greater than 0")
	}

	if userLocation.Horizon < 0 || userLocation.Horizon > 90 {
		return fmt.Errorf("horizon must be between 0 and 90 degrees")
	}

	if startTime == "" || endTime == "" {
		return fmt.Errorf("startTime and endTime cannot be empty")
	}

	return nil
}
