package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
)

const (
	MaxSatellitesPerEventDetectorDefault = 100
	DefaultSimulationInterval            = 30 * time.Second
	DefaultSimulationSteps               = 10
	DefdaultSimulationDuration           = 60 * time.Minute
	DefaultSimulationBufferDuration      = 3 * time.Hour
)

// GlobalPropertyRepository manages retrieval of configuration properties.
type GlobalPropertyRepository struct {
	db *data.Database
}

// NewGlobalPropertyRepository initializes the repository.
func NewGlobalPropertyRepository(db *data.Database) GlobalPropertyRepository {
	return GlobalPropertyRepository{db: db}
}

// GetProperty retrieves a property by key.
func (r *GlobalPropertyRepository) GetProperty(ctx context.Context, key string) (*models.GlobalProperty, error) {
	var prop models.GlobalProperty
	if err := r.db.DbHandler.WithContext(ctx).Where("key = ?", key).First(&prop).Error; err != nil {
		return nil, fmt.Errorf("property [%s] not found: %w", key, err)
	}
	return &prop, nil
}

// GetAllProperties retrieves all properties as a map.
func (r *GlobalPropertyRepository) GetAllProperties(ctx context.Context) (map[string]string, error) {
	var properties []models.GlobalProperty
	result := r.db.DbHandler.WithContext(ctx).Find(&properties)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch properties: %w", result.Error)
	}

	propMap := make(map[string]string)
	for _, prop := range properties {
		propMap[prop.Key] = prop.Value
	}
	return propMap, nil
}

// SetProperty inserts or updates a property.
func (r *GlobalPropertyRepository) SetProperty(ctx context.Context, key, value, valueType, description string) error {
	// Upsert property
	return r.db.DbHandler.WithContext(ctx).Save(&models.GlobalProperty{
		Key:         key,
		Value:       value,
		ValueType:   valueType,
		Description: description,
	}).Error
}

// GetBool retrieves a boolean property.
func (r *GlobalPropertyRepository) GetBool(ctx context.Context, key string, defaultValue bool) (bool, error) {
	prop, err := r.GetProperty(ctx, key)
	if err != nil {
		return defaultValue, err
	}
	parsed, err := strconv.ParseBool(prop.Value)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid boolean value for [%s]: %w", key, err)
	}
	return parsed, nil
}

// GetInt retrieves an integer property.
func (r *GlobalPropertyRepository) GetInt(ctx context.Context, key string, defaultValue int64) (int64, error) {
	prop, err := r.GetProperty(ctx, key)
	if err != nil {
		return defaultValue, err
	}
	parsed, err := strconv.ParseInt(prop.Value, 10, 64)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid integer value for [%s]: %w", key, err)
	}
	return parsed, nil
}

// GetFloat retrieves a float property.
func (r *GlobalPropertyRepository) GetFloat(ctx context.Context, key string, defaultValue float64) (float64, error) {
	prop, err := r.GetProperty(ctx, key)
	if err != nil {
		return defaultValue, err
	}
	parsed, err := strconv.ParseFloat(prop.Value, 64)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid float value for [%s]: %w", key, err)
	}
	return parsed, nil
}

// GetDuration retrieves a time.Duration property (milliseconds).
func (r *GlobalPropertyRepository) GetDuration(ctx context.Context, key string, defaultValue time.Duration) (time.Duration, error) {
	prop, err := r.GetProperty(ctx, key)
	if err != nil {
		return defaultValue, err
	}
	parsed, err := strconv.ParseInt(prop.Value, 10, 64)
	if err != nil {
		return defaultValue, fmt.Errorf("invalid duration value for [%s]: %w", key, err)
	}
	return time.Duration(parsed) * time.Millisecond, nil
}

// GetMaxSatellitesPerEventDetector retrieves the max number of satellites an event detector can handle.
func (r *GlobalPropertyRepository) GetMaxSatellitesPerEventDetector(ctx context.Context, defaultValue int64) (int64, error) {
	return r.GetInt(ctx, "event_detector_per_max_satellites", defaultValue)
}

// GetEventDetectorSimulationInterval retrieves the simulation interval for event detectors.
func (r *GlobalPropertyRepository) GetEventDetectorSimulationInterval(ctx context.Context, defaultValue time.Duration) (time.Duration, error) {
	return r.GetDuration(ctx, "event_detector_simulation_interval", defaultValue)
}

// GetEventDetectorSimulationInterval retrieves the simulation duration for event detectors.
func (r *GlobalPropertyRepository) GetEventDetectorSimulationDuration(ctx context.Context, defaultValue time.Duration) (time.Duration, error) {
	return r.GetDuration(ctx, "event_detector_simulation_duration", defaultValue)
}

// GetEventDetectorSimulationSteps retrieves the simulation steps for event detectors.
func (r *GlobalPropertyRepository) GetEventDetectorSimulationSteps(ctx context.Context, defaultValue int64) (int64, error) {
	return r.GetInt(ctx, "event_detector_simulation_steps", defaultValue)
}

// GetEventDetectorSimulationBufferDuration retrieves the buffer duration for event detectors.
func (r *GlobalPropertyRepository) GetEventDetectorSimulationBufferDuration(ctx context.Context, defaultValue time.Duration) (time.Duration, error) {
	return r.GetDuration(ctx, "event_detector_simulation_buffer_duration", defaultValue)
}
