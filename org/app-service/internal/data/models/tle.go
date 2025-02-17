package models

import (
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
)

// TLE Model
type TLE struct {
	ModelBase
	SpaceID string    `gorm:"size:255;not null;index"` // Foreign key to Satellite table via Norad ID
	Line1   string    `gorm:"size:255;not null"`
	Line2   string    `gorm:"size:255;not null"`
	Epoch   time.Time `gorm:"not null"` // Time associated with the TLE
}

// MapToDomain converts a models.Tile to a domain.Tile.
func MapToTLEDomain(t TLE) domain.TLE {
	// Convert to domain
	return domain.TLE{
		ModelBase: domain.ModelBase{
			ID:          t.ID,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   &t.UpdatedAt,
			DeleteAt:    t.DeleteAt,
			ProcessedAt: t.ProcessedAt,
			IsActive:    t.IsActive,
			IsFavourite: t.IsFavourite,
			DisplayName: t.DisplayName,
		},
		SpaceID: t.SpaceID,
		Line1:   t.Line1,
		Line2:   t.Line2,
		Epoch:   t.Epoch,
	}
}

// MapToTLEModel converts a domain.TLE to a models.TLE.
func MapToTLEModel(t domain.TLE) TLE {
	return TLE{
		ModelBase: ModelBase{
			ID:          t.ModelBase.ID,
			CreatedAt:   t.ModelBase.CreatedAt,
			UpdatedAt:   *t.ModelBase.UpdatedAt,
			DeleteAt:    t.ModelBase.DeleteAt,
			ProcessedAt: t.ModelBase.ProcessedAt,
			IsActive:    t.ModelBase.IsActive,
			IsFavourite: t.ModelBase.IsFavourite,
			DisplayName: t.ModelBase.DisplayName,
		},
		SpaceID: t.SpaceID,
		Line1:   t.Line1,
		Line2:   t.Line2,
		Epoch:   t.Epoch,
	}
}
