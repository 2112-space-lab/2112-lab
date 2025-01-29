package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
)

type TleServiceClient interface {
	FetchTLEFromSatCatByCategory(ctx context.Context, category string, contextName domain.GameContextName) ([]domain.TLE, error)
}

type CelestrackTleUploadHandler struct {
	satelliteRepo domain.SatelliteRepository
	tleRepo       repository.TleRepository
	tleService    TleServiceClient
}

func NewCelestrackTleUploadHandler(
	satelliteRepo domain.SatelliteRepository,
	tleRepo repository.TleRepository,
	tleService TleServiceClient) CelestrackTleUploadHandler {
	return CelestrackTleUploadHandler{
		satelliteRepo: satelliteRepo,
		tleRepo:       tleRepo,
		tleService:    tleService,
	}
}

func (h *CelestrackTleUploadHandler) GetTask() Task {
	return Task{
		Name:         "celestrack_tle_upload",
		Description:  "Fetch TLE from CelesTrak and upsert it in the database",
		RequiredArgs: []string{"category", "maxCount", "contextName"},
	}
}

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
	if !ok || nbTles == "" {
		return fmt.Errorf("missing required argument: maxCount")
	}

	// Convert nbTles to an integer (assuming it's a string or a similar type in args)
	maxCount, err := strconv.Atoi(nbTles)
	if err != nil {
		return fmt.Errorf("invalid value for max: %v", err)
	}

	tles, err := h.tleService.FetchTLEFromSatCatByCategory(ctx, category, domain.GameContextName(contextName))
	if err != nil {
		return fmt.Errorf("failed to fetch TLE catalog for category %s: %v", category, err)
	}

	// Retain only the maximum nbTles elements
	if len(tles) > maxCount {
		tles = tles[:maxCount] // Slice to keep only the first maxCount elements
	}

	// Ensure satellites exist for the batch
	// for _, tle := range tles {
	// 	if err := h.ensureSatelliteExists(ctx, tle.NoradID, category); err != nil {
	// 		return fmt.Errorf("failed to ensure satellite existence for NORAD ID %s: %v", tle.NoradID, err)
	// 	}
	// }

	log.Printf("Returning %d TLEs for category %s", len(tles), category)
	err = h.tleRepo.UpdateTleBatch(ctx, tles)
	if err != nil {
		return fmt.Errorf("failed to upsert TLE for NORAD ID %s", err)
	}
	return nil
}

func (h *CelestrackTleUploadHandler) ensureSatelliteExists(ctx context.Context, noradID string, category string) error {
	satellite, err := h.satelliteRepo.FindByNoradID(ctx, noradID)
	if err == nil && satellite.NoradID == noradID {
		return nil
	}

	nowUtc := time.Now().UTC()

	newSatellite, err := domain.NewSatellite(
		fmt.Sprintf("Unknown Satellite %s", noradID),
		noradID,
		domain.SatelliteType(strings.ToUpper(category)),
		true,  //active by default
		false, // not in favourite
		nowUtc,
	)

	if err != nil {
		return err
	}

	err = h.satelliteRepo.Save(ctx, newSatellite)
	if err != nil {
		return fmt.Errorf("failed to create satellite: %v", err)
	}

	return nil
}
