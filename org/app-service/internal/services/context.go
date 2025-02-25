package services

import (
	"context"
	"time"

	"github.com/org/2112-space-lab/org/app-service/internal/domain"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	event_builder "github.com/org/2112-space-lab/org/app-service/internal/events/builder"
	repository "github.com/org/2112-space-lab/org/app-service/internal/repositories"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// ContextService definition
type ContextService struct {
	repo    repository.ContextRepository
	emitter *events.EventEmitter
}

// NewContextService creates a new instance of ContextService.
func NewContextService(repo repository.ContextRepository, emitter *events.EventEmitter) ContextService {
	return ContextService{repo: repo, emitter: emitter}
}

// Create creates a new GameContext.
func (c *ContextService) Create(ctx context.Context, context domain.GameContext) (cc domain.GameContext, err error) {
	ctx, span := tracing.NewSpan(ctx, "Create")
	defer span.EndWithError(err)
	err = c.repo.Save(ctx, context)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// Update updates an existing GameContext.
func (c *ContextService) Update(ctx context.Context, context domain.GameContext) (cc domain.GameContext, err error) {
	ctx, span := tracing.NewSpan(ctx, "Update")
	defer span.EndWithError(err)
	err = c.repo.Update(ctx, context)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// GetByUniqueName retrieves a GameContext by its unique name.
func (c *ContextService) GetByUniqueName(ctx context.Context, name domain.GameContextName) (cc domain.GameContext, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetByUniqueName")
	defer span.EndWithError(err)
	context, err := c.repo.FindByUniqueName(ctx, name)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// DeleteByUniqueName deletes a GameContext by its unique name.
func (c *ContextService) DeleteByUniqueName(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "DeleteByUniqueName")
	defer span.EndWithError(err)
	err = c.repo.DeleteByUniqueName(ctx, string(name))
	if err != nil {
		return err
	}
	return nil
}

// ActiveContext activates a GameContext by its unique name.
func (c *ContextService) ActiveContext(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "ActiveContext")
	defer span.EndWithError(err)
	err = c.repo.ActivateContext(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

// DisableContext deactivates a GameContext by its unique name.
func (c *ContextService) DisableContext(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "DisableContext")
	defer span.EndWithError(err)
	err = c.repo.DesactiveContext(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

// GetActiveContext retrieves the currently active GameContext.
func (c *ContextService) GetActiveContext(ctx context.Context) (cc domain.GameContext, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetActiveContext")
	defer span.EndWithError(err)
	context, err := c.repo.GetActiveContext(ctx)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// GetAllContexts retrieves all GameContexts.
func (c *ContextService) GetAllContexts(ctx context.Context) (cs []domain.GameContext, err error) {
	ctx, span := tracing.NewSpan(ctx, "GetAllContexts")
	defer span.EndWithError(err)
	contexts, err := c.repo.FindAll(ctx)
	if err != nil {
		return []domain.GameContext{}, err
	}
	return contexts, nil
}

// FindBySatelliteID retrieves all GameContexts associated with a specific satellite.
func (c *ContextService) FindBySatelliteID(ctx context.Context, satelliteID domain.SatelliteID) (cc domain.GameContext, err error) {
	ctx, span := tracing.NewSpan(ctx, "FindBySatelliteID")
	defer span.EndWithError(err)
	context, err := c.repo.FindActiveBySatelliteID(ctx, satelliteID)
	if err != nil {
		return domain.GameContext{}, err
	}
	return context, nil
}

// AssignSatellite associates a satellite with a GameContext.
func (c *ContextService) AssignSatellite(ctx context.Context, name domain.GameContextName, satelliteID domain.SatelliteID) (err error) {
	ctx, span := tracing.NewSpan(ctx, "AssignSatellite")
	defer span.EndWithError(err)
	err = c.repo.AssignSatellite(ctx, name, satelliteID)
	ctxLog := log.WithFields(log.Fields{"func": "AssignSatellite"})
	if err != nil {
		ctxLog.WithError(err).Error("failed to assign satellite")
		return err
	}
	return nil
}

// AssignSatellite associates a satellite with a GameContext.
func (c *ContextService) AssignSatellites(ctx context.Context, name domain.GameContextName, satelliteID []domain.SatelliteID) (err error) {
	ctx, span := tracing.NewSpan(ctx, "AssignSatellite")
	defer span.EndWithError(err)

	for _, id := range satelliteID {
		err = c.repo.AssignSatellite(ctx, name, id)
		ctxLog := log.WithFields(log.Fields{"func": "AssignSatellite"})
		if err != nil {
			ctxLog.WithError(err).Error("failed to assign satellite")
			return err
		}
	}
	return nil
}

// RemoveSatellite removes the association between a satellite and a GameContext.
func (c *ContextService) RemoveSatellite(ctx context.Context, name domain.GameContextName, satelliteID domain.SatelliteID) (err error) {
	ctx, span := tracing.NewSpan(ctx, "RemoveSatellite")
	defer span.EndWithError(err)
	err = c.repo.RemoveSatellite(ctx, name, satelliteID)
	if err != nil {
		return err
	}
	return nil
}

// FindAllWithPagination retrieves all contexts with pagination and optional filtering by name or description.
func (c *ContextService) FindAllWithPagination(ctx context.Context, page int, pageSize int, wildcard string) (cs []domain.GameContext, err error) {
	ctx, span := tracing.NewSpan(ctx, "FindAllWithPagination")
	defer span.EndWithError(err)
	contexts, err := c.repo.FindAllWithPagination(ctx, page, pageSize, wildcard)
	if err != nil {
		return []domain.GameContext{}, err
	}
	return contexts, nil
}

func (c *ContextService) SetActivatedAt(ctx context.Context, name domain.GameContextName, activatedAt time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx, "SetActivatedAt")
	defer span.EndWithError(err)
	return c.repo.SetActivatedAt(ctx, name, activatedAt)
}

func (c *ContextService) UnsetActivatedAt(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "UnsetActivatedAt")
	defer span.EndWithError(err)
	return c.repo.UnsetActivatedAt(ctx, name)
}

func (c *ContextService) SetDesactivatedAt(ctx context.Context, name domain.GameContextName, desactivatedAt time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx, "SetDesactivatedAt")
	defer span.EndWithError(err)
	return c.repo.SetDesactivatedAt(ctx, name, desactivatedAt)
}

func (c *ContextService) UnsetDesactivatedAt(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "UnsetDesactivatedAt")
	defer span.EndWithError(err)
	return c.repo.UnsetDesactivatedAt(ctx, name)
}

func (c *ContextService) SetTriggerGeneratedMappingAt(ctx context.Context, name domain.GameContextName, timestamp time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx, "SetTriggerGeneratedMappingAt")
	defer span.EndWithError(err)
	return c.repo.SetTriggerGeneratedMappingAt(ctx, name, timestamp)
}

func (c *ContextService) UnsetTriggerGeneratedMappingAt(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "UnsetTriggerGeneratedMappingAt")
	defer span.EndWithError(err)
	return c.repo.UnsetTriggerGeneratedMappingAt(ctx, name)
}

func (c *ContextService) SetTriggerImportedTLEAt(ctx context.Context, name domain.GameContextName, timestamp time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx, "SetTriggerImportedTLEAt")
	defer span.EndWithError(err)
	return c.repo.SetTriggerImportedTLEAt(ctx, name, timestamp)
}

func (c *ContextService) UnsetTriggerImportedTLEAt(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "UnsetTriggerImportedTLEAt")
	defer span.EndWithError(err)
	return c.repo.UnsetTriggerImportedTLEAt(ctx, name)
}

func (c *ContextService) SetTriggerImportedSatelliteAt(ctx context.Context, name domain.GameContextName, timestamp time.Time) (err error) {
	ctx, span := tracing.NewSpan(ctx, "UnsetTriggerImportedSatelliteAt")
	defer span.EndWithError(err)
	return c.repo.SetTriggerImportedSatelliteAt(ctx, name, timestamp)
}

func (c *ContextService) UnsetTriggerImportedSatelliteAt(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "UnsetTriggerImportedSatelliteAt")
	defer span.EndWithError(err)
	return c.repo.UnsetTriggerImportedSatelliteAt(ctx, name)
}

func (c *ContextService) Rehydrate(ctx context.Context, name domain.GameContextName) (err error) {
	ctx, span := tracing.NewSpan(ctx, "Rehydrate")
	defer span.EndWithError(err)

	ev, err := event_builder.NewRehydrateEvent(string(name), nil)
	if err != nil {
		return err
	}

	return c.emitter.PublishEvent(ctx, *ev)
}
