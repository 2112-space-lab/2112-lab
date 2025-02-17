package models

import (
	"github.com/org/2112-space-lab/org/app-service/internal/domain"
)

// GlobalProperty represents a configurable system property stored in the database.
type GlobalProperty struct {
	ModelBase
	Key         string `gorm:"primaryKey;size:255;not null" json:"key"` // Unique configuration key
	Value       string `gorm:"type:text;not null" json:"value"`         // Configuration value
	Description string `gorm:"type:text" json:"description"`            // Description of the property
	ValueType   string `gorm:"size:50;not null" json:"value_type"`      // Type of value (bool, int, float, string, duration)
}

// MapToGlobalPropertyDomain converts an GlobalProperty database model to a domain GlobalProperty model.
func MapToGlobalPropertyDomain(a GlobalProperty) domain.GlobalProperty {
	return domain.GlobalProperty{
		ModelBase: domain.ModelBase{
			ID:          a.ID,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   &a.UpdatedAt,
			DeleteAt:    a.DeleteAt,
			ProcessedAt: a.ProcessedAt,
			IsActive:    a.IsActive,
			IsFavourite: a.IsFavourite,
			DisplayName: a.DisplayName,
		},
		Key:         a.Key,
		Value:       a.Value,
		Description: a.Description,
		ValueType:   a.ValueType,
	}
}

// MapToGlobalPropertyModel converts a domain GlobalProperty model to an GlobalProperty database model.
func MapToGlobalPropertyModel(a domain.GlobalProperty) GlobalProperty {
	return GlobalProperty{
		ModelBase: ModelBase{
			ID:          a.ModelBase.ID,
			CreatedAt:   a.ModelBase.CreatedAt,
			UpdatedAt:   *a.ModelBase.UpdatedAt,
			DeleteAt:    a.ModelBase.DeleteAt,
			ProcessedAt: a.ModelBase.ProcessedAt,
			IsActive:    a.ModelBase.IsActive,
			IsFavourite: a.ModelBase.IsFavourite,
			DisplayName: a.ModelBase.DisplayName,
		},
		Key:         a.Key,
		Value:       a.Value,
		Description: a.Description,
		ValueType:   a.ValueType,
	}
}
