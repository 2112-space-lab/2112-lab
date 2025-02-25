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

// NewRehydrateEvent creates an EventRoot for a RehydrateGameContext event
func NewRehydrateEvent(name string, comment *string) (*model.EventRoot, error) {
	payload := model.RehydrateGameContext{
		Name: name,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize rehydrate event payload: %w", err)
	}

	return &model.EventRoot{
		EventTimeUtc: generateEventTimestamp(),
		EventUID:     generateEventUID(),
		EventType:    model.EventTypeRehydrateGameContext.String(),
		Comment:      comment,
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
