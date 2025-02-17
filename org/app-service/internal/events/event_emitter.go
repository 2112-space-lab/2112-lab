package events

import (
	"encoding/json"
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/rabbitmq"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

// EventEmitter wraps the RabbitMQ client.
type EventEmitter struct {
	rabbitClient *rabbitmq.RabbitMQClient
}

// NewEventEmitter initializes a new EventEmitter using an existing RabbitMQ client.
func NewEventEmitter(rabbitClient *rabbitmq.RabbitMQClient) *EventEmitter {
	return &EventEmitter{
		rabbitClient: rabbitClient,
	}
}

// PublishEvent sends an event to RabbitMQ.
func (e *EventEmitter) PublishEvent(event model.EventRoot) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = e.rabbitClient.PublishMessage(body, rabbitmq.NewHeader())
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Tracef("[x] Event published: %s", string(body))
	return nil
}
