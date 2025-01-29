package models

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

// Satellite represents a satellite database model.
type Satellite struct {
	ModelBase
	Name           string     `gorm:"size:255;not null"`        // Satellite name
	NoradID        string     `gorm:"size:255;unique;not null"` // NORAD ID
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
}

// MapToDomain converts a Satellite database model to a Satellite domain model.
func MapToSatelliteDomain(s Satellite) domain.Satellite {
	return domain.Satellite{
		ModelBase: domain.ModelBase{
			ID:          s.ID,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   &s.UpdatedAt,
			DeleteAt:    s.DeleteAt,
			ProcessedAt: s.ProcessedAt,
			IsActive:    s.IsActive,
			IsFavourite: s.IsFavourite,
			DisplayName: s.DisplayName,
		},
		Name:           s.Name,
		NoradID:        s.NoradID,
		Type:           domain.SatelliteType(s.Type),
		LaunchDate:     s.LaunchDate,
		DecayDate:      s.DecayDate,
		IntlDesignator: s.IntlDesignator,
		Owner:          s.Owner,
		ObjectType:     s.ObjectType,
		Period:         s.Period,
		Inclination:    s.Inclination,
		Apogee:         s.Apogee,
		Perigee:        s.Perigee,
		RCS:            s.RCS,
		Altitude:       s.Altitude,
	}
}

// MapToSatelliteModel converts a Satellite domain model to a Satellite database model.
func MapToSatelliteModel(d domain.Satellite) Satellite {
	return Satellite{
		ModelBase: ModelBase{
			ID:          d.ModelBase.ID,
			CreatedAt:   d.ModelBase.CreatedAt,
			UpdatedAt:   *d.ModelBase.UpdatedAt,
			DeleteAt:    d.ModelBase.DeleteAt,
			ProcessedAt: d.ModelBase.ProcessedAt,
			IsActive:    d.ModelBase.IsActive,
			IsFavourite: d.ModelBase.IsFavourite,
			DisplayName: d.ModelBase.DisplayName,
		},
		Name:           d.Name,
		NoradID:        d.NoradID,
		Type:           string(d.Type),
		LaunchDate:     d.LaunchDate,
		DecayDate:      d.DecayDate,
		IntlDesignator: d.IntlDesignator,
		Owner:          d.Owner,
		ObjectType:     d.ObjectType,
		Period:         d.Period,
		Inclination:    d.Inclination,
		Apogee:         d.Apogee,
		Perigee:        d.Perigee,
		RCS:            d.RCS,
		Altitude:       d.Altitude,
	}
}
