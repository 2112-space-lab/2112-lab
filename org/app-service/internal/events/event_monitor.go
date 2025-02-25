package events

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/rabbitmq"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	domainenum "github.com/org/2112-space-lab/org/app-service/internal/domain/domain-enums"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DefaultEventQueueSize = 100
)

// EventMonitor handles event subscription, processing, and persistence.
type EventMonitor struct {
	rabbitClient  *rabbitmq.RabbitMQClient
	eventQueue    chan []byte
	eventHandlers map[model.EventType][]EventHandler
	mutex         sync.Mutex
	eventRepo     repository.EventRepository
	handlerRepo   repository.EventHandlerRepository
}

// NewEventMonitor initializes an EventMonitor with persistence.
func NewEventMonitor(ctx context.Context, rabbitClient *rabbitmq.RabbitMQClient, eventRepo repository.EventRepository, handlerRepo repository.EventHandlerRepository) (e *EventMonitor, err error) {
	_, span := tracing.NewSpan(ctx, "EventMonitor.NewEventMonitor")
	defer span.EndWithError(err)

	return &EventMonitor{
		rabbitClient:  rabbitClient,
		eventQueue:    make(chan []byte, DefaultEventQueueSize),
		eventHandlers: make(map[model.EventType][]EventHandler),
		eventRepo:     eventRepo,
		handlerRepo:   handlerRepo,
	}, nil
}

// RegisterHandler associates an event type with a handler.
func (m *EventMonitor) RegisterHandler(ctx context.Context, eventType model.EventType, handler EventHandler) (err error) {
	_, span := tracing.NewSpan(ctx, "EventMonitor.RegisterHandler")
	defer span.EndWithError(err)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.eventHandlers == nil {
		m.eventHandlers = make(map[model.EventType][]EventHandler)
	}

	m.eventHandlers[eventType] = append(m.eventHandlers[eventType], handler)
	log.Infof("✅ Registered handler for event: %s", eventType)
	return nil
}

// StartMonitoring continuously listens for RabbitMQ messages and processes events.
func (m *EventMonitor) StartMonitoring(ctx context.Context, header *rabbitmq.Header) (err error) {
	ctx, span := tracing.NewSpan(ctx, "StartMonitoring")
	defer span.EndWithError(err)

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

// processEvents processes events asynchronously from the queue and stores them.
func (m *EventMonitor) processEvents(ctx context.Context) {
	for msgBody := range m.eventQueue {
		var eventRoot model.EventRoot
		if err := json.Unmarshal(msgBody, &eventRoot); err != nil {
			log.Errorf("❌ Failed to parse event: %v", err)
			continue
		}

		log.Infof("🔹 Received event: %s | UID: %s", eventRoot.EventType, eventRoot.EventUID)

		event, err := ConvertToDomainEvent(eventRoot)
		if err != nil {
			log.Errorf("❌ Failed to convert event to domain: %v", err)
			continue
		}

		if err := m.eventRepo.Save(ctx, event); err != nil {
			log.Errorf("❌ Failed to store event in database: %v", err)
			continue
		}

		m.mutex.Lock()
		handlers, exists := m.eventHandlers[model.EventType(eventRoot.EventType)]
		m.mutex.Unlock()

		if !exists {
			log.Warnf("⚠️ No handler registered for event type: %s", eventRoot.EventType)
			continue
		}

		for _, handler := range handlers {
			go m.executeHandler(ctx, handler, eventRoot)
		}
	}
}

// executeHandler processes an event with a handler and logs execution.
func (m *EventMonitor) executeHandler(ctx context.Context, handler EventHandler, event model.EventRoot) {
	handlerLog := domain.EventHandler{
		EventID:     event.EventUID,
		HandlerName: handler.HandlerName(),
		StartedAt:   xtime.UtcNow(),
		Status:      domainenum.HandlerStates.Started(),
	}

	if err := m.handlerRepo.Save(ctx, handlerLog); err != nil {
		log.Errorf("❌ Failed to log handler start: %v", err)
	}

	if err := handler.Run(ctx, event); err != nil {
		log.Errorf("❌ Error processing event %s: %v", event.EventType, err)

		errorMsg := err.Error()
		handlerLog.Status = domainenum.HandlerStates.Failed()
		handlerLog.Error = fx.NewValueOption(errorMsg)
	} else {
		handlerLog.Status = domainenum.HandlerStates.Completed()
	}

	handlerLog.CompletedAt = fx.NewValueOption(xtime.UtcNow())

	if err := m.handlerRepo.Save(ctx, handlerLog); err != nil {
		log.Errorf("❌ Failed to update handler execution log: %v", err)
	}
}
