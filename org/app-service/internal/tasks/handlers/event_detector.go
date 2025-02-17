package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/rabbitmq"
	"github.com/org/2112-space-lab/org/app-service/internal/dependencies"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	event_handlers "github.com/org/2112-space-lab/org/app-service/internal/events/handlers"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
)

// EventDetector is a task handler that listens indefinitely for events
type EventDetector struct {
	eventEmitter *events.EventEmitter
	eventMonitor *events.EventMonitor
	dependencies *dependencies.Dependencies
}

// NewEventDetector creates a new Event Detector Task
func NewEventDetector(
	eventEmitter *events.EventEmitter,
	eventMonitor *events.EventMonitor,
	dependencies *dependencies.Dependencies,
) EventDetector {
	return EventDetector{
		eventEmitter: eventEmitter,
		eventMonitor: eventMonitor,
		dependencies: dependencies,
	}
}

// GetTask returns task details
func (d *EventDetector) GetTask() Task {
	return Task{
		Name:         "event_detector",
		Description:  "Monitors events and processes TLE propagated events",
		RequiredArgs: []string{"name"},
	}
}

// Run starts monitoring events indefinitely
func (d *EventDetector) Run(ctx context.Context, args map[string]string) error {
	log.Println("🔄 Event Detector started. Listening for events indefinitely...")

	processorName, ok := args["name"]
	if !ok || processorName == "" {
		return fmt.Errorf("missing required argument: name")
	}

	positionsUpdatedHandler := event_handlers.NewSatellitePositionHandler(d.dependencies.Services.SatelliteService, d.eventEmitter, d.dependencies.Clients.RedisClient)

	satelliteKeys, err := d.dependencies.Services.SatelliteService.GetAndLockSatellites(ctx, processorName)
	if err != nil {
		return err
	}

	d.eventMonitor.RegisterHandler(ctx, model.EventTypeSatelliteTlePropagated, positionsUpdatedHandler)

	header := rabbitmq.NewHeader()
	for _, s := range satelliteKeys {
		header.AddField("satellite_id", s)
	}

	d.eventMonitor.StartMonitoring(ctx, header)

	return nil
}
