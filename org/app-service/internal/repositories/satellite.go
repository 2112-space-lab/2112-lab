package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/redis"
	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SatelliteTLEAggregate is used for joining Satellite and TLE data.
type SatelliteTLEAggregate struct {
	// Satellite fields
	ID             string     `gorm:"column:id"`
	Name           string     `gorm:"column:name"`
	SpaceID        string     `gorm:"column:space_id"`
	Owner          string     `gorm:"column:owner"`
	LaunchDate     *time.Time `gorm:"column:launch_date"`
	DecayDate      *time.Time `gorm:"column:decay_date"`
	IntlDesignator string     `gorm:"column:international_designator"`
	ObjectType     string     `gorm:"column:object_type"`
	Period         *float64   `gorm:"column:period"`
	Inclination    *float64   `gorm:"column:inclination"`
	Apogee         *float64   `gorm:"column:apogee"`
	Perigee        *float64   `gorm:"column:perigee"`
	RCS            *float64   `gorm:"column:rcs"`
	Altitude       *float64   `gorm:"column:altitude"`
	IsActive       bool       `gorm:"column:is_active"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at"`
	ProcessedAt    *time.Time `gorm:"column:processed_at"`
	IsFavourite    bool       `gorm:"column:is_favourite"`

	// TLE fields
	Line1        *string    `gorm:"column:line1"`
	Line2        *string    `gorm:"column:line2"`
	TLEUpdatedAt *time.Time `gorm:"column:tle_updated_at"`
}

// SatelliteRepository manages satellite data access.
type SatelliteRepository struct {
	db      *data.Database
	redis   *redis.RedisClient
	lockTTL time.Duration
}

// NewSatelliteRepository creates a new SatelliteRepository instance.
func NewSatelliteRepository(db *data.Database, redis *redis.RedisClient, lockTTL time.Duration) SatelliteRepository {
	return SatelliteRepository{db: db, redis: redis, lockTTL: lockTTL}
}

// FindBySpaceID retrieves a satellite by its SPACE ID, excluding deleted ones.
func (r *SatelliteRepository) FindBySpaceID(ctx context.Context, spaceID string) (domain.Satellite, error) {
	var satellite models.Satellite
	result := r.db.DbHandler.Where("space_id = ? AND deleted_at IS NULL", spaceID).First(&satellite)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.Satellite{}, nil
	}
	return models.MapToSatelliteDomain(satellite), result.Error
}

// FindAll retrieves all satellites excluding deleted ones.
func (r *SatelliteRepository) FindAll(ctx context.Context) ([]domain.Satellite, error) {
	var satellites []models.Satellite
	result := r.db.DbHandler.Where("deleted_at IS NULL").Find(&satellites)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainSatellites []domain.Satellite
	for _, satellite := range satellites {
		domainSatellites = append(domainSatellites, models.MapToSatelliteDomain(satellite))
	}
	return domainSatellites, nil
}

// Save creates a new satellite record.
func (r *SatelliteRepository) Save(ctx context.Context, satellite domain.Satellite) error {
	model := models.MapToSatelliteModel(satellite)
	return r.db.DbHandler.Create(&model).Error
}

// Update modifies an existing satellite record.
func (r *SatelliteRepository) Update(ctx context.Context, satellite domain.Satellite) error {
	model := models.MapToSatelliteModel(satellite)
	return r.db.DbHandler.Save(&model).Error
}

// DeleteBySpaceID marks a satellite record as deleted.
func (r *SatelliteRepository) DeleteBySpaceID(ctx context.Context, spaceID string) error {
	return r.db.DbHandler.Model(&models.Satellite{}).
		Where("space_id = ?", spaceID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// SaveBatch performs a batch insert or update (upsert) for satellites.
func (r *SatelliteRepository) SaveBatch(ctx context.Context, satellites []domain.Satellite) error {
	if len(satellites) == 0 {
		return nil
	}

	var modelsBatch []models.Satellite
	for _, satellite := range satellites {
		modelsBatch = append(modelsBatch, models.MapToSatelliteModel(satellite))
	}

	return r.db.DbHandler.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "space_id"}},
			UpdateAll: true,
		}).
		CreateInBatches(modelsBatch, 100).Error
}

// FindSatelliteInfoWithPagination retrieves satellites and their TLEs with pagination.
func (r *SatelliteRepository) FindSatelliteInfoWithPagination(ctx context.Context, page, pageSize int, searchRequest *domain.SearchRequest) ([]domain.SatelliteInfo, int64, error) {
	var results []SatelliteTLEAggregate
	var totalRecords int64

	// Calculate the offset for pagination
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// Build the query to retrieve satellites and their most recent TLEs
	query := r.db.DbHandler.Table("satellites").
		Select(`
			satellites.id, satellites.name, satellites.space_id, satellites.owner,
			satellites.launch_date, satellites.decay_date, satellites.international_designator,
			satellites.object_type, satellites.period, satellites.inclination, satellites.apogee,
			satellites.perigee, satellites.rcs, satellites.altitude, satellites.is_active,
			satellites.created_at, satellites.updated_at, satellites.processed_at, satellites.is_favourite,
			latest_tles.line1, latest_tles.line2, latest_tles.updated_at AS tle_updated_at
		`).
		Joins(`LEFT JOIN (
			SELECT t1.space_id, t1.line1, t1.line2, t1.updated_at
			FROM tles t1
			WHERE t1.updated_at = (
				SELECT MAX(t2.updated_at)
				FROM tles t2
				WHERE t2.space_id = t1.space_id
			)
		) AS latest_tles ON satellites.space_id = latest_tles.space_id`).
		Where("satellites.deleted_at IS NULL")

	// Apply search filters if a wildcard search is provided
	if searchRequest != nil && searchRequest.Wildcard != "" {
		wildcard := "%" + searchRequest.Wildcard + "%"
		query = query.Where("LOWER(satellites.name) LIKE LOWER(?) OR LOWER(satellites.space_id) LIKE LOWER(?)", wildcard, wildcard)
	}

	// Count total records
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Retrieve paginated results
	if err := query.Limit(pageSize).Offset(offset).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	// Map results to domain objects
	var satelliteInfos []domain.SatelliteInfo
	for _, result := range results {
		// Map satellite fields from the aggregate struct
		satellite := domain.Satellite{
			Name:           result.Name,
			SpaceID:        result.SpaceID,
			Owner:          result.Owner,
			LaunchDate:     result.LaunchDate,
			DecayDate:      result.DecayDate,
			IntlDesignator: result.IntlDesignator,
			ObjectType:     result.ObjectType,
			Period:         result.Period,
			Inclination:    result.Inclination,
			Apogee:         result.Apogee,
			Perigee:        result.Perigee,
			RCS:            result.RCS,
			Altitude:       result.Altitude,
			ModelBase: domain.ModelBase{
				ID:          result.ID,
				IsActive:    result.IsActive,
				CreatedAt:   result.CreatedAt,
				UpdatedAt:   result.UpdatedAt,
				ProcessedAt: result.ProcessedAt,
				IsFavourite: result.IsFavourite,
			},
		}

		// Map TLE data if available
		var tles []domain.TLE
		if result.Line1 != nil && result.Line2 != nil && result.TLEUpdatedAt != nil {
			tles = append(tles, domain.TLE{
				Line1: *result.Line1,
				Line2: *result.Line2,
				Epoch: *result.TLEUpdatedAt,
			})
		}

		// Create SatelliteInfo and append to result list
		satelliteInfos = append(satelliteInfos, domain.NewSatelliteInfo(satellite, tles))
	}

	return satelliteInfos, totalRecords, nil
}

// AssignSatelliteToContext associates a satellite with a context.
func (r *SatelliteRepository) AssignSatelliteToContext(ctx context.Context, contextID, satelliteID string) error {
	association := models.ContextSatellite{
		ContextID:   contextID,
		SatelliteID: satelliteID,
	}
	return r.db.DbHandler.Create(&association).Error
}

// RemoveSatelliteFromContext removes the association between a satellite and a context.
func (r *SatelliteRepository) RemoveSatelliteFromContext(ctx context.Context, contextID, satelliteID string) error {
	return r.db.DbHandler.Where("context_id = ? AND satellite_id = ?", contextID, satelliteID).
		Delete(&models.ContextSatellite{}).Error
}

// FindContextsBySatellite retrieves contexts associated with a given satellite.
func (r *SatelliteRepository) FindContextsBySatellite(ctx context.Context, satelliteID string) ([]domain.GameContext, error) {
	var contexts []models.Context
	result := r.db.DbHandler.Table("contexts").
		Joins("JOIN context_satellites ON contexts.id = context_satellites.context_id").
		Where("context_satellites.satellite_id = ?", satelliteID).
		Find(&contexts)

	if result.Error != nil {
		return nil, result.Error
	}

	var domainContexts []domain.GameContext
	for _, contextModel := range contexts {
		domainContexts = append(domainContexts, models.MapToContextDomain(contextModel))
	}
	return domainContexts, nil
}

// FindSatellitesByContext retrieves satellites associated with a given context.
func (r *SatelliteRepository) FindSatellitesByContext(ctx context.Context, contextID string) ([]domain.Satellite, error) {
	var satellites []models.Satellite
	result := r.db.DbHandler.Table("satellites").
		Joins("JOIN context_satellites ON satellites.id = context_satellites.satellite_id").
		Where("context_satellites.context_id = ?", contextID).
		Find(&satellites)

	if result.Error != nil {
		return nil, result.Error
	}

	var domainSatellites []domain.Satellite
	for _, satellite := range satellites {
		domainSatellites = append(domainSatellites, models.MapToSatelliteDomain(satellite))
	}
	return domainSatellites, nil
}

func (r *SatelliteRepository) FindAllWithPagination(ctx context.Context, page int, pageSize int, searchRequest *domain.SearchRequest) ([]domain.Satellite, int64, error) {
	var results []models.Satellite
	var totalRecords int64

	// Calculate the offset
	offset := (page - 1) * pageSize

	query := r.db.DbHandler.Table("satellites").
		Where("deleted_at IS NULL")

	// Apply search filtering if a wildcard is provided
	if searchRequest != nil && searchRequest.Wildcard != "" {
		wildcard := "%" + searchRequest.Wildcard + "%"
		query = query.Where(
			"LOWER(name) LIKE LOWER(?) OR LOWER(space_id) LIKE LOWER(?)",
			wildcard, wildcard,
		)
	}

	// Count total records
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Paginate and retrieve results
	if err := query.Limit(pageSize).Offset(offset).Find(&results).Error; err != nil {
		return nil, 0, err
	}

	// Map results to domain.Satellite
	var satellites []domain.Satellite
	for _, result := range results {
		satellites = append(satellites, models.MapToSatelliteDomain(result))
	}

	return satellites, totalRecords, nil
}

// FetchAndLockSatellites retrieves and locks available satellites within active contexts or returns satellites already locked by the same lockedBy.
func (r *SatelliteRepository) FetchAndLockSatellites(ctx context.Context, lockedBy string, maxNbSatellites int64) ([]string, error) {
	query := `
		WITH existing_locks AS (
			SELECT s.space_id
			FROM satellites s
			INNER JOIN context_satellites cs ON s.space_id = cs.satellite_id
			INNER JOIN contexts c ON cs.context_id = c.id
			WHERE s.locked = TRUE
			AND s.locked_by = $1
			AND c.is_active = TRUE
			AND c.deleted_at IS NULL
		), new_locks AS (
			UPDATE satellites s
			SET locked = TRUE, locked_by = $1
			FROM context_satellites cs
			INNER JOIN contexts c ON cs.context_id = c.id
			WHERE s.space_id = cs.satellite_id
			AND s.locked = FALSE
			AND c.is_active = TRUE
			AND c.deleted_at IS NULL
			RETURNING s.space_id
		)
		SELECT space_id FROM existing_locks
		UNION ALL
		SELECT space_id FROM new_locks
		LIMIT $2;
	`

	var lockedSatellites []string

	// Execute query and retrieve locked satellite IDs
	err := r.db.DbHandler.Select(ctx, &lockedSatellites, query, lockedBy, maxNbSatellites)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch and lock satellites: %w", err.Error)
	}

	if len(lockedSatellites) == 0 {
		log.Warn("No satellites available for locking in active contexts")
		return nil, nil
	}

	log.Infof("🔒 Successfully locked %d satellites (including previously locked by %s)", len(lockedSatellites), lockedBy)
	return lockedSatellites, nil
}

// ReleaseSatellites unlocks the satellites after processing.
func (r *SatelliteRepository) ReleaseSatellites(ctx context.Context, satelliteIDs []string) error {
	if len(satelliteIDs) == 0 {
		log.Warn("No satellites to release")
		return nil
	}

	query := `
		UPDATE satellites
		SET locked = FALSE
		WHERE space_id = ANY($1);
	`

	err := r.db.DbHandler.Exec(query, pq.Array(satelliteIDs))
	if err != nil {
		return fmt.Errorf("failed to release satellite locks: %w", err.Error)
	}

	log.Infof("🔓 Released lock on %d satellites", len(satelliteIDs))
	return nil
}
