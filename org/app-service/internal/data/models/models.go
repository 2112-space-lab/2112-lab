package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db *gorm.DB

// Init initializes the database connection in the models package.
func Init(database *gorm.DB) {
	db = database
}

// ModelBase provides a base structure with common fields for all models.
type ModelBase struct {
	ID          string     `gorm:"type:char(36);primary_key;"` // Unique identifier
	DisplayName string     `gorm:"type:varchar(255);not null"` // Human-readable name
	CreatedAt   time.Time  // Record creation timestamp
	UpdatedAt   time.Time  // Record update timestamp
	DeleteAt    *time.Time // Soft delete timestamp
	ProcessedAt *time.Time // Custom timestamp for processing
	IsActive    bool       `gorm:"not null;default:true"`  // Active status
	IsFavourite bool       `gorm:"not null;default:false"` // Favorite status
}

// BeforeCreate ensures a UUID is generated for the ID field before creating the record.
func (base *ModelBase) BeforeCreate(tx *gorm.DB) error {
	if base.ID == "" {
		base.ID = uuid.NewString()
	}
	return nil
}

// Paginate adds pagination to queries.
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Ensure valid page number
		if page <= 0 {
			page = 1
		}
		// Ensure page size is within limits
		if pageSize <= 0 {
			pageSize = 10
		} else if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
