package proc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/celestrack"
	propagator "github.com/Elbujito/2112/src/app-service/internal/clients/propagate"
	"github.com/Elbujito/2112/src/app-service/internal/clients/redis"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/data"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	"github.com/Elbujito/2112/src/app-service/internal/tasks"
	"github.com/Elbujito/2112/src/app-service/internal/tasks/handlers"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xutils"
)

func TaskExec(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a task name")
		return
	}
	taskName := args[0]
	taskArgs := xutils.ResolveArgs(args[1:])

	database := data.NewDatabase()

	propagteClient := propagator.NewPropagatorClient(config.Env)
	redisClient, err := redis.NewRedisClient(config.Env)
	if err != nil {
		log.Println(err.Error())
		return
	}

	tleRepo := repository.NewTLERepository(&database, redisClient, 3600*time.Hour)
	celestrackClient := celestrack.NewCelestrackClient(config.Env)
	satelliteRepo := repository.NewSatelliteRepository(&database)
	visibilityRepo := repository.NewTileSatelliteMappingRepository(&database)
	tileRepo := repository.NewTileRepository(&database)
	contextRepo := repository.NewContextRepository(&database)

	tleService := services.NewTleService(celestrackClient, tleRepo, contextRepo)
	satService := services.NewSatelliteService(tleRepo, propagteClient, celestrackClient, satelliteRepo)

	monitor, err := tasks.NewTaskMonitor(satelliteRepo, tleRepo, tileRepo, visibilityRepo, tleService, satService, redisClient)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = monitor.Process(ctx, handlers.TaskName(taskName), taskArgs)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
