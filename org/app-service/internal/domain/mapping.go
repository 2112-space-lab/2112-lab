package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MappingRepository interface {
	FindByNoradIDAndTile(ctx context.Context, contextID, noradID, tileID string) ([]TileSatelliteMapping, error)
	FindAll(ctx context.Context, contextID string) ([]TileSatelliteMapping, error) // Updated to include contextID
	Save(ctx context.Context, visibility TileSatelliteMapping) error
	Update(ctx context.Context, visibility TileSatelliteMapping) error
	Delete(ctx context.Context, id string) error
	SaveBatch(ctx context.Context, visibilities []TileSatelliteMapping) error
	FindSatellitesForTiles(ctx context.Context, contextID string, tileIDs []string) ([]Satellite, error)
	FindAllVisibleTilesByNoradIDSortedByAOSTime(ctx context.Context, contextID, noradID string) ([]TileSatelliteInfo, error)
	ListSatellitesMappingWithPagination(ctx context.Context, contextID string, page int, pageSize int, search *SearchRequest) ([]TileSatelliteInfo, int64, error)
	GetSatelliteMappingsByNoradID(ctx context.Context, contextID, noradID string) ([]TileSatelliteInfo, error)
	DeleteMappingsByNoradID(ctx context.Context, contextID, noradID string) error
}

// TileSatelliteMapping represents the domain entity TileSatelliteMapping
type TileSatelliteMapping struct {
	ModelBase
	NoradID               string
	TileID                string
	IntersectionLongitude float64
	IntersectionLatitude  float64
	IntersectedAt         time.Time
	ComputationID         string
}

// NewMapping constructor
func NewMapping(noradID string,
	tileID string, intersection Point, interestedTime time.Time, createdAt time.Time, displayName string, isActive bool, isFavourite bool) TileSatelliteMapping {

	return TileSatelliteMapping{
		ModelBase: ModelBase{
			ID:          uuid.NewString(),
			CreatedAt:   createdAt,
			UpdatedAt:   &createdAt,
			DisplayName: displayName,
			IsActive:    isActive,
			IsFavourite: isFavourite,
			ProcessedAt: &createdAt,
		},
		NoradID:               noradID,
		TileID:                tileID,
		IntersectionLongitude: intersection.Longitude,
		IntersectionLatitude:  intersection.Latitude,
		IntersectedAt:         interestedTime,
	}

}

// TileSatelliteInfo represents the aggregated data of a tile and satellite, sorted by AOS time.
type TileSatelliteInfo struct {
	MappingID     string
	TileID        string  // The ID of the tile
	TileQuadkey   string  // The Quadkey of the tile
	TileCenterLat float64 // Latitude of the tile center
	TileCenterLon float64 // Longitude of the tile center
	TileZoomLevel int     // Zoom level of the tile
	NoradID       string  // The NORAD ID of the satellite
	Intersection  Point
}

type Point struct {
	Longitude float64
	Latitude  float64
}
