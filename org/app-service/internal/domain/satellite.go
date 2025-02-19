package domain

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xspace"
)

// SatelliteType represents the type of a satellite.
type SatelliteType string
type SatelliteID string

const (
	// Active satellite type.
	Active SatelliteType = "ACTIVE"
	// Other satellite type from SATCAT catalogue.
	Other SatelliteType = "OTHER"
)

// IsValid checks if the SatelliteType is valid.
func (t SatelliteType) IsValid() error {
	switch t {
	case Active, Other:
		return nil
	default:
		return errors.New("invalid satellite type")
	}
}

type Satellite struct {
	ModelBase
	Name                 string
	SpaceID              string
	Type                 SatelliteType
	LaunchDate           fx.Option[xtime.UtcTime] // Added field for launch date
	DecayDate            fx.Option[xtime.UtcTime] // Added field for decay date, if applicable
	IntlDesignator       string                   // Added field for international designator
	Owner                string                   // Added field for ownership information
	ObjectType           string                   // Added field for object type (e.g., PAYLOAD)
	PeriodInMinutes      fx.Option[float64]       // Added field for orbital period in minutes
	InclinationInDegrees fx.Option[float64]       // Added field for orbital inclination in degrees
	ApogeeInKm           fx.Option[float64]       // Added field for apogee altitude in kilometers
	PerigeeInKm          fx.Option[float64]       // Added field for perigee altitude in kilometers
	RCS                  fx.Option[float64]       // Added field for radar cross-section in square meters
	TleUpdatedAt         fx.Option[xtime.UtcTime] `gorm:"-"`
	Altitude             fx.Option[float64]
	OrbitType            xspace.OrbitType
}

// NewSatelliteFromParameters creates a new Satellite instance with optional SATCAT data.
func NewSatelliteFromParameters(
	name string,
	spaceID string,
	satType SatelliteType,
	launchDate *time.Time,
	decayDate *time.Time,
	intlDesignator string,
	owner string,
	objectType string,
	period *float64,
	inclination *float64,
	apogee *float64,
	perigee *float64,
	rcs *float64,
	altitude *float64,
) (Satellite, error) {
	nowUtc := time.Now().UTC()
	if err := satType.IsValid(); err != nil {
		return Satellite{}, err
	}

	return Satellite{
		ModelBase: ModelBase{
			ID:          uuid.NewString(),
			CreatedAt:   nowUtc,
			UpdatedAt:   &nowUtc,
			DisplayName: name,
			IsActive:    true,
			ProcessedAt: &nowUtc,
			IsFavourite: false,
		},
		Name:                 name,
		SpaceID:              spaceID,
		Type:                 satType,
		LaunchDate:           xtime.ConvertToUtcTime(launchDate),
		DecayDate:            xtime.ConvertToUtcTime(decayDate),
		IntlDesignator:       intlDesignator,
		Owner:                owner,
		ObjectType:           objectType,
		PeriodInMinutes:      fx.ConvertToFloatOption(period),
		InclinationInDegrees: fx.ConvertToFloatOption(inclination),
		ApogeeInKm:           fx.ConvertToFloatOption(apogee),
		PerigeeInKm:          fx.ConvertToFloatOption(perigee),
		RCS:                  fx.ConvertToFloatOption(rcs),
		Altitude:             fx.ConvertToFloatOption(altitude),
		OrbitType:            xspace.ComputeOrbitType(*altitude),
	}, nil
}

// NewSatellite creates a new Satellite instance.
func NewSatellite(name string, spaceID string, satType SatelliteType, isFavourite bool, isActive bool, createdAt time.Time) (Satellite, error) {
	if err := satType.IsValid(); err != nil {
		return Satellite{}, err
	}
	return Satellite{
		ModelBase: ModelBase{
			ID:          uuid.NewString(),
			CreatedAt:   createdAt,
			UpdatedAt:   &createdAt,
			DisplayName: name,
			IsActive:    isActive,
			ProcessedAt: &createdAt,
			IsFavourite: isFavourite,
		},
		Name:    name,
		SpaceID: spaceID,
		Type:    satType,
	}, nil
}

// SatelliteRepository defines the interface for Satellite operations.
type SatelliteRepository interface {
	FindBySpaceID(ctx context.Context, spaceID string) (Satellite, error)
	FindAll(ctx context.Context) ([]Satellite, error)
	Save(ctx context.Context, satellite Satellite) error
	Update(ctx context.Context, satellite Satellite) error
	DeleteBySpaceID(ctx context.Context, spaceID string) error
	SaveBatch(ctx context.Context, satellites []Satellite) error
	FindAllWithPagination(ctx context.Context, page int, pageSize int, searchRequest *SearchRequest) ([]Satellite, int64, error)
	FindSatelliteInfoWithPagination(ctx context.Context, page int, pageSize int, searchRequest *SearchRequest) ([]SatelliteInfo, int64, error)
	AssignSatelliteToContext(ctx context.Context, contextID, satelliteID string) error
	RemoveSatelliteFromContext(ctx context.Context, contextID, satelliteID string) error
	FindContextsBySatellite(ctx context.Context, satelliteID string) ([]GameContext, error)
	FindSatellitesByContext(ctx context.Context, contextID string) ([]Satellite, error)
}

type SatellitePosition struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Timestamp time.Time `json:"time"`
	CreatedAt time.Time `json:"created_at"`
}

type SatelliteInfo struct {
	Satellite Satellite // The associated Satellite
	TLEs      []TLE     // List of TLEs ordered by most recent
}

// NewSatelliteInfo creates a new SatelliteInfo instance.
func NewSatelliteInfo(satellite Satellite, tles []TLE) SatelliteInfo {
	sort.Slice(tles, func(i, j int) bool {
		return tles[i].Epoch.After(tles[j].Epoch)
	})

	return SatelliteInfo{
		Satellite: satellite,
		TLEs:      tles,
	}
}

// GetMostRecentTLE returns the most recent TLE from the SatelliteInfo.
func (info *SatelliteInfo) GetMostRecentTLE() *TLE {
	if len(info.TLEs) == 0 {
		return nil
	}
	return &info.TLEs[0]
}

// AddTLE adds a new TLE to the SatelliteInfo and keeps the list sorted by most recent.
func (info *SatelliteInfo) AddTLE(tle TLE) {
	info.TLEs = append(info.TLEs, tle)
	sort.Slice(info.TLEs, func(i, j int) bool {
		return info.TLEs[i].Epoch.After(info.TLEs[j].Epoch)
	})
}
