package dependencies

import (
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/data"
)

// Dependencies holds all dependencies in one place
type Dependencies struct {
	Clients      *Clients
	Repositories *Repositories
	Services     *Services
}

// NewDependencies initializes and returns a Dependencies struct
func NewDependencies(env *config.SEnv) *Dependencies {
	// Initialize Database
	database := data.NewDatabase()

	// Initialize Clients
	clients := NewClients(env)

	// Initialize Repositories
	repositories := NewRepositories(&database, clients)

	// Initialize Services
	services := NewServices(repositories, clients)

	return &Dependencies{
		Clients:      clients,
		Repositories: repositories,
		Services:     services,
	}
}
