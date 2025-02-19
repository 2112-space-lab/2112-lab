package handlers

import (
	"context"
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/rabbitmq"
	"github.com/org/2112-space-lab/org/app-service/internal/dependencies"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	event_handlers "github.com/org/2112-space-lab/org/app-service/internal/events/handlers"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// EventDetector is a task handler that listens indefinitely for events
type EventDetector struct {
	eventEmitter *events.EventEmitter
	eventMonitor *events.EventMonitor
	dependencies *dependencies.Dependencies
}

// NewEventDetector creates a new Event Detector Task
func NewEventDetector(
	ctx context.Context,
	eventEmitter *events.EventEmitter,
	eventMonitor *events.EventMonitor,
	dependencies *dependencies.Dependencies,
) (e EventDetector, err error) {
	_, span := tracing.NewSpan(ctx, "EventDetector.NewEventDetector")
	defer span.EndWithError(err)

	return EventDetector{
		eventEmitter: eventEmitter,
		eventMonitor: eventMonitor,
		dependencies: dependencies,
	}, nil
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
func (d *EventDetector) Run(ctx context.Context, args map[string]string) (err error) {
	ctx, span := tracing.NewSpan(ctx, "EventDetector.Run")
	defer span.EndWithError(err)

	log.Debug("🔄 Event Detector started. Listening for events indefinitely...")

	processorName, ok := args["name"]
	if !ok || processorName == "" {
		return fmt.Errorf("missing required argument: name")
	}

	positionsUpdatedHandler := event_handlers.NewSatellitePositionHandler(d.dependencies.Services.SatelliteService, d.eventEmitter, d.dependencies.Clients.RedisClient)

	satelliteKeys, err := d.dependencies.Services.SatelliteService.GetAndLockSatellites(ctx, processorName)
	if err != nil {
		return err
	}

	err = d.eventMonitor.RegisterHandler(ctx, model.EventTypeSatelliteTlePropagated, positionsUpdatedHandler)
	if err != nil {
		return err
	}

	header := rabbitmq.NewHeader()
	for _, s := range satelliteKeys {
		header.AddField("satellite_id", s)
	}

	err = d.eventMonitor.StartMonitoring(ctx, header)
	return err
}
