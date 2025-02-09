package services

import (
	"log"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/celestrack"
	propagator "github.com/org/2112-space-lab/org/app-service/internal/clients/propagate"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/redis"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/data"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
)

// ServiceComponent holds all service instances for dependency injection.
type ServiceComponent struct {
	SatelliteService  SatelliteService
	TileService       TileService
	ContextService    ContextService
	AuditTrailService AuditTrailService
}

// NewServiceComponent initializes and returns a new ServiceComponent.
func NewServiceComponent(env *config.SEnv) *ServiceComponent {
	database := data.NewDatabase()

	redisClient, err := redis.NewRedisClient(config.Env)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	tleRepo := repository.NewTLERepository(&database, redisClient, 24*7*time.Hour)
	satelliteRepo := repository.NewSatelliteRepository(&database)
	tileRepo := repository.NewTileRepository(&database)
	mappingRepo := repository.NewTileSatelliteMappingRepository(&database)
	contextRepo := repository.NewContextRepository(&database)
	auditTrailRepo := repository.NewAuditTrailRepository(&database)

	propagteClient := propagator.NewPropagatorClient(env)
	celestrackClient := celestrack.NewCelestrackClient(env)

	satelliteService := NewSatelliteService(tleRepo, propagteClient, celestrackClient, satelliteRepo)
	tileService := NewTileService(tileRepo, tleRepo, satelliteRepo, mappingRepo)
	contextService := NewContextService(contextRepo)
	auditTrailService := NewAuditTrailService(auditTrailRepo)

	return &ServiceComponent{
		SatelliteService:  satelliteService,
		TileService:       tileService,
		ContextService:    contextService,
		AuditTrailService: auditTrailService,
	}
}
