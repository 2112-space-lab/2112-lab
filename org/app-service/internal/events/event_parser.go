package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
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

// ConvertToDomainEvent converts a model.EventRoot (GraphQL model) to a domain.Event (Domain model).
func ConvertToDomainEvent(eventRoot model.EventRoot) (domain.Event, error) {
	eventTime, err := time.Parse(time.RFC3339, eventRoot.EventTimeUtc)
	if err != nil {
		return domain.Event{}, fmt.Errorf("failed to parse event timestamp: %w", err)
	}

	domainEvent := domain.Event{
		ModelBase:   domain.NewModelBaseDefault(),
		EventType:   domain.EventType(eventRoot.EventType),
		EventUID:    eventRoot.EventUID,
		Payload:     fx.NewValueOption(eventRoot.Payload),
		PublishedAt: xtime.NewUtcTimeIgnoreZone(eventTime),
		Comment:     fx.AsOption(eventRoot.Comment),
	}

	return domainEvent, nil
}
