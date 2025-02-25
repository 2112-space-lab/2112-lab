package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/rabbitmq"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// EventEmitter wraps the RabbitMQ client and integrates with the EventProcessor.
type EventEmitter struct {
	rabbitClient *rabbitmq.RabbitMQClient
	processor    *EventProcessor
}

// NewEventEmitter initializes a new EventEmitter using RabbitMQ and an EventProcessor.
func NewEventEmitter(ctx context.Context, rabbitClient *rabbitmq.RabbitMQClient, processor *EventProcessor) (*EventEmitter, error) {
	_, span := tracing.NewSpan(ctx, "EventEmitter.NewEventEmitter")
	defer span.End()

	return &EventEmitter{
		rabbitClient: rabbitClient,
		processor:    processor,
	}, nil
}

// PublishEvent sends an event to the local EventProcessor.
func (e *EventEmitter) PublishEventToProcessorOnly(ctx context.Context, event model.EventRoot) error {
	ctx, span := tracing.NewSpan(ctx, "PublishEvent")
	defer span.End()

	e.processor.EmitEvent(ctx, event)
	log.Tracef("[x] Event processed locally: %s", event.EventType)

	return nil
}

// PublishEvent sends an event to RabbitMQ and also to the local EventProcessor.
func (e *EventEmitter) PublishEvent(ctx context.Context, event model.EventRoot) error {
	ctx, span := tracing.NewSpan(ctx, "PublishEvent")
	defer span.End()

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = e.rabbitClient.PublishMessage(ctx, body, rabbitmq.NewHeader())
	if err != nil {
		log.Errorf("❌ Failed to publish event to RabbitMQ: %v", err)
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
