package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/redis"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
)

type SatellitesTilesMappingsHandler struct {
	tileRepo      domain.TileRepository
	tleRepo       repository.TleRepository
	satelliteRepo domain.SatelliteRepository
	mappingRepo   domain.MappingRepository
	redisClient   *redis.RedisClient
	workerCount   int
}

// NewSatellitesTilesMappingsHandler creates a new instance of the handler.
func NewSatellitesTilesMappingsHandler(
	tileRepo domain.TileRepository,
	tleRepo repository.TleRepository,
	satelliteRepo domain.SatelliteRepository,
	mappingRepo domain.MappingRepository,
	redisClient *redis.RedisClient,
	workerCount int, // Number of workers
) SatellitesTilesMappingsHandler {
	return SatellitesTilesMappingsHandler{
		tileRepo:      tileRepo,
		tleRepo:       tleRepo,
		satelliteRepo: satelliteRepo,
		mappingRepo:   mappingRepo,
		redisClient:   redisClient,
		workerCount:   workerCount,
	}
}

func (h *SatellitesTilesMappingsHandler) GetTask() Task {
	return Task{
		Name:         "satellites_tiles_mappings",
		Description:  "Computes satellite visibilities for all tiles by satellite path",
		RequiredArgs: []string{},
	}
}

// Run executes the visibility computation process.
func (h *SatellitesTilesMappingsHandler) Run(ctx context.Context, args map[string]string) error {
	log.Println("Starting Run method")
	log.Println("Subscribing to event_satellite_positions_updated channel")
	return h.Subscribe(ctx, "event_satellite_positions_updated")
}

// Exec executes the visibility computation process, considering satellite paths.
func (h *SatellitesTilesMappingsHandler) Exec(ctx context.Context, id string, startTime time.Time, endTime time.Time) error {
	log.Printf("Starting Exec method for satellite ID: %s, from %s to %s\n", id, startTime, endTime)
	sat, err := h.satelliteRepo.FindByNoradID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to fetch satellite: %w", err)
	}

	positions, err := h.tleRepo.QuerySatellitePositions(ctx, sat.NoradID, startTime, endTime)
	if err != nil {
		return fmt.Errorf("error querying satellite positions for satellite %s: %w", sat.NoradID, err)
	}

	if len(positions) < 2 {
		log.Printf("Not enough positions to compute mappings for satellite %s\n", sat.NoradID)
		return nil
	}

	log.Printf("Computing mappings for satellite %s\n", sat.NoradID)
	if err := h.computeTileMappings(ctx, "todoSatellitesTilesMappingsHandler", sat, positions); err != nil {
		return fmt.Errorf("error computing mappings for satellite %s: %w", sat.NoradID, err)
	}

	log.Printf("Completed Exec method for satellite ID: %s\n", id)
	return nil
}

// computeTileMappings computes visibility for a list of satellite positions.
func (h *SatellitesTilesMappingsHandler) computeTileMappings(
	ctx context.Context,
	contextID string,
	sat domain.Satellite,
	positions []domain.SatellitePosition,
) error {
	log.Printf("Finding visible tiles for satellite %s along its path\n", sat.NoradID)

	err := h.mappingRepo.DeleteMappingsByNoradID(ctx, contextID, sat.NoradID)
	if err != nil {
		return fmt.Errorf("failed to delete visible tiles along the path: %w", err)
	}

	mappings, err := h.tileRepo.FindTilesVisibleFromLine(ctx, sat, positions)
	if err != nil {
		return fmt.Errorf("failed to find visible tiles along the path: %w", err)
	}

	if len(mappings) == 0 {
		log.Printf("No visible tiles found for satellite %s along its path\n", sat.NoradID)
		return nil
	}

	if err := h.mappingRepo.SaveBatch(ctx, mappings); err != nil {
		return fmt.Errorf("failed to save mappings: %w", err)
	}
	log.Printf("Saved %d mappings for satellite %s\n", len(mappings), sat.NoradID)

	return nil
}

// Subscribe listens for satellite position updates and computes visibility using a worker pool.
func (h *SatellitesTilesMappingsHandler) Subscribe(ctx context.Context, channel string) error {
	log.Printf("Subscribing to Redis channel: %s\n", channel)

	// Create a channel for incoming messages
	messageChan := make(chan string, 100)

	// Start worker pool
	for i := 0; i < h.workerCount; i++ {
		go h.worker(ctx, messageChan)
	}

	// Subscribe to the Redis channel
	err := h.redisClient.Subscribe(ctx, channel, func(message string) error {
		select {
		case messageChan <- message:
			// Successfully passed message to the channel
		case <-ctx.Done():
			// Context canceled
			return ctx.Err()
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %w", channel, err)
	}

	log.Printf("Subscribed to Redis channel: %s\n", channel)
	return nil
}

// worker processes incoming messages in the messageChan
func (h *SatellitesTilesMappingsHandler) worker(ctx context.Context, messageChan <-chan string) {
	for {
		select {
		case message := <-messageChan:
			var update struct {
				SatelliteID string `json:"satellite_id"`
				StartTime   string `json:"start_time"`
				EndTime     string `json:"end_time"`
			}
			if err := json.Unmarshal([]byte(message), &update); err != nil {
				log.Printf("Failed to parse update message: %v\n", err)
				continue
			}

			startTime, err := time.Parse(time.RFC3339, update.StartTime)
			if err != nil {
				log.Printf("Failed to parse start time: %v\n", err)
				continue
			}
			endTime, err := time.Parse(time.RFC3339, update.EndTime)
			if err != nil {
				log.Printf("Failed to parse end time: %v\n", err)
				continue
			}

			log.Printf("Processing update for satellite ID: %s, from %s to %s\n", update.SatelliteID, startTime, endTime)
			if err := h.Exec(ctx, update.SatelliteID, startTime, endTime); err != nil {
				log.Printf("Failed to execute computation for satellite ID %s: %v\n", update.SatelliteID, err)
			}

		case <-ctx.Done():
			log.Println("Worker shutting down due to context cancellation")
			return
		}
	}
}
