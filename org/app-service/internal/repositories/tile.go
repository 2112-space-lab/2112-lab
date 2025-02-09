package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"gorm.io/gorm"
)

type TileRepository struct {
	db *data.Database
}

// NewTileRepository creates a new instance of TileRepository.
func NewTileRepository(db *data.Database) domain.TileRepository {
	return &TileRepository{db: db}
}

// FindTilesInRegion retrieves tiles that intersect a given bounding box and belong to a specific context.
func (r *TileRepository) FindTilesInRegion(ctx context.Context, contextID string, minLat, minLon, maxLat, maxLon float64) ([]domain.Tile, error) {
	var tiles []models.Tile

	// Execute the query with context filtering
	result := r.db.DbHandler.WithContext(ctx).Raw(`
		SELECT t.*
		FROM tiles t
		INNER JOIN context_tiles ct ON t.id = ct.tile_id
		WHERE ct.context_id = ?
		AND ST_Intersects(
			t.spatial_index,
			ST_MakeEnvelope(?, ?, ?, ?, 4326)
		)
	`, contextID, minLon, minLat, maxLon, maxLat).Scan(&tiles)

	if result.Error != nil {
		if errors.Is(result.Error, context.Canceled) {
			return nil, fmt.Errorf("query canceled: %w", result.Error)
		}
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, fmt.Errorf("query deadline exceeded: %w", result.Error)
		}
		return nil, fmt.Errorf("failed to find tiles in region for context %s: %w", contextID, result.Error)
	}

	// Map tiles to domain models
	var domainTiles []domain.Tile
	for _, tile := range tiles {
		domainTiles = append(domainTiles, models.MapToTileDomain(tile))
	}

	return domainTiles, nil
}

// FindTilesIntersectingLocation retrieves all tiles that intersect the given location and belong to a specific context.
func (r *TileRepository) FindTilesIntersectingLocation(ctx context.Context, contextID string, lat, lon, radius float64) ([]domain.Tile, error) {
	var tiles []models.Tile

	// Query tiles that intersect the user's location and are associated with the given context
	result := r.db.DbHandler.WithContext(ctx).Raw(`
		SELECT t.*
		FROM tiles t
		INNER JOIN context_tiles ct ON t.id = ct.tile_id
		WHERE ct.context_id = ?
		AND ST_DWithin(
			t.spatial_index,
			ST_MakePoint(?, ?)::geography,
			?
		)
	`, contextID, lon, lat, radius).Scan(&tiles)

	if result.Error != nil {
		if errors.Is(result.Error, context.Canceled) {
			return nil, fmt.Errorf("query canceled: %w", result.Error)
		}
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, fmt.Errorf("query deadline exceeded: %w", result.Error)
		}
		return nil, fmt.Errorf("failed to find tiles intersecting location for context %s: %w", contextID, result.Error)
	}

	// Map tiles to domain models
	var domainTiles []domain.Tile
	for _, tile := range tiles {
		domainTiles = append(domainTiles, models.MapToTileDomain(tile))
	}

	return domainTiles, nil
}

// FindByQuadkey retrieves a Tile by its quadkey.
func (r *TileRepository) FindByQuadkey(ctx context.Context, quadkey string) (*domain.Tile, error) {
	var tile models.Tile
	result := r.db.DbHandler.Where("quadkey = ?", quadkey).First(&tile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	tileMapped := models.MapToTileDomain(tile)
	return &tileMapped, nil
}

// FindBySpatialLocation retrieves a Tile by a geographical location using spatial indexing.
func (r *TileRepository) FindBySpatialLocation(ctx context.Context, lat, lon float64) (*domain.Tile, error) {
	var tile models.Tile
	result := r.db.DbHandler.Raw(`
		SELECT *
		FROM tiles
		WHERE ST_Contains(spatial_index, ST_SetSRID(ST_Point(?, ?), 4326))
		LIMIT 1
	`, lon, lat).Scan(&tile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	tileMapped := models.MapToTileDomain(tile)
	return &tileMapped, nil
}

// FindAll retrieves all Tiles.
func (r *TileRepository) FindAll(ctx context.Context) ([]domain.Tile, error) {
	var tiles []models.Tile
	result := r.db.DbHandler.Find(&tiles)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainTiles []domain.Tile
	for _, t := range tiles {
		domainTiles = append(domainTiles, models.MapToTileDomain(t))
	}
	return domainTiles, nil
}

// Save creates a new Tile record.
func (r *TileRepository) Save(ctx context.Context, tile domain.Tile) error {
	modelTile := models.MapFromDomain(tile)
	return r.db.DbHandler.Create(&modelTile).Error
}

// SaveBatch allows batch insertion of tiles for optimized performance.
func (r *TileRepository) SaveBatch(ctx context.Context, tiles []domain.Tile) error {
	modelTiles := make([]models.Tile, len(tiles))
	for i, t := range tiles {
		modelTiles[i] = models.MapFromDomain(t)
	}
	return r.db.DbHandler.Create(&modelTiles).Error
}

// Update modifies an existing Tile record.
func (r *TileRepository) Update(ctx context.Context, tile domain.Tile) error {
	modelTile := models.MapFromDomain(tile)
	return r.db.DbHandler.Save(&modelTile).Error
}

// DeleteByQuadkey removes a Tile record by its quadkey.
func (r *TileRepository) DeleteByQuadkey(ctx context.Context, key string) error {
	return r.db.DbHandler.Where("quadkey = ?", key).Delete(&models.Tile{}).Error
}

// AssociateTileWithContext associates a Tile with a specific Context.
func (r *TileRepository) AssociateTileWithContext(ctx context.Context, contextID string, tileID string) error {
	contextTile := models.ContextTile{
		ContextID: contextID,
		TileID:    tileID,
	}

	if err := r.db.DbHandler.Create(&contextTile).Error; err != nil {
		return fmt.Errorf("failed to associate Tile with context: %w", err)
	}
	return nil
}

// GetTilesByContext retrieves all Tiles associated with a specific Context.
func (r *TileRepository) GetTilesByContext(ctx context.Context, contextID string) ([]domain.Tile, error) {
	var contextTiles []models.ContextTile

	if err := r.db.DbHandler.Where("context_id = ?", contextID).Find(&contextTiles).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve Tiles by context: %w", err)
	}

	tileIDs := make([]string, len(contextTiles))
	for i, contextTile := range contextTiles {
		tileIDs[i] = contextTile.TileID
	}

	var tiles []models.Tile
	if err := r.db.DbHandler.Where("id IN ?", tileIDs).Find(&tiles).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve Tile details: %w", err)
	}

	var domainTiles []domain.Tile
	for _, tile := range tiles {
		domainTiles = append(domainTiles, models.MapToTileDomain(tile))
	}

	return domainTiles, nil
}

// RemoveTileFromContext removes the association between a Tile and a Context.
func (r *TileRepository) RemoveTileFromContext(ctx context.Context, contextID string, tileID string) error {
	if err := r.db.DbHandler.Where("context_id = ? AND tile_id = ?", contextID, tileID).
		Delete(&models.ContextTile{}).Error; err != nil {
		return fmt.Errorf("failed to remove Tile from context: %w", err)
	}
	return nil
}

// FindTilesVisibleFromLine retrieves Tiles intersecting a satellite's trajectory.
func (r *TileRepository) FindTilesVisibleFromLine(ctx context.Context, sat domain.Satellite, points []domain.SatellitePosition) ([]domain.TileSatelliteMapping, error) {
	if len(points) < 2 {
		return nil, fmt.Errorf("at least two points are required to create a line")
	}

	wktPoints := make([]string, len(points))
	for i, point := range points {
		wktPoints[i] = fmt.Sprintf("%f %f", point.Longitude, point.Latitude)
	}
	lineString := fmt.Sprintf("LINESTRING(%s)", strings.Join(wktPoints, ", "))

	query := `
        WITH line_geom AS (
            SELECT ST_GeomFromText(?, 4326) AS geom
        )
        SELECT 
            tiles.*,
            ST_AsText(ST_PointOnSurface(ST_Intersection(line_geom.geom, spatial_index))) AS intersection_geom
        FROM tiles, line_geom
        WHERE ST_Intersects(spatial_index, line_geom.geom)
    `

	var results []struct {
		models.Tile
		IntersectionGeom string `gorm:"column:intersection_geom"`
	}
	result := r.db.DbHandler.Raw(query, lineString).Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}

	var mappings []domain.TileSatelliteMapping
	for _, res := range results {
		tile := models.MapToTileDomain(res.Tile)
		interestPoint, err := parseIntersectionGeometry(res.IntersectionGeom)
		if err != nil {
			log.Printf("Failed to parse intersection geometry for TileID %s: %v\n", tile.ID, err)
			continue
		}

		nowUtc := time.Now().UTC()
		mapping := domain.NewMapping(
			sat.NoradID,
			tile.ID,
			interestPoint,
			nowUtc,
			nowUtc,
			"",
			true,
			false,
		)
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

// parseIntersectionGeometry parses WKT intersection points.
func parseIntersectionGeometry(wkt string) (domain.Point, error) {
	if wkt == "" || !strings.HasPrefix(wkt, "POINT(") || !strings.HasSuffix(wkt, ")") {
		return domain.Point{}, fmt.Errorf("invalid WKT format: %s", wkt)
	}

	coordinates := strings.TrimPrefix(strings.TrimSuffix(wkt, ")"), "POINT(")
	parts := strings.Fields(coordinates)
	if len(parts) != 2 {
		return domain.Point{}, fmt.Errorf("invalid WKT coordinates: %s", coordinates)
	}

	longitude, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return domain.Point{}, fmt.Errorf("error parsing longitude: %w", err)
	}

	latitude, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return domain.Point{}, fmt.Errorf("error parsing latitude: %w", err)
	}

	return domain.Point{
		Longitude: longitude,
		Latitude:  latitude,
	}, nil
}

// DeleteBySpatialLocation removes a Tile record by its geographical location.
func (r *TileRepository) DeleteBySpatialLocation(ctx context.Context, lat, lon float64) error {
	var tile models.Tile
	result := r.db.DbHandler.Raw(`
		SELECT *
		FROM tiles
		WHERE ST_Contains(spatial_index, ST_SetSRID(ST_Point(?, ?), 4326))
		LIMIT 1
	`, lon, lat).Scan(&tile)

	if result.Error != nil {
		if result.RowsAffected == 0 {
			return fmt.Errorf("tile not found at the specified location")
		}
		return result.Error
	}

	// Delete the tile
	return r.db.DbHandler.Delete(&tile).Error
}

// Upsert inserts or updates a Tile record in the database.
func (r *TileRepository) Upsert(ctx context.Context, tile domain.Tile) error {
	// Check for an existing tile using the spatial location or quadkey
	existingTile, err := r.FindByQuadkey(ctx, tile.Quadkey)
	if err != nil {
		return err
	}

	if existingTile != nil {
		// Update if the tile exists
		return r.Update(ctx, tile)
	}

	// Save if the tile doesn't exist
	return r.Save(ctx, tile)
}
