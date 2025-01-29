package models

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	fx "github.com/Elbujito/2112/src/app-service/pkg/option"
	xtime "github.com/Elbujito/2112/src/app-service/pkg/time"
)

// AuditTrail represents the database model for audit trails.
type AuditTrail struct {
	ModelBase
	TableName   string    `gorm:"size:255;not null"` // Name of the table affected
	RecordID    string    `gorm:"size:255;not null"` // ID of the record affected
	Action      string    `gorm:"size:50;not null"`  // Action performed (e.g., "INSERT", "UPDATE", "DELETE")
	ChangesJSON string    `gorm:"type:json"`         // JSON representation of changes made
	PerformedBy string    `gorm:"size:255;not null"` // User or system that performed the action
	PerformedAt time.Time `gorm:"not null"`          // Timestamp of the action
}

// MapToAuditTrailDomain converts an AuditTrail database model to a domain AuditTrail model.
func MapToAuditTrailDomain(a AuditTrail) domain.AuditTrail {
	return domain.AuditTrail{
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
		TableName:   a.TableName,
		RecordID:    a.RecordID,
		Action:      domain.ActionType(a.Action),
		ChangesJSON: fx.AsOption(&a.ChangesJSON),
		PerformedBy: domain.UserID(a.PerformedBy),
		PerformedAt: xtime.NewUtcTimeIgnoreZone(a.PerformedAt),
	}
}

// MapToAuditTrailModel converts a domain AuditTrail model to an AuditTrail database model.
func MapToAuditTrailModel(a domain.AuditTrail) AuditTrail {
	return AuditTrail{
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
		TableName:   a.TableName,
		RecordID:    a.RecordID,
		Action:      string(a.Action),
		ChangesJSON: fx.GetOrDefault(a.ChangesJSON, ""),
		PerformedBy: string(a.PerformedBy),
		PerformedAt: a.PerformedAt.Inner(),
	}
}
