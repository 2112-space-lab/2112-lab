package event_handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	event_builder "github.com/org/2112-space-lab/org/app-service/internal/events/builder"
	model "github.com/org/2112-space-lab/org/app-service/internal/graphql/models/generated"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	"github.com/org/2112-space-lab/org/app-service/internal/services"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	xtime "github.com/org/2112-space-lab/org/app-service/pkg/time"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// RehydrateGameContextHandler listens for REHYDRATE_GAME_CONTEXT events.
type RehydrateGameContextHandler struct {
	events.BaseHandler[model.RehydrateGameContextRequested]
	gameContextService services.ContextService
	globalRepo         repository.GlobalPropertyRepository
	eventEmitter       *events.EventEmitter
	tleRepo            repository.TleRepository
}

// NewRehydrateGameContextHandler creates a new instance of the handler.
func NewRehydrateGameContextHandler(
	gameContextService services.ContextService,
	eventEmitter *events.EventEmitter,
	globalRepo repository.GlobalPropertyRepository,
	tleRepo repository.TleRepository,
) *RehydrateGameContextHandler {
	return &RehydrateGameContextHandler{
		gameContextService: gameContextService,
		eventEmitter:       eventEmitter,
		globalRepo:         globalRepo,
		tleRepo:            tleRepo,
	}
}

func (h *RehydrateGameContextHandler) HandlerName() string {
	return "SatellitePositionHandler"
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
func (h *RehydrateGameContextHandler) HandleRehydrateGameContextEvent(ctx context.Context, event model.EventRoot, payload *model.RehydrateGameContextRequested) (err error) {
	_, span := tracing.NewSpan(ctx, "HandleRehydrateGameContextEvent")
	defer span.EndWithError(err)

	log.Infof("🔄 Processing rehydration request for game context: %s", payload.Name)

	var failureReason string
	var tleCount int32
	defer func() {
		if err != nil {
			ev, eventErr := event_builder.NewRehydrateGameContextFailedEvent(payload.Name, failureReason, tleCount)
			if eventErr == nil {
				_ = h.eventEmitter.PublishEvent(ctx, *ev)
			}
			log.Errorf("❌ Rehydration failed for context: %s | Reason: %s", payload.Name, failureReason)
		} else {
			ev, eventErr := event_builder.NewRehydrateGameContextSuccessEvent(payload.Name, tleCount)
			if eventErr == nil {
				_ = h.eventEmitter.PublishEvent(ctx, *ev)
			}
			log.Infof("✅ Rehydration process completed successfully for context: %s", payload.Name)
		}
	}()

	gameContext, err := h.gameContextService.GetByUniqueName(ctx, domain.GameContextName(payload.Name))
	if err != nil {
		failureReason = fmt.Sprintf("Failed to fetch game context: %v", err)
		log.Errorf("❌ %s", failureReason)
		return err
	}

	tles, err := h.tleRepo.GetTLEsByContextName(ctx, gameContext.Name)
	if err != nil {
		failureReason = fmt.Sprintf("Failed to retrieve TLEs for context %s: %v", gameContext.Name, err)
		log.Errorf("❌ %s", failureReason)
		return err
	}

	if len(tles) == 0 {
		log.Warnf("⚠️ No TLEs found for game context: %s", gameContext.Name)
	}

	tleCount = int32(len(tles))

	for _, tle := range tles {
		msg := fmt.Sprintf("🛰 Rehydrating TLE for SPACE ID %s", tle.SpaceID)

		propagationPayload := model.SatelliteTlePropagated{
			SpaceID:      tle.SpaceID,
			TleLine1:     tle.Line1,
			TleLine2:     tle.Line2,
			RedisKey:     fmt.Sprintf("tle:%s", tle.SpaceID),
			StartTimeUtc: xtime.UtcNow().Inner().Format(time.RFC3339),
		}

		ev, err := event_builder.NewSatelliteTlePropagationRequestedEvent(propagationPayload, &msg)
		if err != nil {
			failureReason = fmt.Sprintf("Failed to create propagation event for SPACE ID %s: %v", tle.SpaceID, err)
			log.Errorf("❌ %s", failureReason)
			return err
		}

		err = h.eventEmitter.PublishEvent(ctx, *ev)
		if err != nil {
			failureReason = fmt.Sprintf("Failed to publish propagation event for SPACE ID %s: %v", tle.SpaceID, err)
			log.Errorf("❌ %s", failureReason)
			return err
		}

		log.Infof("📤 Propagation event sent for SPACE ID %s", tle.SpaceID)
	}

	return nil
}
