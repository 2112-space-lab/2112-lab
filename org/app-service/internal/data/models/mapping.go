package models

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

// TileSatelliteMapping defines the relationship between a satellite and a tile.
type TileSatelliteMapping struct {
	ModelBase
	NoradID               string    `gorm:"size:255;not null;index"`                       // Foreign key to Satellite table via NORAD ID
	TileID                string    `gorm:"type:char(36);not null;uniqueIndex:norad_tile"` // Foreign key to Tile table
	IntersectionLatitude  float64   `gorm:"type:double precision;not null;"`               // Latitude of the intersection point
	IntersectionLongitude float64   `gorm:"type:double precision;not null;"`               // Longitude of the intersection point
	IntersectedAt         time.Time `gorm:"not null"`                                      // Time of intersection
	ComputationID         string    `gorm:"size:36;not null;index"`                        // Foreign key to Computation table
}

// MapToTileSatelliteMappingDomain converts a models.TileSatelliteMapping to a domain.TileSatelliteMapping.
func MapToTileSatelliteMappingDomain(t TileSatelliteMapping) (domain.TileSatelliteMapping, error) {
	return domain.TileSatelliteMapping{
		ModelBase: domain.ModelBase{
			ID:          t.ID,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   &t.UpdatedAt,
			DeleteAt:    t.DeleteAt,
			ProcessedAt: t.ProcessedAt,
			IsActive:    t.IsActive,
			IsFavourite: t.IsFavourite,
			DisplayName: t.DisplayName,
		},
		NoradID:               t.NoradID,
		TileID:                t.TileID,
		IntersectionLatitude:  t.IntersectionLatitude,
		IntersectionLongitude: t.IntersectionLongitude,
		IntersectedAt:         t.IntersectedAt,
		ComputationID:         t.ComputationID,
	}, nil
}
