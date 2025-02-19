package models

import (
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
)

// Satellite represents a satellite database model.
type Satellite struct {
	ModelBase
	Name           string     `gorm:"size:255;not null"`        // Satellite name
	SpaceID        string     `gorm:"size:255;unique;not null"` // SPACE ID
	Type           string     `gorm:"size:255"`                 // Satellite type (e.g., telescope, communication)
	LaunchDate     *time.Time `gorm:"type:date"`                // Launch date
	DecayDate      *time.Time `gorm:"type:date"`                // Decay date (optional)
	IntlDesignator string     `gorm:"size:255"`                 // International designator
	Owner          string     `gorm:"size:255"`                 // Ownership information
	ObjectType     string     `gorm:"size:255"`                 // Object type (e.g., "PAYLOAD")
	Period         *float64   `gorm:"type:float"`               // Orbital period in minutes (optional)
	Inclination    *float64   `gorm:"type:float"`               // Orbital inclination in degrees (optional)
	Apogee         *float64   `gorm:"type:float"`               // Apogee altitude in kilometers (optional)
	Perigee        *float64   `gorm:"type:float"`               // Perigee altitude in kilometers (optional)
	RCS            *float64   `gorm:"type:float"`               // Radar cross-section in square meters (optional)
	Altitude       *float64   `gorm:"type:float"`               // Altitude in kilometers (optional)
	OrbitType      string     `gorm:"size:255;not null"`
}

// MapToSatelliteDomain converts a Satellite database model to a Satellite domain model.
func MapToSatelliteDomain(s Satellite) domain.Satellite {
	domainSatellite, err := domain.NewSatelliteFromParameters(
		s.Name,
		s.SpaceID,
		domain.SatelliteType(s.Type),
		s.LaunchDate,
		s.DecayDate,
		s.IntlDesignator,
		s.Owner,
		s.ObjectType,
		s.Period,
		s.Inclination,
		s.Apogee,
		s.Perigee,
		s.RCS,
		s.Altitude,
	)

	if err != nil {
		return domain.Satellite{}
	}

	return domainSatellite
}

// MapToSatelliteModel converts a Satellite domain model to a Satellite database model.
func MapToSatelliteModel(d domain.Satellite) Satellite {
	return Satellite{
		ModelBase: ModelBase{
			ID:          d.ModelBase.ID,
			CreatedAt:   d.ModelBase.CreatedAt,
			UpdatedAt:   *d.ModelBase.UpdatedAt, // Ensure it's not nil
			DeleteAt:    d.ModelBase.DeleteAt,
			ProcessedAt: d.ModelBase.ProcessedAt,
			IsActive:    d.ModelBase.IsActive,
			IsFavourite: d.ModelBase.IsFavourite,
			DisplayName: d.ModelBase.DisplayName,
		},
		Name:           d.Name,
		SpaceID:        d.SpaceID,
		Type:           string(d.Type),
		LaunchDate:     xtime.ConvertToTimePtr(d.LaunchDate),
		DecayDate:      xtime.ConvertToTimePtr(d.DecayDate),
		IntlDesignator: d.IntlDesignator,
		Owner:          d.Owner,
		ObjectType:     d.ObjectType,
		Period:         fx.ConvertToFloatPtr(d.PeriodInMinutes),
		Inclination:    fx.ConvertToFloatPtr(d.InclinationInDegrees),
		Apogee:         fx.ConvertToFloatPtr(d.ApogeeInKm),
		Perigee:        fx.ConvertToFloatPtr(d.PerigeeInKm),
		RCS:            fx.ConvertToFloatPtr(d.RCS),
		Altitude:       fx.ConvertToFloatPtr(d.Altitude),
		OrbitType:      string(d.OrbitType),
	}
}
