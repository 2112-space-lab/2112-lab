package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
)

type TileSatelliteMappingRepository struct {
	db *data.Database
}

// NewTileSatelliteMappingRepository creates a new instance of TileSatelliteMappingRepository.
func NewTileSatelliteMappingRepository(db *data.Database) TileSatelliteMappingRepository {
	return TileSatelliteMappingRepository{db: db}
}

func (r *TileSatelliteMappingRepository) FindByNoradIDAndTile(ctx context.Context, contextID, noradID, tileID string) ([]domain.TileSatelliteMapping, error) {
	var mappings []domain.TileSatelliteMapping
	result := r.db.DbHandler.WithContext(ctx).
		Where("context_id = ? AND norad_id = ? AND tile_id = ?", contextID, noradID, tileID).
		Find(&mappings)
	return mappings, result.Error
}

func (r *TileSatelliteMappingRepository) FindAll(ctx context.Context, contextID string) ([]domain.TileSatelliteMapping, error) {
	var mappings []domain.TileSatelliteMapping
	result := r.db.DbHandler.WithContext(ctx).
		Where("context_id = ?", contextID).
		Find(&mappings)
	return mappings, result.Error
}

func (r *TileSatelliteMappingRepository) Save(ctx context.Context, mapping domain.TileSatelliteMapping) error {
	if mapping.ID == "" {
		mapping.ID = uuid.NewString()
	}
	return r.db.DbHandler.WithContext(ctx).Create(&mapping).Error
}

func (r *TileSatelliteMappingRepository) Update(ctx context.Context, mapping domain.TileSatelliteMapping) error {
	return r.db.DbHandler.WithContext(ctx).Save(&mapping).Error
}

func (r *TileSatelliteMappingRepository) Delete(ctx context.Context, id string) error {
	return r.db.DbHandler.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.TileSatelliteMapping{}).Error
}

func (r *TileSatelliteMappingRepository) SaveBatch(ctx context.Context, mappings []domain.TileSatelliteMapping) error {
	if len(mappings) == 0 {
		return nil
	}
	for i := range mappings {
		if mappings[i].ID == "" {
			mappings[i].ID = uuid.NewString()
		}
	}
	return r.db.DbHandler.WithContext(ctx).Create(&mappings).Error
}

func (r *TileSatelliteMappingRepository) FindSatellitesForTiles(ctx context.Context, contextID string, tileIDs []string) ([]domain.Satellite, error) {
	var satellites []models.Satellite
	err := r.db.DbHandler.WithContext(ctx).
		Table("tile_satellite_mappings").
		Select("satellites.*").
		Joins("JOIN satellites ON tile_satellite_mappings.norad_id = satellites.norad_id").
		Where("tile_satellite_mappings.context_id = ? AND tile_satellite_mappings.tile_id IN ?", contextID, tileIDs).
		Find(&satellites).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find satellites for tiles: %w", err)
	}
	var domainSatellites []domain.Satellite
	for _, satellite := range satellites {
		domainSatellites = append(domainSatellites, models.MapToSatelliteDomain(satellite))
	}
	return domainSatellites, nil
}

func (r *TileSatelliteMappingRepository) FindAllVisibleTilesByNoradIDSortedByAOSTime(ctx context.Context, contextID, noradID string) ([]domain.TileSatelliteInfo, error) {
	var mappings []domain.TileSatelliteMapping
	result := r.db.DbHandler.WithContext(ctx).
		Where("context_id = ? AND norad_id = ?", contextID, noradID).
		Order("aos ASC").
		Find(&mappings)
	if result.Error != nil {
		return nil, result.Error
	}

	tileIDs := make([]string, len(mappings))
	for i, mapping := range mappings {
		tileIDs[i] = mapping.TileID
	}

	var tiles []models.Tile
	err := r.db.DbHandler.WithContext(ctx).Where("id IN ?", tileIDs).Find(&tiles).Error
	if err != nil {
		return nil, err
	}

	tileMap := make(map[string]models.Tile)
	for _, tile := range tiles {
		tileMap[tile.ID] = tile
	}

	var infos []domain.TileSatelliteInfo
	for _, mapping := range mappings {
		tile := tileMap[mapping.TileID]
		infos = append(infos, domain.TileSatelliteInfo{
			MappingID:     mapping.ID,
			TileID:        tile.ID,
			TileQuadkey:   tile.Quadkey,
			TileCenterLat: tile.CenterLat,
			TileCenterLon: tile.CenterLon,
			TileZoomLevel: tile.ZoomLevel,
			NoradID:       mapping.NoradID,
		})
	}
	return infos, nil
}

func (r *TileSatelliteMappingRepository) ListSatellitesMappingWithPagination(ctx context.Context, contextID string, page, pageSize int, search *domain.SearchRequest) ([]domain.TileSatelliteInfo, int64, error) {
	var (
		mappings           []domain.TileSatelliteMapping
		totalRecords       int64
		tileSatelliteInfos []domain.TileSatelliteInfo
	)

	offset := (page - 1) * pageSize

	query := r.db.DbHandler.WithContext(ctx).Table("tile_satellite_mappings").
		Where("context_id = ?", contextID)

	if search != nil && search.Wildcard != "" {
		wildcard := "%" + search.Wildcard + "%"
		query = query.Where("norad_id LIKE ? OR tile_id LIKE ?", wildcard, wildcard)
	}

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(pageSize).Offset(offset).Find(&mappings).Error; err != nil {
		return nil, 0, err
	}

	tileIDs := make([]string, len(mappings))
	for i, mapping := range mappings {
		tileIDs[i] = mapping.TileID
	}

	var tiles []models.Tile
	err := r.db.DbHandler.WithContext(ctx).Where("id IN ?", tileIDs).Find(&tiles).Error
	if err != nil {
		return nil, 0, err
	}

	tileMap := make(map[string]models.Tile)
	for _, tile := range tiles {
		tileMap[tile.ID] = tile
	}

	for _, mapping := range mappings {
		tile := tileMap[mapping.TileID]
		tileSatelliteInfos = append(tileSatelliteInfos, domain.TileSatelliteInfo{
			MappingID:     mapping.ID,
			TileID:        tile.ID,
			TileQuadkey:   tile.Quadkey,
			TileCenterLat: tile.CenterLat,
			TileCenterLon: tile.CenterLon,
			TileZoomLevel: tile.ZoomLevel,
			NoradID:       mapping.NoradID,
		})
	}

	return tileSatelliteInfos, totalRecords, nil
}

func (r *TileSatelliteMappingRepository) GetSatelliteMappingsByNoradID(ctx context.Context, contextID, noradID string) ([]domain.TileSatelliteInfo, error) {
	var mappings []domain.TileSatelliteMapping
	err := r.db.DbHandler.WithContext(ctx).
		Where("context_id = ? AND norad_id = ?", contextID, noradID).
		Find(&mappings).Error
	if err != nil {
		return nil, err
	}

	tileIDs := make([]string, len(mappings))
	for i, mapping := range mappings {
		tileIDs[i] = mapping.TileID
	}

	var tiles []models.Tile
	err = r.db.DbHandler.WithContext(ctx).Where("id IN ?", tileIDs).Find(&tiles).Error
	if err != nil {
		return nil, err
	}

	tileMap := make(map[string]models.Tile)
	for _, tile := range tiles {
		tileMap[tile.ID] = tile
	}

	var infos []domain.TileSatelliteInfo
	for _, mapping := range mappings {
		tile := tileMap[mapping.TileID]
		infos = append(infos, domain.TileSatelliteInfo{
			MappingID:     mapping.ID,
			TileID:        tile.ID,
			TileQuadkey:   tile.Quadkey,
			TileCenterLat: tile.CenterLat,
			TileCenterLon: tile.CenterLon,
			TileZoomLevel: tile.ZoomLevel,
			NoradID:       mapping.NoradID,
		})
	}
	return infos, nil
}

func (r *TileSatelliteMappingRepository) DeleteMappingsByNoradID(ctx context.Context, contextID, noradID string) error {
	return r.db.DbHandler.WithContext(ctx).
		Where("context_id = ? AND norad_id = ?", contextID, noradID).
		Delete(&domain.TileSatelliteMapping{}).Error
}
