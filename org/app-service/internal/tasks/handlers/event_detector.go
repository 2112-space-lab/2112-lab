package handlers

import (
	"context"
	"log"

	"github.com/org/2112-space-lab/org/app-service/internal/events"
	event_handlers "github.com/org/2112-space-lab/org/app-service/internal/events/handlers"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
)

// EventDetector is a task handler that listens indefinitely for events
type EventDetector struct {
	eventEmitter *events.EventEmitter
	eventMonitor *events.EventMonitor
}

// NewEventDetector creates a new Event Detector Task
func NewEventDetector(
	eventEmitter *events.EventEmitter,
	eventMonitor *events.EventMonitor,
) EventDetector {
	return EventDetector{
		eventEmitter: eventEmitter,
		eventMonitor: eventMonitor,
	}
}

// GetTask returns task details
func (d *EventDetector) GetTask() Task {
	return Task{
		Name:         "event_detector",
		Description:  "Monitors events and processes TLE propagated events",
		RequiredArgs: []string{},
	}
}

// Run starts monitoring events indefinitely
func (d *EventDetector) Run(ctx context.Context, args map[string]string) error {
	log.Println("🔄 Event Detector started. Listening for events indefinitely...")

	positionsUpdatedHandler := event_handlers.NewSatellitePositionHandler()
	d.eventMonitor.RegisterHandler(model.EventTypeSatellitePositionUpdated, positionsUpdatedHandler)

	d.eventMonitor.StartMonitoring(ctx)

	return nil
}
