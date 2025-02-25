package events

import (
	"context"
	"log"
	"sync"

	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
)

// EventProcessor manages event dispatching and execution.
type EventProcessor struct {
	eventQueue    chan model.EventRoot
	eventHandlers map[model.EventType][]EventHandler
	mutex         sync.Mutex
	wg            sync.WaitGroup
}

// NewEventProcessor initializes an EventProcessor with a predefined queue size.
func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		eventQueue:    make(chan model.EventRoot, 100),
		eventHandlers: make(map[model.EventType][]EventHandler),
	}
}

// RegisterHandler registers an event handler for a specific event type.
func (ep *EventProcessor) RegisterHandler(eventType model.EventType, handler EventHandler) {
	ep.mutex.Lock()
	defer ep.mutex.Unlock()

	ep.eventHandlers[eventType] = append(ep.eventHandlers[eventType], handler)
	log.Printf("✅ Registered handler for event: %s", eventType)
}

// EmitEvent adds an event to the queue for processing.
func (ep *EventProcessor) EmitEvent(ctx context.Context, event model.EventRoot) {
	select {
	case ep.eventQueue <- event:
		log.Printf("📤 Event emitted: %s | UID: %s", event.EventType, event.EventUID)
	default:
		log.Printf("⚠️ Event queue is full, dropping event: %s", event.EventType)
	}
}

// StartProcessing continuously listens for events and executes their handlers.
func (ep *EventProcessor) StartProcessing(ctx context.Context) {
	log.Println("📡 Event Processor started, listening for events...")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Event Processor shutting down...")
			ep.shutdown()
			return
		case event := <-ep.eventQueue:
			ep.processEvent(ctx, event)
		}
	}
}

// processEvent dispatches the event to all registered handlers.
func (ep *EventProcessor) processEvent(ctx context.Context, event model.EventRoot) {
	ep.mutex.Lock()
	handlers, exists := ep.eventHandlers[model.EventType(event.EventType)]
	ep.mutex.Unlock()

	if !exists {
		log.Printf("⚠️ No handlers registered for event type: %s", event.EventType)
		return
	}

	for _, handler := range handlers {
		ep.wg.Add(1)
		go func(h EventHandler) {
			defer ep.wg.Done()
			if err := h.Run(ctx, event); err != nil {
				log.Printf("❌ Error processing event %s: %v", event.EventType, err)
			}
		}(handler)
	}
}

// shutdown ensures all ongoing event processing completes before shutting down.
func (ep *EventProcessor) shutdown() {
	log.Println("⚠️ Draining event queue before shutdown...")
	close(ep.eventQueue)
	ep.wg.Wait()
	log.Println("✅ Event Processor shutdown complete.")
}
