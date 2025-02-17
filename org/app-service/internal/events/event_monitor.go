package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/rabbitmq"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultEventQueueSize = 100
)

// EventMonitor handles event subscription and processing.
type EventMonitor struct {
	rabbitClient  *rabbitmq.RabbitMQClient
	eventQueue    chan []byte
	eventHandlers map[model.EventType][]EventHandler
	mutex         sync.Mutex
}

// NewEventMonitor initializes an EventMonitor.
func NewEventMonitor(rabbitClient *rabbitmq.RabbitMQClient) *EventMonitor {
	return &EventMonitor{
		rabbitClient:  rabbitClient,
		eventQueue:    make(chan []byte, DefaultEventQueueSize),
		eventHandlers: make(map[model.EventType][]EventHandler),
	}
}

// RegisterHandler associates an event type with a handler.
func (m *EventMonitor) RegisterHandler(ctx context.Context, eventType model.EventType, handler EventHandler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.eventHandlers[eventType] = append(m.eventHandlers[eventType], handler)
	log.Infof("✅ Registered handler for event: %s", eventType)
}

// StartMonitoring continuously listens for RabbitMQ messages and dispatches events with automatic reconnection.
func (m *EventMonitor) StartMonitoring(ctx context.Context, header *rabbitmq.Header) error {
	log.Info("📡 Event Monitor started. Waiting for messages...")

	go m.processEvents(ctx)

	retryPolicy := createBackoff()
	for {
		select {
		case <-ctx.Done():
			log.Info("🛑 Event Monitor shutting down gracefully...")
			m.rabbitClient.Close()
			return nil
		default:
		}

		var msgs <-chan amqp.Delivery
		err := backoff.Retry(func() error {
			var err error
			msgs, err = m.rabbitClient.ConsumeMessages(header)
			if err != nil {
				log.Warnf("❌ Failed to consume messages: %v. Retrying...", err)
				return err
			}
			return nil
		}, retryPolicy)

		if err != nil {
			log.Errorf("🚨 Permanent failure consuming messages: %v", err)
			return err
		}

		retryPolicy.Reset()
		for {
			select {
			case <-ctx.Done():
				log.Info("🛑 Event Monitor stopping message consumption...")
				m.rabbitClient.Close()
				return nil
			case msg, ok := <-msgs:
				if !ok {
					log.Warn("⚠️ Message channel closed unexpectedly. Reconnecting to RabbitMQ...")
					break
				}
				m.eventQueue <- msg.Body
			}
		}
	}
}

func createBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.MaxInterval = 120 * time.Second
	b.MaxElapsedTime = 20 * time.Minute
	b.RandomizationFactor = .75

	b.Reset()
	return b
}

// processEvents processes events asynchronously from the queue.
func (m *EventMonitor) processEvents(ctx context.Context) {
	for msgBody := range m.eventQueue {
		var event model.EventRoot
		if err := json.Unmarshal(msgBody, &event); err != nil {
			log.Errorf("❌ Failed to parse event: %v", err)
			continue
		}

		log.Infof("🔹 Received event: %s | UID: %s", event.EventType, event.EventUID)

		m.mutex.Lock()
		handlers, exists := m.eventHandlers[model.EventType(event.EventType)]
		m.mutex.Unlock()

		if !exists {
			log.Warnf("⚠️ No handler registered for event type: %s", event.EventType)
			continue
		}

		for _, handler := range handlers {
			go handler.Run(ctx, event)
		}
	}
}
