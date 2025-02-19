package services

import (
	"context"
	"fmt"
	"time"

	propagator "github.com/org/2112-space-lab/org/app-service/internal/clients/propagate"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xspace"
)

type SatelliteService struct {
	tleRepo          repository.TleRepository
	propagateClient  *propagator.PropagatorClient
	celestrackClient celestrackClient
	repo             repository.SatelliteRepository
	globalPropRepo   repository.GlobalPropertyRepository
}

// NewSatelliteService creates a new instance of SatelliteService.
func NewSatelliteService(tleRepo repository.TleRepository, propagateClient *propagator.PropagatorClient, celestrackClient celestrackClient, repo repository.SatelliteRepository) SatelliteService {
	return SatelliteService{tleRepo: tleRepo, propagateClient: propagateClient, celestrackClient: celestrackClient, repo: repo}
}

func (s *SatelliteService) Propagate(ctx context.Context, spaceID string, duration time.Duration, interval time.Duration) (pos []xspace.SatellitePosition, err error) {
	ctx, span := tracing.NewSpan(ctx, "Propagate")
	defer span.EndWithError(err)
	if spaceID == "" {
		return nil, fmt.Errorf("SPACE ID is required")
	}
	if duration <= 0 || interval <= 0 {
		return nil, fmt.Errorf("invalid duration or interval: both must be greater than zero")
	}

	tle, err := s.tleRepo.GetTle(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLE data for SPACE ID %s: %w", spaceID, err)
	}

	startTime := time.Now().UTC()

	resultChan, errorChan := s.propagateClient.FetchPropagation(
		ctx,
		tle.Line1,
		tle.Line2,
		startTime.Format(time.RFC3339),
		int(duration.Minutes()),
		int(interval.Seconds()),
		spaceID,
	)

	select {
	case propagatedPositions := <-resultChan:
		if propagatedPositions == nil {
			return nil, fmt.Errorf("received nil propagated positions for SPACE ID %s", spaceID)
		}

		if len(propagatedPositions) > 0 {
			firstPos := propagatedPositions[0]
			lastPos := propagatedPositions[len(propagatedPositions)-1]

			log.Tracef("First Position for SPACE ID %s: Latitude: %f, Longitude: %f, Altitude: %f, Time: %s",
				spaceID, firstPos.Latitude, firstPos.Longitude, firstPos.Altitude, firstPos.Time)

			log.Tracef("Last Position for SPACE ID %s: Latitude: %f, Longitude: %f, Altitude: %f, Time: %s",
				spaceID, lastPos.Latitude, lastPos.Longitude, lastPos.Altitude, lastPos.Time)
		}

		var positions []xspace.SatellitePosition
		for _, pos := range propagatedPositions {
			parsedTime, err := time.Parse(time.RFC3339, pos.Time)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time %s for SPACE ID %s: %w", pos.Time, spaceID, err)
			}

			positions = append(positions, xspace.SatellitePosition{
				Latitude:  pos.Latitude,
				Longitude: pos.Longitude,
				Altitude:  pos.Altitude,
				Time:      parsedTime,
			})
		}
		return positions, nil

	case err := <-errorChan:
		if err != nil {
			return nil, fmt.Errorf("failed to fetch propagated positions for SPACE ID %s: %w", spaceID, err)
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("operation canceled for SPACE ID %s: %w", spaceID, ctx.Err())
	}

	return nil, fmt.Errorf("unexpected end of Propagate function for SPACE ID %s", spaceID)
}

// GetSatelliteBySpaceID retrieves a satellite by SPACE ID.
func (s *SatelliteService) GetSatelliteBySpaceID(ctx context.Context, spaceID string) (satellite domain.Satellite, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetSatelliteBySpaceID")
	defer span.EndWithError(err)
	return s.repo.FindBySpaceID(ctx, spaceID)
}

// ListAllSatellites retrieves all stored satellites.
func (s *SatelliteService) ListAllSatellites(ctx context.Context) (satellite []domain.Satellite, err error) {
	ctx, span := tracing.NewSpan(ctx, "ListAllSatellites")
	defer span.EndWithError(err)
	return s.repo.FindAll(ctx)
}
func (s *SatelliteService) FetchAndStoreAllSatellites(ctx context.Context, maxCount int) (satellite []domain.Satellite, err error) {
	ctx, span := tracing.NewSpan(ctx, "FetchAndStoreAllSatellites")
	defer span.EndWithError(err)
	rawSatellites, err := s.celestrackClient.FetchSatelliteMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch satellite metadata: %w", err)
	}

	if len(rawSatellites) == 0 {
		return nil, fmt.Errorf("no satellite metadata available")
	}

	var storedSatellites []domain.Satellite
	for _, rawSatellite := range rawSatellites {

		satellite, err := domain.NewSatelliteFromParameters(
			rawSatellite.Name,
			rawSatellite.SpaceID,
			domain.Other,
			&rawSatellite.LaunchDate,
			rawSatellite.DecayDate,
			rawSatellite.IntlDesignator,
			rawSatellite.Owner,
			rawSatellite.ObjectType,
			rawSatellite.Period,
			rawSatellite.Inclination,
			rawSatellite.Apogee,
			rawSatellite.Apogee,
			rawSatellite.RCS,
			rawSatellite.Altitude,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create satellite for SPACE ID %s: %w", rawSatellite.SpaceID, err)
		}
		storedSatellites = append(storedSatellites, satellite)
	}

	if len(storedSatellites) > maxCount {
		storedSatellites = storedSatellites[:maxCount] // Slice to keep only the first maxCount elements
	}

	if err := s.repo.SaveBatch(ctx, storedSatellites); err != nil {
		return nil, fmt.Errorf("failed to save satellite to database: %w", err)
	}

	return storedSatellites, nil
}

// ListSatellitesWithPaginationAndTLE retrieves satellites with pagination and includes a flag indicating if a TLE is present.
func (s *SatelliteService) ListSatellitesWithPagination(ctx context.Context, page int, pageSize int, search *domain.SearchRequest) (satellite []domain.Satellite, count int64, err error) {
	ctx, span := tracing.NewSpan(ctx, "ListSatellitesWithPagination")
	defer span.EndWithError(err)
	if page <= 0 {
		return nil, 0, fmt.Errorf("page must be greater than 0")
	}
	if pageSize <= 0 {
		return nil, 0, fmt.Errorf("pageSize must be greater than 0")
	}
	satellites, totalRecords, err := s.repo.FindAllWithPagination(ctx, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve satellites with paginations: %w", err)
	}

	return satellites, totalRecords, nil
}

// ListSatelliteInfoWithPagination retrieves SatelliteInfo objects with pagination.
func (s *SatelliteService) ListSatelliteInfoWithPagination(ctx context.Context, page int, pageSize int, search *domain.SearchRequest) (satellite []domain.SatelliteInfo, count int64, err error) {
	ctx, span := tracing.NewSpan(ctx, "ListSatelliteInfoWithPagination")
	defer span.EndWithError(err)
	if page <= 0 {
		return nil, 0, fmt.Errorf("page must be greater than 0")
	}
	if pageSize <= 0 {
		return nil, 0, fmt.Errorf("pageSize must be greater than 0")
	}
	satelliteInfos, totalRecords, err := s.repo.FindSatelliteInfoWithPagination(ctx, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve satellite info with pagination: %w", err)
	}

	return satelliteInfos, totalRecords, nil
}

// GetAndLockSatellites fetches and locks satellites for processing
func (s *SatelliteService) GetAndLockSatellites(ctx context.Context, processorName string) (keys []string, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetAndLockSatellites")
	defer span.EndWithError(err)

	maxNbSatellites, err := s.globalPropRepo.GetMaxSatellitesPerEventDetector(ctx, repository.MaxSatellitesPerEventDetectorDefault)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch global property: %w", err)
	}

	lockedSatellites, err := s.repo.FetchAndLockSatellites(ctx, processorName, maxNbSatellites)
	if err != nil {
		return nil, fmt.Errorf("failed to lock satellites: %w", err)
	}

	if len(lockedSatellites) == 0 {
		log.Warn("No satellites available for locking")
		return nil, nil
	}

	log.Infof("Successfully locked %d satellites", len(lockedSatellites))
	return lockedSatellites, nil
}

// UnlockSatellites releases the locks after processing is complete
func (s *SatelliteService) UnlockSatellites(ctx context.Context, satelliteIDs []string) (err error) {
	ctx, span := tracing.NewSpan(ctx, "UnlockSatellites")
	defer span.EndWithError(err)
	return s.repo.ReleaseSatellites(ctx, satelliteIDs)
}
