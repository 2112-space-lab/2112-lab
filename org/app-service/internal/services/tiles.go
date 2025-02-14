package services

import (
	"context"
	"fmt"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

type TileService struct {
	repo          repository.TileRepository
	tleRepo       repository.TleRepository
	satelliteRepo repository.SatelliteRepository
	mappingRepo   repository.TileSatelliteMappingRepository
}

// NewTileService creates a new instance of TileService.
func NewTileService(
	tileRepo repository.TileRepository,
	tleRepo repository.TleRepository,
	satelliteRepo repository.SatelliteRepository,
	mappingRepo repository.TileSatelliteMappingRepository,
) TileService {
	return TileService{
		repo:          tileRepo,
		tleRepo:       tleRepo,
		satelliteRepo: satelliteRepo,
		mappingRepo:   mappingRepo,
	}
}

// FindAllTiles retrieves all tiles associated with a specific context.
func (s *TileService) FindAllTiles(ctx context.Context, contextID string) (t []domain.Tile, err error) {
	ctx, span := tracing.NewSpan(ctx, "FindAllTiles")
	defer span.EndWithError(err)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	tiles, err := s.repo.FindTilesInRegion(ctx, contextID, -90, -180, 90, 180) // World bounding box
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tiles for context [%s]: %w", contextID, err)
	}

	if len(tiles) == 0 {
		return nil, fmt.Errorf("no tiles found for context [%s]", contextID)
	}

	return tiles, nil
}

// GetTilesInRegion fetches tiles that intersect with a bounding box and belong to a specific context.
func (s *TileService) GetTilesInRegion(ctx context.Context, contextID string, minLat, minLon, maxLat, maxLon float64) (t []domain.Tile, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetTilesInRegion")
	defer span.EndWithError(err)
	// Validate input
	if minLat >= maxLat || minLon >= maxLon {
		return nil, fmt.Errorf("invalid bounding box coordinates")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	tiles, err := s.repo.FindTilesInRegion(ctx, contextID, minLat, minLon, maxLat, maxLon)
	if err != nil {
		return nil, fmt.Errorf("error fetching tiles in region for context [%s]: %w", contextID, err)
	}

	return tiles, nil
}

// ListSatellitesMappingWithPagination retrieves mappings with pagination for a specific context.
func (s *TileService) ListSatellitesMappingWithPagination(ctx context.Context, contextID string, page int, pageSize int, search *domain.SearchRequest) (ts []domain.TileSatelliteInfo, count int64, err error) {
	ctx, span := tracing.NewSpan(ctx, "ListSatellitesMappingWithPagination")
	defer span.EndWithError(err)
	// Validate inputs
	if page <= 0 {
		return nil, 0, fmt.Errorf("page must be greater than 0")
	}
	if pageSize <= 0 {
		return nil, 0, fmt.Errorf("pageSize must be greater than 0")
	}

	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	mappings, totalRecords, err := s.mappingRepo.ListSatellitesMappingWithPagination(ctx, contextID, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve satellites mapping with pagination for context [%s]: %w", contextID, err)
	}

	return mappings, totalRecords, nil
}

// GetSatelliteMappingsByNoradID retrieves mappings for a specific NORAD ID and context.
func (s *TileService) GetSatelliteMappingsByNoradID(ctx context.Context, contextID, noradID string) (ts []domain.TileSatelliteInfo, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetSatelliteMappingsByNoradID")
	defer span.EndWithError(err)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	mappings, err := s.mappingRepo.GetSatelliteMappingsByNoradID(ctx, contextID, noradID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve satellite mappings for NORAD ID [%s] in context [%s]: %w", noradID, contextID, err)
	}

	return mappings, nil
}

// RecomputeMappings deletes existing mappings for a NORAD ID in a specific context and computes new ones.
func (s *TileService) RecomputeMappings(ctx context.Context, contextID, noradID string, startTime, endTime time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx, "RecomputeMappings")
	defer span.EndWithError(err)
	log.Tracef("Recomputing mappings for NORAD ID: %s in context: %s\n", noradID, contextID)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Step 1: Delete existing mappings
	if err := s.mappingRepo.DeleteMappingsByNoradID(ctx, contextID, noradID); err != nil {
		return fmt.Errorf("failed to delete existing mappings for NORAD ID [%s] in context [%s]: %w", noradID, contextID, err)
	}
	log.Tracef("Deleted existing mappings for NORAD ID: %s in context: %s\n", noradID, contextID)

	// Step 2: Fetch satellite data
	satellite, err := s.satelliteRepo.FindByNoradID(ctx, noradID)
	if err != nil {
		return fmt.Errorf("failed to fetch satellite for NORAD ID [%s]: %w", noradID, err)
	}

	// Step 3: Fetch satellite positions
	positions, err := s.tleRepo.QuerySatellitePositions(ctx, satellite.NoradID, startTime, endTime)
	if err != nil {
		return fmt.Errorf("failed to fetch satellite positions for NORAD ID [%s]: %w", noradID, err)
	}

	if len(positions) < 2 {
		log.Warnf("Not enough positions to compute mappings for NORAD ID: %s\n", noradID)
		return nil
	}

	// Step 4: Compute new mappings
	mappings, err := s.repo.FindTilesVisibleFromLine(ctx, satellite, positions)
	if err != nil {
		return fmt.Errorf("failed to compute tile mappings for NORAD ID [%s]: %w", noradID, err)
	}

	// Step 5: Save new mappings
	if err := s.mappingRepo.SaveBatch(ctx, mappings); err != nil {
		return fmt.Errorf("failed to save new mappings for NORAD ID [%s]: %w", noradID, err)
	}

	log.Debugf("Recomputed and saved %d mappings for NORAD ID: %s in context: %s\n", len(mappings), noradID, contextID)
	return nil
}
