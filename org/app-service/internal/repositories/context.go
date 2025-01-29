package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/data"
	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"gorm.io/gorm"
)

// ContextRepository manages context data access.
type ContextRepository struct {
	db *data.Database
}

// NewContextRepository creates a new ContextRepository instance.
func NewContextRepository(db *data.Database) domain.GameContextRepository {
	return &ContextRepository{db: db}
}

// Save creates a new context record.
func (r *ContextRepository) Save(ctx context.Context, context domain.GameContext) error {
	model := models.MapToContextModel(context)
	return r.db.DbHandler.Create(&model).Error
}

// Update modifies an existing context record.
func (r *ContextRepository) Update(ctx context.Context, context domain.GameContext) error {
	model := models.MapToContextModel(context)
	return r.db.DbHandler.Save(&model).Error
}

// FindByUniqueName retrieves a context by its unique name.
func (r *ContextRepository) FindByUniqueName(ctx context.Context, gameContextName domain.GameContextName) (domain.GameContext, error) {
	var model models.Context
	result := r.db.DbHandler.First(&model, "name = ? AND deleted_at IS NULL", string(gameContextName))
	if result.Error != nil {
		return domain.GameContext{}, result.Error
	}
	return models.MapToContextDomain(model), nil
}

// FindAll retrieves all contexts.
func (r *ContextRepository) FindAll(ctx context.Context) ([]domain.GameContext, error) {
	var results []models.Context
	result := r.db.DbHandler.Find(&results, "deleted_at IS NULL")
	if result.Error != nil {
		return nil, result.Error
	}

	var contexts []domain.GameContext
	for _, model := range results {
		context := models.MapToContextDomain(model)
		contexts = append(contexts, context)
	}
	return contexts, nil
}

// FindAllWithPagination retrieves all contexts with pagination and optional filtering by name or description.
func (r *ContextRepository) FindAllWithPagination(ctx context.Context, page int, pageSize int, wildcard string) ([]domain.GameContext, error) {
	var results []models.Context

	// Construct the query with pagination
	query := r.db.DbHandler.Scopes(models.Paginate(page, pageSize)).Where("deleted_at IS NULL")

	// Apply wildcard filter for both name and description if a wildcard is provided
	if wildcard != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+wildcard+"%", "%"+wildcard+"%")
	}

	// Execute the query
	result := query.Find(&results)

	if result.Error != nil {
		return nil, result.Error
	}

	// Map the results to domain models
	var contexts []domain.GameContext
	for _, model := range results {
		context := models.MapToContextDomain(model)
		contexts = append(contexts, context)
	}
	return contexts, nil
}

// DeleteByUniqueName marks a context record as deleted by unique name.
func (r *ContextRepository) DeleteByUniqueName(ctx context.Context, name string) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", name).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindActiveBySatelliteID retrieves the single active context associated with a satellite.
// Raises an error if multiple active contexts are found.
func (r *ContextRepository) FindActiveBySatelliteID(ctx context.Context, satelliteID domain.SatelliteID) (domain.GameContext, error) {
	var query []models.Context
	result := r.db.DbHandler.Table("contexts").
		Joins("JOIN context_satellites ON contexts.id = context_satellites.context_id").
		Where("context_satellites.satellite_id = ? AND contexts.is_active = TRUE AND contexts.deleted_at IS NULL", string(satelliteID)).
		Find(&query)

	if result.Error != nil {
		return domain.GameContext{}, result.Error
	}

	// Check for multiple active contexts
	if len(query) > 1 {
		return domain.GameContext{}, fmt.Errorf("multiple active contexts found for satellite ID %s", satelliteID)
	}

	// Check if no active contexts were found
	if len(query) == 0 {
		return domain.GameContext{}, fmt.Errorf("no active context found for satellite ID %s", satelliteID)
	}

	// Return the single active context
	return models.MapToContextDomain(query[0]), nil
}

// AssignSatellite associates a satellite with a context.
func (r *ContextRepository) AssignSatellite(ctx context.Context, gameContextName domain.GameContextName, satelliteID domain.SatelliteID) error {
	contextSatellite := models.ContextSatellite{
		ContextID:   string(gameContextName),
		SatelliteID: string(satelliteID),
	}
	return r.db.DbHandler.Create(&contextSatellite).Error
}

// AssignSatellites associates a list of satellites with a GameContext.
func (r *ContextRepository) AssignSatellites(ctx context.Context, gameContextName domain.GameContextName, satelliteIDs []domain.SatelliteID) error {
	var contextSatellites []models.ContextSatellite

	// Create ContextSatellite records for each satellite ID
	for _, satelliteID := range satelliteIDs {
		contextSatellite := models.ContextSatellite{
			ContextID:   string(gameContextName),
			SatelliteID: string(satelliteID),
		}
		contextSatellites = append(contextSatellites, contextSatellite)
	}

	// Use GORM's Create method to insert all records in a single operation
	return r.db.DbHandler.Create(&contextSatellites).Error
}

// RemoveSatellite removes the association between a satellite and a context.
func (r *ContextRepository) RemoveSatellite(ctx context.Context, gameContextName domain.GameContextName, satelliteID domain.SatelliteID) error {
	return r.db.DbHandler.Where("context_id = ? AND satellite_id = ?", string(gameContextName), string(satelliteID)).
		Delete(&models.ContextSatellite{}).Error
}

// RemoveSatellites removes the associations between a list of satellites and a GameContext.
func (r *ContextRepository) RemoveSatellites(ctx context.Context, gameContextName domain.GameContextName, satelliteIDs []domain.SatelliteID) error {
	// Convert satelliteIDs to a slice of strings for query compatibility
	var satelliteIDStrings []string
	for _, satelliteID := range satelliteIDs {
		satelliteIDStrings = append(satelliteIDStrings, string(satelliteID))
	}

	// Perform batch delete operation
	return r.db.DbHandler.
		Where("context_id = ? AND satellite_id IN ?", string(gameContextName), satelliteIDStrings).
		Delete(&models.ContextSatellite{}).Error
}

// DesactiveContext marks a context as inactive.
func (r *ContextRepository) DesactiveContext(ctx context.Context, gameContextName domain.GameContextName) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("is_active", false).Error
}

// ActivateContext marks a context as active.
func (r *ContextRepository) ActivateContext(ctx context.Context, gameContextName domain.GameContextName) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("is_active", true).Error
}

// GetActiveContext retrieves the currently active context.
func (r *ContextRepository) GetActiveContext(ctx context.Context) (domain.GameContext, error) {
	var query []models.Context
	result := r.db.DbHandler.Table("contexts").
		Where("is_active = TRUE AND deleted_at IS NULL").
		Find(&query)

	if result.Error != nil {
		return domain.GameContext{}, result.Error
	}

	if len(query) > 1 {
		return domain.GameContext{}, fmt.Errorf("multiple contexts active at the same time [%+v]", query)
	}

	if len(query) == 0 {
		return domain.GameContext{}, fmt.Errorf("no active context found")
	}

	return models.MapToContextDomain(query[0]), nil
}

// SetActivatedAt sets the ActivatedAt timestamp for a context.
func (r *ContextRepository) SetActivatedAt(ctx context.Context, gameContextName domain.GameContextName, activatedAt time.Time) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("activated_at", activatedAt).Error
}

// UnsetActivatedAt clears the ActivatedAt timestamp for a context.
func (r *ContextRepository) UnsetActivatedAt(ctx context.Context, gameContextName domain.GameContextName) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("activated_at", nil).Error
}

// SetDesactivatedAt sets the DesactivatedAt timestamp for a context.
func (r *ContextRepository) SetDesactivatedAt(ctx context.Context, gameContextName domain.GameContextName, desactivatedAt time.Time) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("desactivated_at", desactivatedAt).Error
}

// UnsetDesactivatedAt clears the DesactivatedAt timestamp for a context.
func (r *ContextRepository) UnsetDesactivatedAt(ctx context.Context, gameContextName domain.GameContextName) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("desactivated_at", nil).Error
}

// SetTriggerGeneratedMappingAt sets the TriggerGeneratedMappingAt timestamp for a context.
func (r *ContextRepository) SetTriggerGeneratedMappingAt(ctx context.Context, gameContextName domain.GameContextName, timestamp time.Time) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("trigger_generated_mapping_at", timestamp).Error
}

// UnsetTriggerGeneratedMappingAt clears the TriggerGeneratedMappingAt timestamp for a context.
func (r *ContextRepository) UnsetTriggerGeneratedMappingAt(ctx context.Context, gameContextName domain.GameContextName) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("trigger_generated_mapping_at", nil).Error
}

// SetTriggerImportedTLEAt sets the TriggerImportedTLEAt timestamp for a context.
func (r *ContextRepository) SetTriggerImportedTLEAt(ctx context.Context, gameContextName domain.GameContextName, timestamp time.Time) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("trigger_imported_tle_at", timestamp).Error
}

// UnsetTriggerImportedTLEAt clears the TriggerImportedTLEAt timestamp for a context.
func (r *ContextRepository) UnsetTriggerImportedTLEAt(ctx context.Context, gameContextName domain.GameContextName) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("trigger_imported_tle_at", nil).Error
}

// SetTriggerImportedSatelliteAt sets the TriggerImportedSatelliteAt timestamp for a context.
func (r *ContextRepository) SetTriggerImportedSatelliteAt(ctx context.Context, gameContextName domain.GameContextName, timestamp time.Time) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("trigger_imported_satellite_at", timestamp).Error
}

// UnsetTriggerImportedSatelliteAt clears the TriggerImportedSatelliteAt timestamp for a context.
func (r *ContextRepository) UnsetTriggerImportedSatelliteAt(ctx context.Context, gameContextName domain.GameContextName) error {
	return r.db.DbHandler.Model(&models.Context{}).
		Where("name = ?", string(gameContextName)).
		Update("trigger_imported_satellite_at", nil).Error
}
