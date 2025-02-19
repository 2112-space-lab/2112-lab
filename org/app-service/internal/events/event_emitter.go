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

// EventEmitter wraps the RabbitMQ client.
type EventEmitter struct {
	rabbitClient *rabbitmq.RabbitMQClient
}

// NewEventEmitter initializes a new EventEmitter using an existing RabbitMQ client.
func NewEventEmitter(ctx context.Context, rabbitClient *rabbitmq.RabbitMQClient) (e *EventEmitter, err error) {
	_, span := tracing.NewSpan(ctx, "EventEmitter.NewEventEmitter")
	defer span.EndWithError(err)

	return &EventEmitter{
		rabbitClient: rabbitClient,
	}, nil
}

// PublishEvent sends an event to RabbitMQ.
func (e *EventEmitter) PublishEvent(ctx context.Context, event model.EventRoot) (err error) {
	ctx, span := tracing.NewSpan(ctx, "PublishEvent")
	defer span.EndWithError(err)

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = e.rabbitClient.PublishMessage(ctx, body, rabbitmq.NewHeader())
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Tracef("[x] Event published: %s", string(body))
	return nil
}
