package domain

import (
	"context"
	"time"

	fx "github.com/Elbujito/2112/src/app-service/pkg/option"
	xtime "github.com/Elbujito/2112/src/app-service/pkg/time"
)

type GameContextName string
type GameContextDescription string
type TenantID string

// GameContext represents a logical grouping for satellites, TLEs, and mappings in the domain layer.
type GameContext struct {
	ModelBase
	Name                       GameContextName
	TenantID                   TenantID
	Description                fx.Option[GameContextDescription]
	ActivatedAt                fx.Option[xtime.UtcTime]
	DesactivatedAt             fx.Option[xtime.UtcTime]
	TriggerGeneratedMappingAt  fx.Option[xtime.UtcTime]
	TriggerImportedTLEAt       fx.Option[xtime.UtcTime]
	TriggerImportedSatelliteAt fx.Option[xtime.UtcTime]
}

// GameContextSatellite represents the relationship between Context and Satellite in the domain layer.
type GameContextSatellite struct {
	ContextID   string      // ID of the Context
	SatelliteID SatelliteID // ID of the Satellite
	Satellite   Satellite
}

// GameContextTLE represents the relationship between Context and TLE in the domain layer.
type GameContextTLE struct {
	ContextID string // ID of the Context
	TLEID     string // ID of the TLE
	TLE       TLE
}

// GameContextTile represents the relationship between a Context and a Tile in the domain layer.
type GameContextTile struct {
	ContextID string // ID of the associated Context
	TileID    string // ID of the associated Tile
	Tile      Tile
}

// GameContextRepository defines the interface for context operations.
type GameContextRepository interface {
	// Core CRUD operations
	Save(ctx context.Context, gameContext GameContext) error
	Update(ctx context.Context, gameContext GameContext) error
	FindByUniqueName(ctx context.Context, gameContextName GameContextName) (GameContext, error)
	FindAll(ctx context.Context) ([]GameContext, error)
	DeleteByUniqueName(ctx context.Context, id string) error

	// Satellite association management
	FindActiveBySatelliteID(ctx context.Context, satelliteID SatelliteID) (GameContext, error)
	AssignSatellite(ctx context.Context, gameContextName GameContextName, satelliteID SatelliteID) error
	RemoveSatellite(ctx context.Context, gameContextName GameContextName, satelliteID SatelliteID) error

	// Context activation and deactivation
	DesactiveContext(ctx context.Context, gameContextName GameContextName) error
	ActivateContext(ctx context.Context, gameContextName GameContextName) error

	// Retrieve active contexts
	GetActiveContext(ctx context.Context) (GameContext, error)

	// Pagination with filtering
	FindAllWithPagination(ctx context.Context, page int, pageSize int, wildcard string) ([]GameContext, error)

	// Setter and Unsetter for timestamps
	SetActivatedAt(ctx context.Context, gameContextName GameContextName, activatedAt time.Time) error
	UnsetActivatedAt(ctx context.Context, gameContextName GameContextName) error
	SetDesactivatedAt(ctx context.Context, gameContextName GameContextName, desactivatedAt time.Time) error
	UnsetDesactivatedAt(ctx context.Context, gameContextName GameContextName) error
	SetTriggerGeneratedMappingAt(ctx context.Context, gameContextName GameContextName, timestamp time.Time) error
	UnsetTriggerGeneratedMappingAt(ctx context.Context, gameContextName GameContextName) error
	SetTriggerImportedTLEAt(ctx context.Context, gameContextName GameContextName, timestamp time.Time) error
	UnsetTriggerImportedTLEAt(ctx context.Context, gameContextName GameContextName) error
	SetTriggerImportedSatelliteAt(ctx context.Context, gameContextName GameContextName, timestamp time.Time) error
	UnsetTriggerImportedSatelliteAt(ctx context.Context, gameContextName GameContextName) error
}
