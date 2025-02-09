package repository

import (
	"context"
	"errors"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SatelliteTLEAggregate is used for joining Satellite and TLE data.
type SatelliteTLEAggregate struct {
	// Satellite fields
	ID             string     `gorm:"column:id"`
	Name           string     `gorm:"column:name"`
	NoradID        string     `gorm:"column:norad_id"`
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
	db *data.Database
}

// NewSatelliteRepository creates a new SatelliteRepository instance.
func NewSatelliteRepository(db *data.Database) SatelliteRepository {
	return SatelliteRepository{db: db}
}

// FindByNoradID retrieves a satellite by its NORAD ID, excluding deleted ones.
func (r *SatelliteRepository) FindByNoradID(ctx context.Context, noradID string) (domain.Satellite, error) {
	var satellite models.Satellite
	result := r.db.DbHandler.Where("norad_id = ? AND deleted_at IS NULL", noradID).First(&satellite)
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

// DeleteByNoradID marks a satellite record as deleted.
func (r *SatelliteRepository) DeleteByNoradID(ctx context.Context, noradID string) error {
	return r.db.DbHandler.Model(&models.Satellite{}).
		Where("norad_id = ?", noradID).
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
			Columns:   []clause.Column{{Name: "norad_id"}},
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
			satellites.id, satellites.name, satellites.norad_id, satellites.owner,
			satellites.launch_date, satellites.decay_date, satellites.international_designator,
			satellites.object_type, satellites.period, satellites.inclination, satellites.apogee,
			satellites.perigee, satellites.rcs, satellites.altitude, satellites.is_active,
			satellites.created_at, satellites.updated_at, satellites.processed_at, satellites.is_favourite,
			latest_tles.line1, latest_tles.line2, latest_tles.updated_at AS tle_updated_at
		`).
		Joins(`LEFT JOIN (
			SELECT t1.norad_id, t1.line1, t1.line2, t1.updated_at
			FROM tles t1
			WHERE t1.updated_at = (
				SELECT MAX(t2.updated_at)
				FROM tles t2
				WHERE t2.norad_id = t1.norad_id
			)
		) AS latest_tles ON satellites.norad_id = latest_tles.norad_id`).
		Where("satellites.deleted_at IS NULL")

	// Apply search filters if a wildcard search is provided
	if searchRequest != nil && searchRequest.Wildcard != "" {
		wildcard := "%" + searchRequest.Wildcard + "%"
		query = query.Where("LOWER(satellites.name) LIKE LOWER(?) OR LOWER(satellites.norad_id) LIKE LOWER(?)", wildcard, wildcard)
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
			NoradID:        result.NoradID,
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
			"LOWER(name) LIKE LOWER(?) OR LOWER(norad_id) LIKE LOWER(?)",
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
