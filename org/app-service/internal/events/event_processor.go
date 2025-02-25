package events

import (
	"context"
	"sync"

	log "github.com/org/2112-space-lab/org/app-service/pkg/log"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	domainenum "github.com/org/2112-space-lab/org/app-service/internal/domain/domain-enums"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
)

// EventProcessor manages event dispatching and execution with persistence.
type EventProcessor struct {
	eventQueue    chan model.EventRoot
	eventHandlers map[model.EventType][]EventHandler
	eventRepo     repository.EventRepository
	handlerRepo   repository.EventHandlerRepository
	mutex         sync.Mutex
	wg            sync.WaitGroup
}

// NewEventProcessor initializes an EventProcessor with event persistence.
func NewEventProcessor(eventRepo repository.EventRepository, handlerRepo repository.EventHandlerRepository) *EventProcessor {
	return &EventProcessor{
		eventQueue:    make(chan model.EventRoot, DefaultEventQueueSize),
		eventHandlers: make(map[model.EventType][]EventHandler),
		eventRepo:     eventRepo,
		handlerRepo:   handlerRepo,
	}
}

// RegisterHandler registers an event handler for a specific event type.
func (ep *EventProcessor) RegisterHandler(eventType model.EventType, handler EventHandler) {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()

	ep.eventHandlers[eventType] = append(ep.eventHandlers[eventType], handler)
	log.Infof("✅ Registered handler for event: %s", eventType)
}

// EmitEvent stores and queues an event for processing.
func (ep *EventProcessor) BroadcastEvent(ctx context.Context, event model.EventRoot) error {
	ev, err := ConvertToDomainEvent(event)
	if err != nil {
		log.Errorf("❌ Failed to convert event to domain event: %v", err)
		return err
	}
	if err := ep.eventRepo.Save(ctx, ev); err != nil {
		log.Errorf("❌ Failed to store event in database: %v", err)
		return err
	}

	select {
	case ep.eventQueue <- event:
		log.Debugf("📤 Event emitted: %s | UID: %s", event.EventType, event.EventUID)
	default:
		log.Warnf("⚠️ Event queue is full, dropping event: %s", event.EventType)
	}
	return nil
}

// StartProcessing continuously listens for events and executes their handlers.
func (ep *EventProcessor) StartProcessing(ctx context.Context) {
	log.Info("📡 Event Processor started, listening for events...")

	for {
		select {
		case <-ctx.Done():
			log.Info("🛑 Event Processor shutting down...")
			ep.shutdown()
			return
		case event := <-ep.eventQueue:
			ep.processEvent(ctx, event)
		}
	}
}

// processEvent dispatches the event to all registered handlers and logs execution.
func (ep *EventProcessor) processEvent(ctx context.Context, event model.EventRoot) {
	ep.mutex.Lock()
	handlers, exists := ep.eventHandlers[model.EventType(event.EventType)]
	ep.mutex.Unlock()

	if !exists {
		log.Warnf("⚠️ No handlers registered for event type: %s", event.EventType)
		return
	}

	for _, handler := range handlers {
		ep.wg.Add(1)

		go func(h EventHandler) {
			defer ep.wg.Done()

			handler := domain.EventHandler{
				EventID:     event.EventUID,
				HandlerName: h.HandlerName(),
				StartedAt:   xtime.UtcNow(),
				Status:      domainenum.HandlerStates.Started(),
			}

			if err := ep.handlerRepo.Save(ctx, handler); err != nil {
				log.Errorf("❌ Failed to log handler start: %v", err)
			}

			if err := h.Run(ctx, event); err != nil {
				log.Errorf("❌ Error processing event %s: %v", event.EventType, err)

				errorMsg := err.Error()
				handler.Status = domainenum.HandlerStates.Failed()
				handler.Error = fx.NewValueOption(errorMsg)
			} else {
				handler.Status = domainenum.HandlerStates.Completed()
			}

			handler.CompletedAt = fx.NewValueOption(xtime.UtcNow())
			if err := ep.handlerRepo.Save(ctx, handler); err != nil {
				log.Errorf("❌ Failed to update handler execution log: %v", err)
			}
		}(handler)
	}
}

// shutdown ensures all ongoing event processing completes before shutting down.
func (ep *EventProcessor) shutdown() {
	log.Info("⚠️ Draining event queue before shutdown...")
	close(ep.eventQueue)
	ep.wg.Wait()
	log.Info("✅ Event Processor shutdown complete.")
}
