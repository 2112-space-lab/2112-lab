package event_builder

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
)

// Generate a new unique event UID
func generateEventUID() string {
	return uuid.New().String()
}

// Generate the current UTC timestamp in ISO 8601 format
func generateEventTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// NewRehydrateGameContextRequestedEvent creates an EventRoot for a RehydrateGameContextRequested event
func NewRehydrateGameContextRequestedEvent(name string) (*model.EventRoot, error) {
	payload := model.RehydrateGameContextRequested{
		Name:        name,
		TriggeredAt: generateEventTimestamp(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize rehydrate requested event payload: %w", err)
	}

	return &model.EventRoot{
		EventTimeUtc: generateEventTimestamp(),
		EventUID:     generateEventUID(),
		EventType:    model.EventTypeRehydrateGameContextRequested.String(),
		Payload:      string(payloadBytes),
	}, nil
}

// NewRehydrateGameContextSuccessEvent creates an EventRoot for a RehydrateGameContextSuccess event
func NewRehydrateGameContextSuccessEvent(name string, nbSatellites int32) (*model.EventRoot, error) {
	payload := model.RehydrateGameContextSuccess{
		Name:         name,
		NbSatellites: nbSatellites,
		CompletedAt:  generateEventTimestamp(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize rehydrate success event payload: %w", err)
	}

	return &model.EventRoot{
		EventTimeUtc: generateEventTimestamp(),
		EventUID:     generateEventUID(),
		EventType:    model.EventTypeRehydrateGameContextSuccess.String(),
		Payload:      string(payloadBytes),
	}, nil
}

// NewRehydrateGameContextFailedEvent creates an EventRoot for a RehydrateGameContextFailed event
func NewRehydrateGameContextFailedEvent(name, reason string, failureCount int32) (*model.EventRoot, error) {
	payload := model.RehydrateGameContextFailed{
		Name:         name,
		Reason:       reason,
		FailureCount: failureCount,
		FailedAt:     generateEventTimestamp(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize rehydrate failed event payload: %w", err)
	}

	return &model.EventRoot{
		EventTimeUtc: generateEventTimestamp(),
		EventUID:     generateEventUID(),
		EventType:    model.EventTypeRehydrateGameContextFailed.String(),
		Payload:      string(payloadBytes),
	}, nil
}

// NewSatelliteTlePropagationRequestedEvent creates an EventRoot for a SatelliteTLEPropagationRequested event
func NewSatelliteTlePropagationRequestedEvent(request model.SatelliteTlePropagated, comment *string) (*model.EventRoot, error) {
	payloadBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize propagation event payload: %w", err)
	}

	return &model.EventRoot{
		EventTimeUtc: generateEventTimestamp(),
		EventUID:     generateEventUID(),
		EventType:    model.EventTypeSatelliteTlePropagationRequested.String(),
		Comment:      comment,
		Payload:      string(payloadBytes),
	}, nil
}
