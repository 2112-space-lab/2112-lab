package services

import (
	"context"
	"fmt"
	"log"
	"time"

	propagator "github.com/org/2112-space-lab/org/app-service/internal/clients/propagate"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xspace"
)

type SatelliteService struct {
	tleRepo          repository.TleRepository
	propagateClient  *propagator.PropagatorClient
	celestrackClient celestrackClient
	repo             domain.SatelliteRepository
}

// NewSatelliteService creates a new instance of SatelliteService.
func NewSatelliteService(tleRepo repository.TleRepository, propagateClient *propagator.PropagatorClient, celestrackClient celestrackClient, repo domain.SatelliteRepository) SatelliteService {
	return SatelliteService{tleRepo: tleRepo, propagateClient: propagateClient, celestrackClient: celestrackClient, repo: repo}
}

func (s *SatelliteService) Propagate(ctx context.Context, noradID string, duration time.Duration, interval time.Duration) (pos []xspace.SatellitePosition, err error) {
	ctx, span := tracing.NewSpan(ctx, "Propagate")
	defer span.EndWithError(err)
	// Validate inputs
	if noradID == "" {
		return nil, fmt.Errorf("NORAD ID is required")
	}
	if duration <= 0 || interval <= 0 {
		return nil, fmt.Errorf("invalid duration or interval: both must be greater than zero")
	}

	// Get the TLE data for the satellite by NORAD ID
	tle, err := s.tleRepo.GetTle(ctx, noradID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLE data for NORAD ID %s: %w", noradID, err)
	}

	// Set up the time range
	startTime := time.Now().UTC()

	// Use PropagatorClient to propagate satellite positions
	resultChan, errorChan := s.propagateClient.FetchPropagation(
		ctx,
		tle.Line1,
		tle.Line2,
		startTime.Format(time.RFC3339),
		int(duration.Minutes()),
		int(interval.Seconds()),
		noradID,
	)

	// Wait for results or errors
	select {
	case propagatedPositions := <-resultChan:
		if propagatedPositions == nil {
			return nil, fmt.Errorf("received nil propagated positions for NORAD ID %s", noradID)
		}

		// Log the first and last positions
		if len(propagatedPositions) > 0 {
			firstPos := propagatedPositions[0]
			lastPos := propagatedPositions[len(propagatedPositions)-1]

			log.Printf("First Position for NORAD ID %s: Latitude: %f, Longitude: %f, Altitude: %f, Time: %s",
				noradID, firstPos.Latitude, firstPos.Longitude, firstPos.Altitude, firstPos.Time)

			log.Printf("Last Position for NORAD ID %s: Latitude: %f, Longitude: %f, Altitude: %f, Time: %s",
				noradID, lastPos.Latitude, lastPos.Longitude, lastPos.Altitude, lastPos.Time)
		}

		// Convert the API response to the internal SatellitePosition format
		var positions []xspace.SatellitePosition
		for _, pos := range propagatedPositions {
			// Parse the time from the API response
			parsedTime, err := time.Parse(time.RFC3339, pos.Time)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time %s for NORAD ID %s: %w", pos.Time, noradID, err)
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
			return nil, fmt.Errorf("failed to fetch propagated positions for NORAD ID %s: %w", noradID, err)
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("operation canceled for NORAD ID %s: %w", noradID, ctx.Err())
	}

	return nil, fmt.Errorf("unexpected end of Propagate function for NORAD ID %s", noradID)
}

// GetSatelliteByNoradID retrieves a satellite by NORAD ID.
func (s *SatelliteService) GetSatelliteByNoradID(ctx context.Context, noradID string) (satellite domain.Satellite, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetSatelliteByNoradID")
	defer span.EndWithError(err)
	return s.repo.FindByNoradID(ctx, noradID)
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
	// Fetch all satellite metadata from CelestrackClient
	rawSatellites, err := s.celestrackClient.FetchSatelliteMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch satellite metadata: %w", err)
	}

	if len(rawSatellites) == 0 {
		return nil, fmt.Errorf("no satellite metadata available")
	}

	var storedSatellites []domain.Satellite
	for _, rawSatellite := range rawSatellites {

		// Use the updated constructor to create a Satellite
		satellite, err := domain.NewSatelliteFromStatCat(
			rawSatellite.Name,
			rawSatellite.NoradID,
			domain.Other, // Default type; adjust based on metadata if available
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
			return nil, fmt.Errorf("failed to create satellite for NORAD ID %s: %w", rawSatellite.NoradID, err)
		}
		// Add the satellite to the result list
		storedSatellites = append(storedSatellites, satellite)
	}

	// Retain only the maximum nbTles elements
	if len(storedSatellites) > maxCount {
		storedSatellites = storedSatellites[:maxCount] // Slice to keep only the first maxCount elements
	}

	// Save the satellite to the repository
	if err := s.repo.SaveBatch(ctx, storedSatellites); err != nil {
		return nil, fmt.Errorf("failed to save satellite to database: %w", err)
	}

	return storedSatellites, nil
}

// ListSatellitesWithPaginationAndTLE retrieves satellites with pagination and includes a flag indicating if a TLE is present.
func (s *SatelliteService) ListSatellitesWithPagination(ctx context.Context, page int, pageSize int, search *domain.SearchRequest) (satellite []domain.Satellite, count int64, err error) {
	ctx, span := tracing.NewSpan(ctx, "ListSatellitesWithPagination")
	defer span.EndWithError(err)
	// Validate inputs
	if page <= 0 {
		return nil, 0, fmt.Errorf("page must be greater than 0")
	}
	if pageSize <= 0 {
		return nil, 0, fmt.Errorf("pageSize must be greater than 0")
	}

	// Fetch satellites with pagination and TLE flag
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
	// Validate inputs
	if page <= 0 {
		return nil, 0, fmt.Errorf("page must be greater than 0")
	}
	if pageSize <= 0 {
		return nil, 0, fmt.Errorf("pageSize must be greater than 0")
	}

	// Fetch SatelliteInfo with pagination
	satelliteInfos, totalRecords, err := s.repo.FindSatelliteInfoWithPagination(ctx, page, pageSize, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve satellite info with pagination: %w", err)
	}

	return satelliteInfos, totalRecords, nil
}
