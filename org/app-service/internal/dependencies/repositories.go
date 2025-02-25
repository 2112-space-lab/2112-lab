package dependencies

import (
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/data"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
)

// Repositories holds all repository instances
type Repositories struct {
	TleRepo          repository.TleRepository
	SatelliteRepo    repository.SatelliteRepository
	TileRepo         repository.TileRepository
	MappingRepo      repository.TileSatelliteMappingRepository
	ContextRepo      repository.ContextRepository
	AuditRepo        repository.AuditTrailRepository
	GlobalPropRepo   repository.GlobalPropertyRepository
	EventRepo        repository.EventRepository
	EventHandlerRepo repository.EventHandlerRepository
}

// NewRepositories initializes and returns a Repositories struct
func NewRepositories(db *data.Database, clients *Clients) *Repositories {
	return &Repositories{
		TleRepo:          repository.NewTLERepository(db, clients.RedisClient),
		SatelliteRepo:    repository.NewSatelliteRepository(db, clients.RedisClient, time.Hour*24),
		TileRepo:         repository.NewTileRepository(db),
		MappingRepo:      repository.NewTileSatelliteMappingRepository(db),
		ContextRepo:      repository.NewContextRepository(db),
		AuditRepo:        repository.NewAuditTrailRepository(db),
		GlobalPropRepo:   repository.NewGlobalPropertyRepository(db),
		EventRepo:        repository.NewEventRepository(db),
		EventHandlerRepo: repository.NewEventHandlerRepository(db),
	}
}

// Get retrieves a specific repository and panics if it's not set
func (r *Repositories) Get(repo interface{}) interface{} {
	if repo == nil {
		panic("Requested repository is not initialized")
	}
	return repo
}
