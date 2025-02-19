package tasks

import (
	"context"
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/dependencies"
	"github.com/org/2112-space-lab/org/app-service/internal/events"
	"github.com/org/2112-space-lab/org/app-service/internal/tasks/handlers"
	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
)

// TaskHandler definition
type TaskHandler interface {
	GetTask() handlers.Task
	Run(ctx context.Context, args map[string]string) error
}

// TaskMonitor definition
type TaskMonitor struct {
	Tasks map[handlers.TaskName]TaskHandler
}

// TaskMonitor constructor
func NewTaskMonitor(ctx context.Context, dependencies *dependencies.Dependencies) (t TaskMonitor, err error) {
	ctx, span := tracing.NewSpan(ctx, "TaskMonitor.NewTaskMonitor")
	defer span.EndWithError(err)

	eventMonitor, err := events.NewEventMonitor(ctx, dependencies.Clients.RabbitMQClient)
	if err != nil {
		return t, err
	}

	eventEmitter, err := events.NewEventEmitter(ctx, dependencies.Clients.RabbitMQClient)
	if err != nil {
		return t, err
	}

	celestrackTleUpload := handlers.NewCelestrackTleUploadHandler(
		dependencies.Repositories.SatelliteRepo,
		dependencies.Repositories.TleRepo,
		&dependencies.Services.TleService,
		eventEmitter,
		eventMonitor,
	)

	generateTilesHandler := handlers.NewGenerateTilesHandler(
		&dependencies.Repositories.TileRepo,
	)

	mappingHandler := handlers.NewSatellitesTilesMappingsHandler(
		&dependencies.Repositories.TileRepo,
		dependencies.Repositories.TleRepo,
		&dependencies.Repositories.SatelliteRepo,
		&dependencies.Repositories.MappingRepo,
		dependencies.Clients.RedisClient,
	)

	celestrackSatelliteUpload := handlers.NewCelesTrackSatelliteUploadHandler(
		&dependencies.Repositories.SatelliteRepo,
		&dependencies.Services.SatelliteService,
	)

	satelliteVisibilities := handlers.NewComputeVisibilitiessHandler(
		&dependencies.Repositories.TileRepo,
		&dependencies.Repositories.MappingRepo,
		dependencies.Repositories.TleRepo,
		dependencies.Clients.RedisClient,
	)

	eventDetector, err := handlers.NewEventDetector(
		ctx, eventEmitter, eventMonitor, dependencies)
	if err != nil {
		return t, err
	}

	tasks := map[handlers.TaskName]TaskHandler{
		celestrackTleUpload.GetTask().Name:       &celestrackTleUpload,
		generateTilesHandler.GetTask().Name:      &generateTilesHandler,
		mappingHandler.GetTask().Name:            &mappingHandler,
		celestrackSatelliteUpload.GetTask().Name: &celestrackSatelliteUpload,
		satelliteVisibilities.GetTask().Name:     &satelliteVisibilities,
		eventDetector.GetTask().Name:             &eventDetector,
	}
	return TaskMonitor{
		Tasks: tasks,
	}, err
}

// Process execute processor
func (t *TaskMonitor) Process(ctx context.Context, taskName handlers.TaskName, args map[string]string) (err error) {
	ctx, span := tracing.NewSpan(ctx, "TaskMonitor.Process")
	defer span.EndWithError(err)

	handler, err := t.GetMatchingTask(ctx, taskName)
	if err != nil {
		return err
	}
	return handler.Run(ctx, args)
}

// GetMatchingTask finds matching task
func (t *TaskMonitor) GetMatchingTask(ctx context.Context, taskName handlers.TaskName) (task TaskHandler, err error) {
	_, span := tracing.NewSpan(ctx, "TaskMonitor.GetMatchingTask")
	defer span.EndWithError(err)

	hh, ok := t.Tasks[taskName]
	if !ok {
		return task, fmt.Errorf("task no found for [%s]", taskName)
	}
	return hh, nil
}

// RunTaskAsGoroutine runs as go routine
func (t *TaskMonitor) RunTaskAsGoroutine(ctx context.Context, taskName handlers.TaskName, args map[string]string) (err error) {
	ctx, span := tracing.NewSpan(ctx, "TaskMonitor.RunTaskAsGoroutine")
	defer span.EndWithError(err)

	handler, err := t.GetMatchingTask(ctx, taskName)
	if err != nil {
		return err
	}
	go func() {
		if runErr := handler.Run(ctx, args); runErr != nil {
			fmt.Printf("Error running task [%s]: %v\n", taskName, runErr)
		}
	}()
	return nil
}
