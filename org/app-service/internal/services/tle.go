package services

import (
	"context"
	"fmt"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/api/mappers"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

type celestrackClient interface {
	FetchTLEFromSatCatByCategory(ctx context.Context, category string) ([]*mappers.RawTLE, error)
	FetchSatelliteMetadata(ctx context.Context) ([]*mappers.SatelliteMetadata, error)
}

type TleService struct {
	celestrackClient celestrackClient
	tleRepo          repository.TleRepository
	contextRepo      domain.GameContextRepository
}

// NewTleService creates a new instance of TleService.
func NewTleService(
	celestrackClient celestrackClient,
	tleRepo repository.TleRepository,
	contextRepo domain.GameContextRepository,
) TleService {
	return TleService{
		celestrackClient: celestrackClient,
		tleRepo:          tleRepo,
		contextRepo:      contextRepo,
	}
}

// FetchTLEFromSatCatByCategory fetches TLEs from a given category and associates them with a context.
func (s *TleService) FetchTLEFromSatCatByCategory(ctx context.Context, category string, contextName domain.GameContextName) (ts []domain.TLE, err error) {
	ctx, span := tracing.NewSpan(ctx, "FetchTLEFromSatCatByCategory")
	defer span.EndWithError(err)
	// Validate the contextID
	if _, err := s.contextRepo.FindByUniqueName(ctx, contextName); err != nil {
		return nil, fmt.Errorf("invalid contextID: %w", err)
	}

	nowUtc := time.Now().UTC()

	// Fetch raw TLEs from the external service
	rawTLEs, err := s.celestrackClient.FetchTLEFromSatCatByCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLEs from category [%s]: %w", category, err)
	}

	tles := make([]domain.TLE, len(rawTLEs))
	for idx, raw := range rawTLEs {
		tle, err := domain.NewTLE(
			raw.NoradID,
			raw.Line1,
			raw.Line2,
			nowUtc,
			string(contextName), // Associate with the context
			true,
			false,
		)

		if err != nil {
			return nil, fmt.Errorf("error creating TLE for NORAD ID [%s]: %w", raw.NoradID, err)
		}
		tles[idx] = tle
	}

	return tles, nil
}

// FetchSatelliteMetadata retrieves metadata about satellites and associates them with a context.
func (s *TleService) FetchSatelliteMetadata(ctx context.Context, contextName domain.GameContextName) (sats []domain.Satellite, err error) {
	ctx, span := tracing.NewSpan(ctx, "FetchSatelliteMetadata")
	defer span.EndWithError(err)
	// Validate the contextID
	if _, err := s.contextRepo.FindByUniqueName(ctx, contextName); err != nil {
		return nil, fmt.Errorf("invalid contextID: %w", err)
	}

	// Fetch satellite metadata from the external client
	metadata, err := s.celestrackClient.FetchSatelliteMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch satellite metadata: %w", err)
	}

	nowUtc := time.Now().UTC()
	satellites := make([]domain.Satellite, len(metadata))
	for idx, raw := range metadata {
		sat := domain.Satellite{
			NoradID:    raw.NoradID,
			Name:       raw.Name,
			Owner:      raw.Owner,
			LaunchDate: &raw.LaunchDate,
			DecayDate:  raw.DecayDate,
			ObjectType: raw.ObjectType,
			ModelBase: domain.ModelBase{
				CreatedAt: nowUtc,
			},
		}
		satellites[idx] = sat
	}

	return satellites, nil
}
