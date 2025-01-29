package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/redis"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
)

// ComputeVisibilitiessHandler handles visibility computation for satellites based on user locations.
type ComputeVisibilitiessHandler struct {
	tileRepo       domain.TileRepository
	mappingRepo    domain.MappingRepository
	tleRepo        repository.TleRepository
	redisClient    *redis.RedisClient
	defaultHorizon int
}

// NewComputeVisibilitiessHandler initializes a new handler instance.
func NewComputeVisibilitiessHandler(
	tileRepo domain.TileRepository,
	mappingRepo domain.MappingRepository,
	tleRepo repository.TleRepository,
	redisClient *redis.RedisClient,
) ComputeVisibilitiessHandler {
	return ComputeVisibilitiessHandler{
		tileRepo:    tileRepo,
		tleRepo:     tleRepo,
		mappingRepo: mappingRepo,
		redisClient: redisClient,
	}
}

func (h *ComputeVisibilitiessHandler) GetTask() Task {
	return Task{
		Name:         "compute_visibilities",
		Description:  "Computes satellite visibilities for all tiles by satellite path",
		RequiredArgs: []string{"defaultHorizon"},
	}
}

// Run executes the visibility computation process.
func (h *ComputeVisibilitiessHandler) Run(ctx context.Context, args map[string]string) error {
	log.Println("Starting Run method")
	log.Println("Subscribing to visibility_requests channel")

	defaultUserHorizon, ok := args["defaultHorizon"]
	if !ok || defaultUserHorizon == "" {
		return fmt.Errorf("missing required argument: defaultUserHorizon")
	}

	horizon, err := strconv.Atoi(defaultUserHorizon)
	if err != nil {
		return fmt.Errorf("invalid value for horizon: %v", err)
	}
	h.defaultHorizon = horizon

	return h.Subscribe(ctx, "visibility_requests")
}

// Subscribe listens for user location updates and computes visibilities.
func (h *ComputeVisibilitiessHandler) Subscribe(ctx context.Context, channel string) error {
	log.Printf("Subscribing to Redis channel: %s\n", channel)

	err := h.redisClient.Subscribe(ctx, channel, func(message string) error {
		var request struct {
			UID       string  `json:"uid"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Radius    float64 `json:"radius"`
			Horizon   float64 `json:"horizon"`
			StartTime string  `json:"startTime"`
			EndTime   string  `json:"endTime"`
		}

		if err := json.Unmarshal([]byte(message), &request); err != nil {
			return fmt.Errorf("failed to parse update message: %w", err)
		}

		startTime, err := time.Parse(time.RFC3339, request.StartTime)
		if err != nil {
			return fmt.Errorf("failed to parse start time: %w", err)
		}
		endTime, err := time.Parse(time.RFC3339, request.EndTime)
		if err != nil {
			return fmt.Errorf("failed to parse end time: %w", err)
		}

		log.Printf("Received visibility request for UID: %s at location (%.6f, %.6f) with radius %.2f, horizon %.2f, from %s to %s\n",
			request.UID, request.Latitude, request.Longitude, request.Radius, request.Horizon, startTime, endTime)

		return h.computeVisibility(ctx, request.UID, "toprovideincomputeVisibility", request.Latitude, request.Longitude, request.Radius, startTime, endTime)
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %w", channel, err)
	}

	log.Printf("Successfully subscribed to Redis channel: %s\n", channel)
	return nil
}

// computeVisibility computes visibilities for a given user location and time range.
func (h *ComputeVisibilitiessHandler) computeVisibility(ctx context.Context, uid string, contextID string, latitude, longitude, radius float64, startTime, endTime time.Time) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude out of bounds: %f", latitude)
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude out of bounds: %f", longitude)
	}
	if radius <= 0 {
		return fmt.Errorf("radius must be greater than 0")
	}

	tiles, err := h.tileRepo.FindTilesIntersectingLocation(ctx, contextID, latitude, longitude, radius)
	if err != nil {
		return fmt.Errorf("failed to find tiles intersecting location: %w", err)
	}
	if len(tiles) == 0 {
		log.Printf("No tiles found for location: (%f, %f) with radius: %f\n", latitude, longitude, radius)
		return nil
	}

	var tileIDs []string
	for _, tile := range tiles {
		tileIDs = append(tileIDs, tile.ID)
	}

	satellites, err := h.mappingRepo.FindSatellitesForTiles(ctx, contextID, tileIDs)
	if err != nil {
		return fmt.Errorf("failed to find satellites for tiles: %w", err)
	}
	if len(satellites) == 0 {
		log.Printf("No satellites found for the identified tiles.\n")
		return nil
	}

	key := fmt.Sprintf("user_visibilities_event:%s", uid)
	var visibilities []map[string]interface{}
	for _, satellite := range satellites {

		// Get the TLE data for the satellite by NORAD ID
		tle, err := h.tleRepo.GetTle(ctx, satellite.NoradID)
		if err != nil {
			return fmt.Errorf("failed to fetch TLE data for NORAD ID %s: %w", satellite.NoradID, err)
		}

		visibility := map[string]interface{}{
			"satelliteID":   satellite.NoradID,
			"satelliteName": satellite.Name,
			"startTime":     startTime.Format(time.RFC3339), // start propgated period
			"endTime":       endTime.Format(time.RFC3339),   // end propgated period
			"tleLine1":      tle.Line1,
			"tleLine2":      tle.Line2,
			"userLocation": map[string]interface{}{
				"latitude":  latitude,
				"longitude": longitude,
				"radius":    radius,
				"horizon":   h.defaultHorizon,
				"uid":       uid,
			},
			"userUID": uid,
		}
		visibilities = append(visibilities, visibility)
	}

	cachedData, err := json.Marshal(visibilities)
	if err != nil {
		return fmt.Errorf("failed to serialize visibilities: %w", err)
	}

	if err := h.redisClient.Publish(ctx, key, cachedData); err != nil {
		return fmt.Errorf("failed to publish visibility event: %w", err)
	}

	return nil
}
