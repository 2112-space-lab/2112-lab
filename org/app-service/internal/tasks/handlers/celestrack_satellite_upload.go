package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
)

type SatelliteServiceClient interface {
	FetchAndStoreAllSatellites(ctx context.Context, maxCount int) ([]domain.Satellite, error)
}

type CelesTrackSatelliteUploadHandler struct {
	satelliteRepo    domain.SatelliteRepository
	satelliteService SatelliteServiceClient
}

func NewCelesTrackSatelliteUploadHandler(
	satelliteRepo domain.SatelliteRepository,
	satelliteService SatelliteServiceClient) CelesTrackSatelliteUploadHandler {
	return CelesTrackSatelliteUploadHandler{
		satelliteRepo:    satelliteRepo,
		satelliteService: satelliteService,
	}
}

func (h *CelesTrackSatelliteUploadHandler) GetTask() Task {
	return Task{
		Name:         "celestrack_satellite_upload",
		Description:  "Fetch Satellite Metadata from CelesTrak and upsert it in the database",
		RequiredArgs: []string{"maxCount"},
	}
}

func (h *CelesTrackSatelliteUploadHandler) Run(ctx context.Context, args map[string]string) error {

	nbTles, ok := args["maxCount"]
	if !ok || nbTles == "" {
		return fmt.Errorf("missing required argument: maxCount")
	}

	// Convert nbTles to an integer (assuming it's a string or a similar type in args)
	maxCount, err := strconv.Atoi(nbTles)
	if err != nil {
		return fmt.Errorf("invalid value for max: %v", err)
	}

	_, err = h.satelliteService.FetchAndStoreAllSatellites(ctx, maxCount)
	if err != nil {
		return err
	}

	return nil
}
