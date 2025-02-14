package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

// TleServiceClient defines an interface for fetching TLE data
type TleServiceClient interface {
	FetchTLEFromSatCatByCategory(ctx context.Context, category string, contextName domain.GameContextName) ([]domain.TLE, error)
}

// CelestrackTleUploadHandler handles TLE uploads from CelesTrak
type CelestrackTleUploadHandler struct {
	satelliteRepo repository.SatelliteRepository
	tleRepo       repository.TleRepository
	tleService    TleServiceClient
	eventEmitter  *events.EventEmitter
	eventsMonitor *events.EventMonitor
}

// NewCelestrackTleUploadHandler creates a new task
func NewCelestrackTleUploadHandler(
	satelliteRepo repository.SatelliteRepository,
	tleRepo repository.TleRepository,
	tleService TleServiceClient,
	eventEmitter *events.EventEmitter,
	eventsMonitor *events.EventMonitor,
) CelestrackTleUploadHandler {
	return CelestrackTleUploadHandler{
		satelliteRepo: satelliteRepo,
		tleRepo:       tleRepo,
		tleService:    tleService,
		eventEmitter:  eventEmitter,
		eventsMonitor: eventsMonitor,
	}
}

func (h *CelestrackTleUploadHandler) GetTask() Task {
	return Task{
		Name:         "celestrack_tle_upload",
		Description:  "Fetch TLE from CelesTrak and upsert it in the database",
		RequiredArgs: []string{"category", "maxCount", "contextName"},
	}
}

// Run executes the task manually or via an event
func (h *CelestrackTleUploadHandler) Run(ctx context.Context, args map[string]string) error {
	category, ok := args["category"]
	if !ok || category == "" {
		return fmt.Errorf("missing required argument: category")
	}

	nbTles, ok := args["maxCount"]
	if !ok || nbTles == "" {
		return fmt.Errorf("missing required argument: maxCount")
	}

	contextName, ok := args["contextName"]
	if !ok || contextName == "" {
		return fmt.Errorf("missing required argument: contextName")
	}

	maxCount, err := strconv.Atoi(nbTles)
	if err != nil {
		return fmt.Errorf("invalid value for maxCount: %v", err)
	}

	tles, err := h.tleService.FetchTLEFromSatCatByCategory(ctx, category, domain.GameContextName(contextName))
	if err != nil {
		return fmt.Errorf("failed to fetch TLE catalog for category %s: %v", category, err)
	}

	if len(tles) > maxCount {
		tles = tles[:maxCount]
	}

	err = h.tleRepo.UpdateTleBatch(ctx, tles)
	if err != nil {
		return fmt.Errorf("failed to upsert TLE batch: %v", err)
	}

	log.Debugf("✅ Successfully processed %d TLEs for category %s", len(tles), category)

	h.emitTleProcessedEvent(category, maxCount, len(tles))
	return nil
}

// emitTleProcessedEvent sends a completion event
func (h *CelestrackTleUploadHandler) emitTleProcessedEvent(category string, maxRequested, processed int) {
	eventPayload := map[string]interface{}{
		"category":      category,
		"maxRequested":  maxRequested,
		"processedTLEs": processed,
		"timestamp":     time.Now().UTC(),
	}

	eventData, _ := json.Marshal(eventPayload)
	event := model.EventRoot{
		EventType: model.EventTypeSatelliteTlePropagated.String(),
		Payload:   string(eventData),
	}

	if err := h.eventEmitter.PublishEvent(event); err != nil {
		log.Errorf("❌ Failed to emit event: %v", err)
	} else {
		log.Tracef("📡 Event emitted: TLE_UPLOAD_COMPLETED for category %s", category)
	}
}
