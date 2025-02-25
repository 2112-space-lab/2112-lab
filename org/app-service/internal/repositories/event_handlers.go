package repository

import (
	"context"

	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EventHandlerRepository manages event handler execution records.
type EventHandlerRepository struct {
	db *data.Database
}

// NewEventHandlerRepository creates a new EventHandlerRepository instance.
func NewEventHandlerRepository(db *data.Database) EventHandlerRepository {
	return EventHandlerRepository{db: db}
}

// Save inserts a new event handler execution record, ignoring duplicates.
func (r *EventHandlerRepository) Save(ctx context.Context, handler domain.EventHandler) error {
	return r.db.DbHandler.Transaction(func(tx *gorm.DB) error {
		model := models.MapToEventHandlerModel(handler)

		err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&model).Error

		return err
	})
}
