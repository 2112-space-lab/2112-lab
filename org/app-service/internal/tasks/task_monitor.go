package tasks

import (
	"context"
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/dependencies"
	"github.com/org/2112-space-lab/org/app-service/internal/tasks/handlers"
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
func NewTaskMonitor(dependencies *dependencies.Dependencies) (TaskMonitor, error) {

	celestrackTleUpload := handlers.NewCelestrackTleUploadHandler(
		dependencies.Repositories.SatelliteRepo,
		dependencies.Repositories.TleRepo,
		&dependencies.Services.TleService,
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

	tasks := map[handlers.TaskName]TaskHandler{
		celestrackTleUpload.GetTask().Name:       &celestrackTleUpload,
		generateTilesHandler.GetTask().Name:      &generateTilesHandler,
		mappingHandler.GetTask().Name:            &mappingHandler,
		celestrackSatelliteUpload.GetTask().Name: &celestrackSatelliteUpload,
		satelliteVisibilities.GetTask().Name:     &satelliteVisibilities,
	}
	return TaskMonitor{
		Tasks: tasks,
	}, nil
}

// Process execute processor
func (t *TaskMonitor) Process(ctx context.Context, taskName handlers.TaskName, args map[string]string) error {
	handler, err := t.GetMatchingTask(taskName)
	if err != nil {
		return err
	}
	return handler.Run(ctx, args)
}

// GetMatchingTask finds matching task
func (t *TaskMonitor) GetMatchingTask(taskName handlers.TaskName) (task TaskHandler, err error) {
	hh, ok := t.Tasks[taskName]
	if !ok {
		return task, fmt.Errorf("task no found for [%s]", taskName)
	}
	return hh, nil
}

// RunTaskAsGoroutine runs as go routine
func (t *TaskMonitor) RunTaskAsGoroutine(ctx context.Context, taskName handlers.TaskName, args map[string]string) error {
	handler, err := t.GetMatchingTask(taskName)
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
