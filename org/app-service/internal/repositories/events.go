package repository

import (
	"context"

	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EventRepository manages audit trail data access.
type EventRepository struct {
	db *data.Database
}

// NewEventRepository creates a new EventRepository instance.
func NewEventRepository(db *data.Database) EventRepository {
	return EventRepository{db: db}
}

// Save inserts a new event handler execution record, ignoring duplicates.
func (r *EventRepository) Save(ctx context.Context, event domain.Event) error {
	return r.db.DbHandler.Transaction(func(tx *gorm.DB) error {
		model := models.MapToEventModel(event)

		err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&model).Error

		return err
	})
}
