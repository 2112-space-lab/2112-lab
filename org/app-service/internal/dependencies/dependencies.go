package dependencies

import (
	"context"

	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/data"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	event_handlers "github.com/org/2112-space-lab/org/app-service/internal/events/handlers"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
)

// Dependencies holds all dependencies in one place
type Dependencies struct {
	Clients      *Clients
	Repositories *Repositories
	Services     *Services
	EventLoop    *events.EventProcessor
	EventEmitter *events.EventEmitter
}

// NewDependencies initializes and returns a Dependencies struct
func NewDependencies(ctx context.Context, env *config.SEnv) (*Dependencies, error) {
	database := data.NewDatabase()
	clients := NewClients(env)

	repositories := NewRepositories(&database, clients)
	eventLoop := events.NewEventProcessor(repositories.EventRepo, repositories.EventHandlerRepo)
	eventEmitter, err := events.NewEventEmitter(ctx, clients.RabbitMQClient, eventLoop)
	if err != nil {
		return &Dependencies{}, err
	}
	services := NewServices(repositories, clients, eventEmitter)
	rehydrateGameContextHandler := event_handlers.NewRehydrateGameContextHandler(services.ContextService, eventEmitter, repositories.GlobalPropRepo)
	eventLoop.RegisterHandler(model.EventTypeRehydrateGameContext, rehydrateGameContextHandler)

	return &Dependencies{
		Clients:      clients,
		Repositories: repositories,
		Services:     services,
		EventLoop:    eventLoop,
		EventEmitter: eventEmitter,
	}, nil
}
