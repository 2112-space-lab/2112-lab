package models

import (
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	fx "github.com/org/2112-space-lab/org/app-service/pkg/option"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
)

// Context represents the database model for logical groupings.
type Context struct {
	ModelBase
	Name                       string     `gorm:"size:255;unique;not null"` // Unique name of the context
	TenantID                   string     `gorm:"size:255;not null;index"`  // Tenant identifier
	Description                string     `gorm:"size:1024"`                // Optional description of the context
	MaxSatellite               int        `gorm:"not null"`                 // Maximum number of satellites allowed
	MaxTiles                   int        `gorm:"not null"`                 // Maximum number of tiles allowed
	ActivatedAt                *time.Time // Time the context was activated
	DesactivatedAt             *time.Time // Time the context was deactivated
	TriggerGeneratedMappingAt  *time.Time // Time the mapping was generated
	TriggerImportedTLEAt       *time.Time // Time the TLE data was imported
	TriggerImportedSatelliteAt *time.Time // Time the satellite data was imported
}

// MapToContextDomain converts a Context database model to a GameContext domain model.
func MapToContextDomain(c Context) domain.GameContext {
	return domain.GameContext{
		ModelBase: domain.ModelBase{
			ID:          c.ID,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   &c.UpdatedAt,
			DeleteAt:    c.DeleteAt,
			ProcessedAt: c.ProcessedAt,
			IsActive:    c.IsActive,
			IsFavourite: c.IsFavourite,
			DisplayName: c.DisplayName,
		},
		TenantID: domain.TenantID(c.TenantID),
		Name:     domain.GameContextName(c.Name),
		Description: fx.ConvertOption(fx.AsOption(&c.Description), func(d string) domain.GameContextDescription {
			return domain.GameContextDescription(d)
		}),
		ActivatedAt:                xtime.ToUtcTime(c.ActivatedAt),
		DesactivatedAt:             xtime.ToUtcTime(c.DesactivatedAt),
		TriggerGeneratedMappingAt:  xtime.ToUtcTime(c.TriggerGeneratedMappingAt),
		TriggerImportedTLEAt:       xtime.ToUtcTime(c.TriggerImportedTLEAt),
		TriggerImportedSatelliteAt: xtime.ToUtcTime(c.TriggerImportedSatelliteAt),
	}
}

// MapToContextModel converts a GameContext domain model to a Context database model.
func MapToContextModel(c domain.GameContext) Context {
	return Context{
		ModelBase: ModelBase{
			ID:          c.ModelBase.ID,
			CreatedAt:   c.ModelBase.CreatedAt,
			UpdatedAt:   *c.ModelBase.UpdatedAt,
			DeleteAt:    c.ModelBase.DeleteAt,
			ProcessedAt: c.ModelBase.ProcessedAt,
			IsActive:    c.ModelBase.IsActive,
			IsFavourite: c.ModelBase.IsFavourite,
			DisplayName: c.ModelBase.DisplayName,
		},
		TenantID:                   string(c.TenantID),
		Name:                       string(c.Name),
		Description:                string(fx.GetOrDefault(c.Description, domain.GameContextDescription(""))),
		ActivatedAt:                xtime.ToTimePointer(c.ActivatedAt),
		DesactivatedAt:             xtime.ToTimePointer(c.DesactivatedAt),
		TriggerGeneratedMappingAt:  xtime.ToTimePointer(c.TriggerGeneratedMappingAt),
		TriggerImportedTLEAt:       xtime.ToTimePointer(c.TriggerImportedTLEAt),
		TriggerImportedSatelliteAt: xtime.ToTimePointer(c.TriggerImportedSatelliteAt),
	}
}

// Common conversion functions for working with database models

// ContextTile defines the many-to-many relationship between Context and Tile.
type ContextTile struct {
	ContextID string  `gorm:"not null;index"` // Foreign key to Context
	TileID    string  `gorm:"not null;index"` // Foreign key to Tile
	Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
	Tile      Tile    `gorm:"constraint:OnDelete:CASCADE;foreignKey:TileID;references:ID"`
}

// MapToContextTileDomain converts a ContextTile database model to a GameContextTile domain model.
func MapToContextTileDomain(ct ContextTile) domain.GameContextTile {
	return domain.GameContextTile{
		ContextID: ct.ContextID,
		TileID:    ct.TileID,
		Tile:      MapToTileDomain(ct.Tile),
	}
}

// MapToContextTileModel converts a GameContextTile domain model to a ContextTile database model.
func MapToContextTileModel(ct domain.GameContextTile) ContextTile {
	return ContextTile{
		ContextID: ct.ContextID,
		TileID:    ct.TileID,
		Tile:      MapFromDomain(ct.Tile),
	}
}

// ContextTLE defines the many-to-many relationship between Context and TLE.
type ContextTLE struct {
	ContextID string  `gorm:"not null;index"` // Foreign key to Context
	TLEID     string  `gorm:"not null;index"` // Foreign key to TLE
	Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
	TLE       TLE     `gorm:"constraint:OnDelete:CASCADE;foreignKey:TLEID;references:ID"`
}

// MapToContextTLEDomain converts a ContextTLE database model to a GameContextTLE domain model.
func MapToContextTLEDomain(ct ContextTLE) domain.GameContextTLE {
	return domain.GameContextTLE{
		ContextID: ct.ContextID,
		TLEID:     ct.TLEID,
		TLE:       MapToTLEDomain(ct.TLE),
	}
}

// MapToContextTLEModel converts a GameContextTLE domain model to a ContextTLE database model.
func MapToContextTLEModel(ct domain.GameContextTLE) ContextTLE {
	return ContextTLE{
		ContextID: ct.ContextID,
		TLEID:     ct.TLEID,
		TLE:       MapToTLEModel(ct.TLE),
	}
}

// ContextSatellite defines the many-to-many relationship between Context and Satellite.
type ContextSatellite struct {
	ContextID   string    `gorm:"not null;index"` // Foreign key to Context
	SatelliteID string    `gorm:"not null;index"` // Foreign key to Satellite
	Context     Context   `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
	Satellite   Satellite `gorm:"constraint:OnDelete:CASCADE;foreignKey:SatelliteID;references:ID"`
	LockedSince time.Time
	LockedBy    string
}

// MapToContextSatelliteDomain converts a ContextSatellite database model to a GameContextSatellite domain model.
func MapToContextSatelliteDomain(cs ContextSatellite) domain.GameContextSatellite {
	return domain.GameContextSatellite{
		ContextID:   cs.ContextID,
		SatelliteID: domain.SatelliteID(cs.SatelliteID),
		Satellite:   MapToSatelliteDomain(cs.Satellite),
		LockedSince: xtime.NewUtcTimeIgnoreZone(cs.LockedSince),
		LockedBy: cs.LockedBy,
	}
}

// MapToContextSatelliteModel converts a GameContextSatellite domain model to a ContextSatellite database model.
func MapToContextSatelliteModel(cs domain.GameContextSatellite) ContextSatellite {
	return ContextSatellite{
		ContextID:   cs.ContextID,
		SatelliteID: string(cs.SatelliteID),
		Satellite:   MapToSatelliteModel(cs.Satellite),
	}
}
