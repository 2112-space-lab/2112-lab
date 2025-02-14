package events

import (
	"encoding/json"
	"fmt"

	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
)

// EventParser handles JSON event parsing
type EventParser struct{}

// NewEventParser initializes an EventParser
func NewEventParser() *EventParser {
	return &EventParser{}
}

// ParseEvent parses the incoming JSON message into an EventRoot struct
func (p *EventParser) ParseEvent(jsonData []byte) (*model.EventRoot, error) {
	var event model.EventRoot
	err := json.Unmarshal(jsonData, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event root: %w", err)
	}
	return &event, nil
}

// Generic type for parsing event payloads
type PayloadParser[T any] struct{}

// NewPayloadParser initializes a PayloadParser for a specific type
func NewPayloadParser[T any]() *PayloadParser[T] {
	return &PayloadParser[T]{}
}

// Parse decodes the JSON payload into the specified struct type
func (p *PayloadParser[T]) Parse(payload string) (*T, error) {
	var result T
	err := json.Unmarshal([]byte(payload), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event payload: %w", err)
	}
	return &result, nil
}
