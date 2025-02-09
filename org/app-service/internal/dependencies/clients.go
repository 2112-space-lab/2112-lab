package dependencies

import (
	"log"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/celestrack"
	propagator "github.com/org/2112-space-lab/org/app-service/internal/clients/propagate"
	"github.com/org/2112-space-lab/org/app-service/internal/clients/redis"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
)

// Clients holds all client instances
type Clients struct {
	RedisClient      *redis.RedisClient
	PropagatorClient *propagator.PropagatorClient
	CelestrackClient *celestrack.CelestrackClient
}

// NewClients initializes and returns a Clients struct
func NewClients(env *config.SEnv) *Clients {
	redisClient, err := redis.NewRedisClient(env)
	if err != nil {
		log.Panic("Failed to initialize Redis client:", err)
	}

	return &Clients{
		RedisClient:      redisClient,
		PropagatorClient: propagator.NewPropagatorClient(env),
		CelestrackClient: celestrack.NewCelestrackClient(env),
	}
}

// Get retrieves a specific client and panics if it's not set
func (c *Clients) Get(client interface{}) interface{} {
	if client == nil {
		log.Panic("Requested client is not initialized")
	}
	return client
}
