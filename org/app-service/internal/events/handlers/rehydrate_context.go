package event_handlers

import (
	"context"

	"github.com/org/2112-space-lab/org/app-service/internal/events"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// RehydrateGameContextHandler listens for REHYDRATE_GAME_CONTEXT events.
type RehydrateGameContextHandler struct {
	events.BaseHandler[model.RehydrateGameContext]
	gameContextService services.ContextService
	globalRepo         repository.GlobalPropertyRepository
	eventEmitter       *events.EventEmitter
}

// NewRehydrateGameContextHandler creates a new instance of the handler.
func NewRehydrateGameContextHandler(
	gameContextService services.ContextService,
	eventEmitter *events.EventEmitter,
	globalRepo repository.GlobalPropertyRepository,
) *RehydrateGameContextHandler {
	return &RehydrateGameContextHandler{
		gameContextService: gameContextService,
		eventEmitter:       eventEmitter,
		globalRepo:         globalRepo,
	}
}

// Run processes the REHYDRATE_GAME_CONTEXT event.
func (h *RehydrateGameContextHandler) Run(ctx context.Context, event model.EventRoot) (err error) {
	ctx, span := tracing.NewSpan(ctx, "RunRehydrateEvent")
	defer span.EndWithError(err)

	log.Infof("🔄 Processing RehydrateGameContext event: UID=%s", event.EventUID)

	payload, err := h.Parse(event.Payload)
	if err != nil {
		log.Errorf("❌ Failed to parse payload for RehydrateGameContext: %v", err)
		return err
	}

	return h.HandleRehydrateGameContextEvent(ctx, event, payload)
}

// HandleRehydrateGameContextEvent rehydrates the game context and updates the repository.
func (h *RehydrateGameContextHandler) HandleRehydrateGameContextEvent(ctx context.Context, event model.EventRoot, payload *model.RehydrateGameContext) (err error) {
	_, span := tracing.NewSpan(ctx, "HandleRehydrateGameContextEvent")
	defer span.EndWithError(err)

	return nil
}
