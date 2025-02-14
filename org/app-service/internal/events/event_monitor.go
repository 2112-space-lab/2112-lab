package events

import (
	"context"
	"encoding/json"
	"math"
	"sync"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/rabbitmq"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
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
		eventQueue:    make(chan []byte, 100),
		eventHandlers: make(map[model.EventType][]EventHandler),
	}
}

// RegisterHandler associates an event type with a handler.
func (m *EventMonitor) RegisterHandler(eventType model.EventType, handler EventHandler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.eventHandlers[eventType] = append(m.eventHandlers[eventType], handler)
	log.Infof("✅ Registered handler for event: %s", eventType)
}

// StartMonitoring continuously listens for RabbitMQ messages and dispatches events with automatic reconnection.
func (m *EventMonitor) StartMonitoring(ctx context.Context) error {
	log.Info("📡 Event Monitor started. Waiting for messages...")

	go m.processEvents()

	backoffFactor := 1

	for {
		select {
		case <-ctx.Done():
			log.Info("🛑 Event Monitor shutting down gracefully...")
			m.rabbitClient.Close()
			return nil
		default:
		}

		msgs, err := m.rabbitClient.ConsumeMessages()
		if err != nil {
			log.Warnf("❌ Failed to consume messages: %v. Retrying in %d seconds...", err, backoffFactor)
			time.Sleep(time.Duration(backoffFactor) * time.Second)

			backoffFactor = int(math.Min(float64(backoffFactor*2), 60))
			continue
		}

		backoffFactor = 1

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

// processEvents processes events asynchronously from the queue.
func (m *EventMonitor) processEvents() {
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
			go handler.Run(event)
		}
	}
}
