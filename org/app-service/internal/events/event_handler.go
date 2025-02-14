package events

import (
	"encoding/json"
	"fmt"

	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

// EventHandler is a non-generic interface that all handlers must implement.
type EventHandler interface {
	Run(event model.EventRoot) error // Parses the event payload and processes it
}

// BaseHandler provides reusable logic for event handlers (uses generics for payloads).
type BaseHandler[T any] struct{}

// NewBaseHandler creates a new BaseHandler.
func NewBaseHandler[T any]() *BaseHandler[T] {
	return &BaseHandler[T]{}
}

// Parse decodes the JSON payload into the expected struct type.
func (h *BaseHandler[T]) Parse(payload string) (*T, error) {
	var result T
	err := json.Unmarshal([]byte(payload), &result)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to parse event payload: %w", err)
	}
	return &result, nil
}

// Run automatically parses the event payload and executes the handler.
func (h *BaseHandler[T]) Run(event model.EventRoot, handler func(event model.EventRoot, payload *T) error) error {
	// Parse payload into expected type
	payload, err := h.Parse(event.Payload)
	if err != nil {
		log.Debugf("❌ Failed to parse payload for event type %s: %v", event.EventType, err)
		return err
	}

	log.Debugf("✅ Successfully parsed event payload for event type: %s", event.EventType)
	return handler(event, payload)
}
