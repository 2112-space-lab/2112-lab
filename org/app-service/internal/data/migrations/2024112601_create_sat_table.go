package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/org/2112-space-lab/org/app-service/internal/data/models"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "2025011102_init",
		Migrate: func(db *gorm.DB) error {
			// Ensure the schema exists
			if err := db.Exec("CREATE SCHEMA IF NOT EXISTS config_schema").Error; err != nil {
				return err
			}

			// Define tables within the schema

			type AuditTrail struct {
				models.ModelBase
				TableName   string    `gorm:"size:255;not null"`
				RecordID    string    `gorm:"size:255;not null"`
				Action      string    `gorm:"size:50;not null"`
				ChangesJSON string    `gorm:"type:json"`
				PerformedBy string    `gorm:"size:255;not null"`
				PerformedAt time.Time `gorm:"not null"`
			}

			type Context struct {
				models.ModelBase
				Name                       string     `gorm:"size:255;unique;not null"`
				TenantID                   string     `gorm:"size:255;not null;index"`
				Description                string     `gorm:"size:1024"`
				MaxSatellite               int        `gorm:"not null"`
				MaxTiles                   int        `gorm:"not null"`
				ActivatedAt                *time.Time `gorm:"null"`
				DesactivatedAt             *time.Time `gorm:"null"`
				TriggerGeneratedMappingAt  *time.Time `gorm:"null"`
				TriggerImportedTLEAt       *time.Time `gorm:"null"`
				TriggerImportedSatelliteAt *time.Time `gorm:"null"`
			}

			type Satellite struct {
				models.ModelBase
				Name           string     `gorm:"size:255;not null"`
				SpaceID        string     `gorm:"size:255;unique;not null"`
				Type           string     `gorm:"size:255"`
				LaunchDate     *time.Time `gorm:"type:date"`
				DecayDate      *time.Time `gorm:"type:date"`
				IntlDesignator string     `gorm:"size:255"`
				Owner          string     `gorm:"size:255"`
				ObjectType     string     `gorm:"size:255"`
				Period         *float64   `gorm:"type:float"`
				Inclination    *float64   `gorm:"type:float"`
				Apogee         *float64   `gorm:"type:float"`
				Perigee        *float64   `gorm:"type:float"`
				RCS            *float64   `gorm:"type:float"`
				Altitude       *float64   `gorm:"type:float"`
				OrbitType      string     `gorm:"size:255;not null"`
			}

			type TLE struct {
				models.ModelBase
				SpaceID string    `gorm:"not null;index"`
				Line1   string    `gorm:"size:255;not null"`
				Line2   string    `gorm:"size:255;not null"`
				Epoch   time.Time `gorm:"not null"`
			}

			type Tile struct {
				models.ModelBase
				Quadkey        string  `gorm:"size:256;unique;not null"`
				ZoomLevel      int     `gorm:"not null"`
				CenterLat      float64 `gorm:"not null"`
				CenterLon      float64 `gorm:"not null"`
				NbFaces        int     `gorm:"not null"`
				Radius         float64 `gorm:"not null"`
				BoundariesJSON string  `gorm:"type:json"`
				SpatialIndex   string  `gorm:"type:geometry(Polygon, 4326);spatialIndex"`
			}

			type TileSatelliteMapping struct {
				models.ModelBase
				SpaceID               string    `gorm:"not null;index"`
				TileID                string    `gorm:"not null;index"`
				TLEID                 string    `gorm:"not null;index"`
				ContextID             string    `gorm:"not null;index"`
				Context               Context   `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				Tile                  Tile      `gorm:"constraint:OnDelete:CASCADE;foreignKey:TileID;references:ID"`
				IntersectionLatitude  float64   `gorm:"type:double precision;not null;"`
				IntersectionLongitude float64   `gorm:"type:double precision;not null;"`
				IntersectedAt         time.Time `gorm:"not null"`
			}

			// Define the many-to-many relationship tables
			type ContextSatellite struct {
				ContextID   string    `gorm:"not null;index;uniqueIndex:unique_context_satellite"`
				SatelliteID string    `gorm:"not null;index;uniqueIndex:unique_context_satellite"`
				Context     Context   `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				Satellite   Satellite `gorm:"constraint:OnDelete:CASCADE;foreignKey:SatelliteID;references:ID"`
				LockedSince time.Time
				LockedBy    string
			}

			type ContextTLE struct {
				ContextID string  `gorm:"not null;index;uniqueIndex:unique_context_tle"`
				TLEID     string  `gorm:"not null;index;uniqueIndex:unique_context_tle"`
				Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				TLE       TLE     `gorm:"constraint:OnDelete:CASCADE;foreignKey:TLEID;references:ID"`
			}

			type ContextTile struct {
				ContextID string  `gorm:"not null;index;uniqueIndex:unique_context_tile"`
				TileID    string  `gorm:"not null;index;uniqueIndex:unique_context_tile"`
				Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				Tile      Tile    `gorm:"constraint:OnDelete:CASCADE;foreignKey:TileID;references:ID"`
			}

			type GlobalProperty struct {
				models.ModelBase
				Key         string `gorm:"primaryKey;size:255;not null" json:"key"`
				Value       string `gorm:"type:text;not null" json:"value"`
				Description string `gorm:"type:text" json:"description"`
				ValueType   string `gorm:"size:50;not null" json:"value_type"`
			}

			type Event struct {
				ID          string    `gorm:"primaryKey;size:255;not null"`
				EventType   string    `gorm:"size:255;not null;index"`
				EventUID    string    `gorm:"size:255;not null;unique"`
				Payload     string    `gorm:"type:json;not null"`
				PublishedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
				Comment     *string   `gorm:"type:text"`
			}

			type EventHandler struct {
				ID          string     `gorm:"primaryKey;size:255;not null"`
				EventID     string     `gorm:"size:255;not null;index"`
				Handler     string     `gorm:"size:255;not null"`
				StartedAt   time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
				CompletedAt *time.Time `gorm:"null"`
				Status      string     `gorm:"size:50;not null"`
				Error       *string    `gorm:"type:text"`

				Event Event `gorm:"constraint:OnDelete:CASCADE;foreignKey:EventID;references:ID"`
			}

			// AutoMigrate with schema
			if err := db.Set("gorm:table_options", "SCHEMA=config_schema").
				AutoMigrate(
					&Context{},
					&Satellite{},
					&TLE{},
					&Tile{},
					&TileSatelliteMapping{},
					&ContextSatellite{},
					&ContextTLE{},
					&ContextTile{},
					&AuditTrail{},
					&GlobalProperty{},
					&Event{},
					&EventHandler{},
				); err != nil {
				return err
			}

			// Unique constraint for contexts
			if err := db.Exec(`
				ALTER TABLE config_schema.contexts
				ADD CONSTRAINT unique_tenant_context_name
				UNIQUE (tenant_id, name);
			`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(db *gorm.DB) error {
			// Drop constraint first
			if err := db.Exec(`
				ALTER TABLE config_schema.contexts
				DROP CONSTRAINT IF EXISTS unique_tenant_context_name;
			`).Error; err != nil {
				return err
			}

			// Drop tables in correct order to avoid dependency errors
			return db.Migrator().DropTable(
				"config_schema.context_satellites",
				"config_schema.context_tles",
				"config_schema.context_tiles",
				"config_schema.tile_satellite_mappings",
				"config_schema.tiles",
				"config_schema.tles",
				"config_schema.satellites",
				"config_schema.contexts",
				"config_schema.audit_trails",
				"config_schema.events",
				"config_schema.event_handlers",
			)
		},
	}

	AddMigration(m)
}
