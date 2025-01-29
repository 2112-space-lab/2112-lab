package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xpolygon"
)

// TileRepository defines the interface for Tile repository operations.
type TileRepository interface {
	// Existing methods
	FindByQuadkey(ctx context.Context, key string) (*Tile, error)                                                            // Find a tile by Quadkey
	FindBySpatialLocation(ctx context.Context, lat, lon float64) (*Tile, error)                                              // Find a tile by spatial location
	FindTilesInRegion(ctx context.Context, contextID string, minLat, minLon, maxLat, maxLon float64) ([]Tile, error)         // Find tiles intersecting a region
	FindAll(ctx context.Context) ([]Tile, error)                                                                             // Retrieve all tiles
	Save(ctx context.Context, tile Tile) error                                                                               // Save a new tile
	Update(ctx context.Context, tile Tile) error                                                                             // Update an existing tile
	Upsert(ctx context.Context, tile Tile) error                                                                             // Upsert (insert or update) a tile
	DeleteByQuadkey(ctx context.Context, key string) error                                                                   // Delete a tile by Quadkey
	DeleteBySpatialLocation(ctx context.Context, lat float64, lon float64) error                                             // Delete a tile by spatial location
	FindTilesVisibleFromLine(ctx context.Context, sat Satellite, points []SatellitePosition) ([]TileSatelliteMapping, error) // Find tiles visible from a satellite trajectory
	FindTilesIntersectingLocation(ctx context.Context, contextID string, lat, lon, radius float64) ([]Tile, error)           // Find tiles intersecting a location with a radius

	// New methods for context support
	AssociateTileWithContext(ctx context.Context, contextID string, tileID string) error // Associate a tile with a context
	GetTilesByContext(ctx context.Context, contextID string) ([]Tile, error)             // Retrieve all tiles associated with a context
	RemoveTileFromContext(ctx context.Context, contextID string, tileID string) error    // Remove a tile from a context
}

// Tile represents the domain entity Tile
type Tile struct {
	ModelBase
	Quadkey   string           // Quadkey representing the tile
	ZoomLevel int              // Zoom level of the tile
	CenterLat float64          // Center latitude of the tile
	CenterLon float64          // Center longitude of the tile
	NbFaces   int              // Number of faces in the tile's geometry
	Radius    float64          // Radius of the tile (in meters or other unit)
	Vertices  []xpolygon.Point // Vertices representing the boundary of the tile
}

// NewTile constructor
// NewTile constructor
func NewTile(polygon xpolygon.Polygon, createdAt time.Time, isFavourite bool, isActive bool, displayName string) Tile {
	return Tile{
		ModelBase: ModelBase{
			ID:          uuid.NewString(),
			CreatedAt:   createdAt,
			UpdatedAt:   &createdAt,
			DisplayName: displayName,
			IsActive:    isActive,
			ProcessedAt: &createdAt,
			IsFavourite: isFavourite,
		},
		Quadkey:   polygon.Center.Key(),     // Extract quadkey from the polygon center
		ZoomLevel: polygon.Center.Level,     // Use the zoom level from the center
		CenterLat: polygon.Center.Latitude,  // Use center latitude
		CenterLon: polygon.Center.Longitude, // Use center longitude
		NbFaces:   polygon.NbFaces,          // Number of faces in the tile
		Radius:    polygon.Radius,           // Tile radius
		Vertices:  polygon.Boundaries,       // Boundary vertices
	}
}

type TileVisibility struct {
	Tile
	AosTime time.Time
}
