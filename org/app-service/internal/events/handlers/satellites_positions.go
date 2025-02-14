package event_handlers

import (
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/events"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

// SatellitePositionHandler listens for SATELLITE_POSITION_UPDATED events.
type SatellitePositionHandler struct {
	events.BaseHandler[model.SatellitePosition]
}

// NewSatellitePositionHandler creates a new handler instance.
func NewSatellitePositionHandler() *SatellitePositionHandler {
	return &SatellitePositionHandler{}
}

// Run processes the SATELLITE_POSITION_UPDATED event.
func (h *SatellitePositionHandler) Run(event model.EventRoot) error {
	log.Infof("📡 Processing SATELLITE_POSITION_UPDATED event: UID=%s", event.EventUID)

	// Parse the payload
	payload, err := h.Parse(event.Payload)
	if err != nil {
		log.Errorf("❌ Failed to parse payload for SATELLITE_POSITION_UPDATED: %v", err)
		return err
	}

	// Handle the satellite position update
	return h.HandleSatellitePosition(event, payload)
}

// HandleSatellitePosition processes the parsed satellite position data.
func (h *SatellitePositionHandler) HandleSatellitePosition(event model.EventRoot, payload *model.SatellitePosition) error {
	if payload == nil {
		return fmt.Errorf("payload is nil for event UID: %s", event.EventUID)
	}

	log.Infof("📍 Satellite Position Updated: %s | Lat: %.6f | Lon: %.6f | Alt: %.2f km",
		payload.Name, payload.Latitude, payload.Longitude, payload.Altitude)

	return nil
}
