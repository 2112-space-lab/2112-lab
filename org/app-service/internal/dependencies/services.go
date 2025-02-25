package dependencies

import (
	"log"

	"github.com/org/2112-space-lab/org/app-service/internal/events"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
)

// Services holds all service instances
type Services struct {
	SatelliteService  services.SatelliteService
	TileService       services.TileService
	ContextService    services.ContextService
	AuditTrailService services.AuditTrailService
	TleService        services.TleService
}

// NewServices initializes and returns a Services struct
func NewServices(repos *Repositories, clients *Clients, emitter *events.EventEmitter) *Services {
	return &Services{
		SatelliteService:  services.NewSatelliteService(repos.TleRepo, clients.PropagatorClient, clients.CelestrackClient, repos.SatelliteRepo),
		TileService:       services.NewTileService(repos.TileRepo, repos.TleRepo, repos.SatelliteRepo, repos.MappingRepo),
		ContextService:    services.NewContextService(repos.ContextRepo, emitter),
		AuditTrailService: services.NewAuditTrailService(repos.AuditRepo),
		TleService:        services.NewTleService(clients.CelestrackClient, repos.TleRepo, &repos.ContextRepo),
	}
}

// Get retrieves a specific service and panics if it's not set
func (s *Services) Get(service interface{}) interface{} {
	if service == nil {
		log.Panic("Requested service is not initialized")
	}
	return service
}
