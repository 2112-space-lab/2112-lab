package event_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clients "github.com/org/2112-space-lab/org/app-service/internal/clients/redis"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// SatellitePositionHandler listens for SATELLITE_TLE_PROPAGATED events.
type SatellitePositionHandler struct {
	events.BaseHandler[model.SatelliteTlePropagated]
	satelliteService     services.SatelliteService
	globalRepo           repository.GlobalPropertyRepository
	eventEmitter         *events.EventEmitter
	redisClient          *clients.RedisClient
	mSatellitesPositions map[domain.SatelliteID][]model.SatellitePosition
}

// NewSatellitePositionHandler creates a new handler instance.
func NewSatellitePositionHandler(
	satelliteService services.SatelliteService,
	eventEmitter *events.EventEmitter,
	redisClient *clients.RedisClient,
) *SatellitePositionHandler {
	return &SatellitePositionHandler{
		satelliteService:     satelliteService,
		eventEmitter:         eventEmitter,
		redisClient:          redisClient,
		mSatellitesPositions: make(map[domain.SatelliteID][]model.SatellitePosition),
	}
}

func (h *SatellitePositionHandler) HandlerName() string {
	return "SatellitePositionHandler"
}

// Run processes the SATELLITE_TLE_PROPAGATED event.
func (h *SatellitePositionHandler) Run(ctx context.Context, event model.EventRoot) (err error) {
	ctx, span := tracing.NewSpan(ctx, "Run")
	defer span.EndWithError(err)

	log.Infof("📡 Processing SatelliteTlePropagated event: UID=%s", event.EventUID)

	payload, err := h.Parse(event.Payload)
	if err != nil {
		log.Errorf("❌ Failed to parse payload for SatelliteTlePropagated: %v", err)
		return err
	}

	return h.HandleSatellitePositionEvent(ctx, event, payload)
}

// HandleSatellitePositionEvent fetches positions from Redis, updates the in-memory list, and starts simulation.
func (h *SatellitePositionHandler) HandleSatellitePositionEvent(ctx context.Context, event model.EventRoot, payload *model.SatelliteTlePropagated) (err error) {
	ctx, span := tracing.NewSpan(ctx, "HandleSatellitePositionEvent")
	defer span.EndWithError(err)

	log.Infof("🔍 Fetching positions from Redis for key: %s", payload.RedisKey)
	positionsJSON, err := h.redisClient.Get(ctx, payload.RedisKey)
	if err != nil {
		log.Errorf("❌ Failed to fetch positions from Redis: %v", err)
		return err
	}

	var newPositions []model.SatellitePosition
	err = json.Unmarshal([]byte(positionsJSON), &newPositions)
	if err != nil {
		log.Errorf("❌ Failed to parse satellite positions JSON: %v", err)
		return err
	}

	bufferDuration, err := h.globalRepo.GetEventDetectorSimulationBufferDuration(ctx, repository.DefaultSimulationBufferDuration)
	if err != nil {
		log.Errorf("❌ Failed to fetch GetEventDetectorSimulationBufferDuration. Using default: %s", bufferDuration)
	}

	cutoffTime := time.Now().UTC().Add(-bufferDuration)
	satelliteID := domain.SatelliteID(payload.SpaceID)
	existingPositions, exists := h.mSatellitesPositions[satelliteID]
	if exists {
		existingPositions = append(existingPositions, newPositions...)
		filteredPositions := []model.SatellitePosition{}
		for _, pos := range existingPositions {

			posTimeUtc, err := xtime.FromString(xtime.DateTimeFormat(pos.Timestamp))
			if err != nil {
				return err
			}

			if posTimeUtc.Inner().After(cutoffTime) {
				filteredPositions = append(filteredPositions, pos)
			}
		}

		h.mSatellitesPositions[satelliteID] = filteredPositions
	} else {
		filteredPositions := []model.SatellitePosition{}

		for _, pos := range newPositions {
			posTimeUtc, err := xtime.FromString(xtime.DateTimeFormat(pos.Timestamp))
			if err != nil {
				return err
			}
			if posTimeUtc.Inner().After(cutoffTime) {
				filteredPositions = append(filteredPositions, pos)
			}
		}
		h.mSatellitesPositions[satelliteID] = filteredPositions
	}

	log.Infof("✅ Stored %d positions in-memory for satellite %s (Buffer Duration: %s)", len(h.mSatellitesPositions[satelliteID]), payload.RedisKey, bufferDuration)

	startTimeUtc, err := xtime.FromString(xtime.DateTimeFormat(payload.StartTimeUtc))
	if err != nil {
		return err
	}

	timeInterval := time.Duration(*payload.IntervalSeconds * int32(time.Second))
	simulationDuration := time.Duration(*payload.DurationMinutes * int32(time.Second))
	startIndex := 0

	go h.startSimulation(ctx, payload.RedisKey, startTimeUtc, timeInterval, simulationDuration, startIndex, h.mSatellitesPositions[satelliteID])

	return nil
}

// startSimulation simulates the movement of the satellite over time.
func (h *SatellitePositionHandler) startSimulation(ctx context.Context, satelliteKey string, startTime xtime.UtcTime, simulationInterval time.Duration, simulationDuration time.Duration, positionIndex int, positions []model.SatellitePosition) (err error) {
	ctx, span := tracing.NewSpan(ctx, "startSimulation")
	defer span.EndWithError(err)

	log.Infof("🛰 Starting simulation for satellite %s", satelliteKey)

	simulationSteps, err := h.globalRepo.GetEventDetectorSimulationSteps(ctx, repository.DefaultSimulationSteps)
	if err != nil {
		log.Errorf("❌ Failed to fetch GetEventDetectorSimulationInterval take default value: %t", simulationSteps)
	}

	totalPositions := len(positions)

	for time.Since(startTime.Inner()) < simulationDuration {
		if positionIndex >= totalPositions {
			log.Warnf("⏳ No more positions to simulate for %s", satelliteKey)
			break
		}

		currentPosition := positions[positionIndex]
		log.Tracef("📍 Satellite %s - Position: Lat=%.6f, Lon=%.6f, Alt=%.2f", satelliteKey, currentPosition.Latitude, currentPosition.Longitude, currentPosition.Altitude)

		eventJSON, err := json.Marshal(currentPosition)
		if err != nil {
			log.Errorf("❌ Failed to parse satellite positions JSON: %v", err)
			return err
		}

		h.eventEmitter.PublishEvent(ctx, model.EventRoot{
			EventType: string(model.EventTypeSatellitePositionUpdated),
			EventUID:  fmt.Sprintf("%s-%d", satelliteKey, positionIndex),
			Payload:   string(eventJSON),
		})

		time.Sleep(simulationInterval)
		positionIndex += int(simulationSteps)
	}

	log.Infof("✅ Simulation completed for satellite %s", satelliteKey)
	return nil
}
