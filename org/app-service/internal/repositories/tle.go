package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/redis"
	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xtime"
	"gorm.io/gorm/clause"
)

// TleRepository implements the TLERepository interface with caching and database operations.
type TleRepository struct {
	db          *data.Database
	redisClient *redis.RedisClient
}

// NewTLERepository initializes the repository with a cache TTL.
func NewTLERepository(db *data.Database, redisClient *redis.RedisClient) TleRepository {
	return TleRepository{db: db, redisClient: redisClient}
}

// mapToDomainTLE converts a models.TLE to a domain.TLE.
func mapToDomainTLE(model models.TLE) domain.TLE {
	return domain.TLE{
		ID:      model.ID,
		NoradID: model.NoradID,
		Line1:   model.Line1,
		Line2:   model.Line2,
		Epoch:   model.Epoch,
	}
}

// mapToModelTLE converts a domain.TLE to a models.TLE.
func mapToModelTLE(domainTLE domain.TLE) models.TLE {
	return models.TLE{
		NoradID: domainTLE.NoradID,
		Line1:   domainTLE.Line1,
		Line2:   domainTLE.Line2,
		Epoch:   domainTLE.Epoch,
	}
}

// GetTle retrieves a TLE from cache or database.
func (r *TleRepository) GetTle(ctx context.Context, noradID string) (domain.TLE, error) {
	key := fmt.Sprintf("satellite:tle:%s", noradID)

	// Check Redis cache
	data, err := r.redisClient.HGetAll(ctx, key)
	if err == nil && len(data) > 0 {
		epoch, parseErr := xtime.ParseEpoch(data["epoch"])
		if parseErr == nil {
			return domain.TLE{
				ID:    noradID,
				Line1: data["line_1"],
				Line2: data["line_2"],
				Epoch: epoch,
			}, nil
		}
	}

	// Fallback to database
	var modelTLE models.TLE
	result := r.db.DbHandler.First(&modelTLE, "norad_id = ?", noradID)
	if result.Error != nil {
		return domain.TLE{}, result.Error
	}

	tle := mapToDomainTLE(modelTLE)

	// Update Redis cache
	r.updateCache(ctx, key, tle)

	return tle, nil
}

// SaveTle saves a TLE to the database and updates the cache.
func (r *TleRepository) SaveTle(ctx context.Context, tle domain.TLE) error {
	modelTLE := mapToModelTLE(tle)
	if err := r.db.DbHandler.Create(&modelTLE).Error; err != nil {
		return err
	}

	key := fmt.Sprintf("satellite:tle:%s", tle.ID)
	r.updateCache(ctx, key, tle)

	return r.publishTleToBroker(ctx, tle)
}

func (r *TleRepository) UpdateTleBatch(ctx context.Context, tles []domain.TLE) error {
	if len(tles) == 0 {
		return fmt.Errorf("no TLEs to update")
	}

	const batchSize = 1000 // Process TLEs in batches of 1000
	for i := 0; i < len(tles); i += batchSize {
		end := i + batchSize
		if end > len(tles) {
			end = len(tles)
		}

		batch := tles[i:end]

		// Map TLEs to the database model
		modelTLEs := make([]models.TLE, len(batch))
		for j, tle := range batch {
			modelTLEs[j] = mapToModelTLE(tle)
		}

		// Batch upsert into the database
		if err := r.db.DbHandler.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Save(&modelTLEs).Error; err != nil {
			log.Errorf("Failed to batch upsert TLEs: %v\n", err)
			return err
		}

		// Process Redis caching and broker publishing
		for _, tle := range batch {
			key := fmt.Sprintf("satellite:tle:%s", tle.ID)
			cacheData := map[string]interface{}{
				"line_1": tle.Line1,
				"line_2": tle.Line2,
				"epoch":  tle.Epoch,
				"id":     tle.NoradID,
			}

			// Update Redis cache
			if err := r.redisClient.HSet(ctx, key, cacheData); err != nil {
				log.Errorf("Failed to update Redis cache for key %s: %v\n", key, err)
			}
			// Publish to the message broker
			if err := r.publishTleToBroker(ctx, tle); err != nil {
				log.Errorf("Failed to publish TLE to message broker for NORAD ID %s: %v\n", tle.NoradID, err)
			}
		}
	}

	return nil
}

// DeleteTle deletes a TLE from the database and invalidates the cache.
func (r *TleRepository) DeleteTle(ctx context.Context, id string) error {
	if err := r.db.DbHandler.Delete(&models.TLE{}, "id = ?", id).Error; err != nil {
		return err
	}

	key := fmt.Sprintf("satellite:tle:%s", id)
	if err := r.redisClient.Del(ctx, key); err != nil {
		log.Errorf("Failed to delete Redis cache for key %s: %v\n", key, err)
	}
	return nil
}

// AssociateTLEWithContext associates a TLE with a specific context.
func (r *TleRepository) AssociateTLEWithContext(ctx context.Context, contextID string, tleID string) error {
	contextTLE := models.ContextTLE{
		ContextID: contextID,
		TLEID:     tleID,
	}

	// Insert into context_tles table
	if err := r.db.DbHandler.Create(&contextTLE).Error; err != nil {
		return fmt.Errorf("failed to associate TLE with context: %w", err)
	}
	return nil
}

// GetTLEsByContext retrieves all TLEs associated with a specific context.
func (r *TleRepository) GetTLEsByContext(ctx context.Context, contextID string) ([]domain.TLE, error) {
	var contextTLEs []models.ContextTLE

	// Query the many-to-many relationship
	if err := r.db.DbHandler.Where("context_id = ?", contextID).Find(&contextTLEs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve TLEs by context: %w", err)
	}

	// Fetch the TLEs based on the IDs
	tleIDs := make([]string, len(contextTLEs))
	for i, contextTLE := range contextTLEs {
		tleIDs[i] = contextTLE.TLEID
	}

	var tles []models.TLE
	if err := r.db.DbHandler.Where("id IN ?", tleIDs).Find(&tles).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve TLE details: %w", err)
	}

	// Convert to domain models
	domainTLEs := make([]domain.TLE, len(tles))
	for i, tle := range tles {
		domainTLEs[i] = mapToDomainTLE(tle)
	}

	return domainTLEs, nil
}

// RemoveTLEFromContext removes the association between a TLE and a context.
func (r *TleRepository) RemoveTLEFromContext(ctx context.Context, contextID string, tleID string) error {
	if err := r.db.DbHandler.Where("context_id = ? AND tle_id = ?", contextID, tleID).
		Delete(&models.ContextTLE{}).Error; err != nil {
		return fmt.Errorf("failed to remove TLE from context: %w", err)
	}
	return nil
}

// QuerySatellitePositions retrieves satellite positions from Redis within a time range.
func (r *TleRepository) QuerySatellitePositions(ctx context.Context, satelliteID string, startTime, endTime time.Time) ([]domain.SatellitePosition, error) {
	key := fmt.Sprintf("satellite_positions:%s", satelliteID)

	startTimestamp := strconv.FormatInt(startTime.Unix(), 10)
	endTimestamp := strconv.FormatInt(endTime.Unix(), 10)

	results, err := r.redisClient.ZRangeByScore(ctx, key, startTimestamp, endTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to query Redis for satellite positions: %w", err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	var positions []domain.SatellitePosition
	for _, result := range results {
		var position domain.SatellitePosition
		if err := json.Unmarshal([]byte(result), &position); err != nil {
			log.Errorf("Failed to parse satellite position: %v\n", err)
			continue
		}
		positions = append(positions, position)
	}

	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Timestamp.Before(positions[j].Timestamp)
	})

	return positions, nil
}

// updateCache updates the Redis cache for a TLE.
func (r *TleRepository) updateCache(ctx context.Context, key string, tle domain.TLE) {
	cacheData := map[string]interface{}{
		"line_1": tle.Line1,
		"line_2": tle.Line2,
		"epoch":  tle.Epoch,
		"id":     tle.NoradID,
	}
	if err := r.redisClient.HSet(ctx, key, cacheData); err != nil {
		log.Errorf("Failed to update Redis cache for key %s: %v\n", key, err)
	}
}

// publishTleToBroker sends TLE updates to the message broker.
func (r *TleRepository) publishTleToBroker(ctx context.Context, tle domain.TLE) error {
	message := map[string]interface{}{
		"id":     tle.NoradID,
		"line_1": tle.Line1,
		"line_2": tle.Line2,
		"epoch":  tle.Epoch,
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	channel := "satellite_tle_updates"
	if err := r.redisClient.Publish(ctx, channel, messageJSON); err != nil {
		return fmt.Errorf("failed to publish message to channel %s: %w", channel, err)
	}

	log.Debugf("Successfully published TLE update to channel %s\n", channel)
	return nil
}
