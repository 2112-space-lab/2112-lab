package tasks

import (
	"context"
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/clients/redis"
	"github.com/Elbujito/2112/src/app-service/internal/domain"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	"github.com/Elbujito/2112/src/app-service/internal/tasks/handlers"
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
func NewTaskMonitor(satelliteRepo domain.SatelliteRepository, tleRepo repository.TleRepository, tileRepo domain.TileRepository, visibilityRepo domain.MappingRepository, tleService services.TleService, satelliteService services.SatelliteService, redisClient *redis.RedisClient) (TaskMonitor, error) {

	celestrackTleUpload := handlers.NewCelestrackTleUploadHandler(
		satelliteRepo,
		tleRepo,
		&tleService,
	)

	generateTilesHandler := handlers.NewGenerateTilesHandler(
		tileRepo,
	)

	mappingHandler := handlers.NewSatellitesTilesMappingsHandler(
		tileRepo,
		tleRepo,
		satelliteRepo,
		visibilityRepo,
		redisClient,
		4,
	)

	celestrackSatelliteUpload := handlers.NewCelesTrackSatelliteUploadHandler(
		satelliteRepo,
		&satelliteService,
	)

	satelliteVisibilities := handlers.NewComputeVisibilitiessHandler(
		tileRepo,
		visibilityRepo,
		tleRepo,
		redisClient,
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
