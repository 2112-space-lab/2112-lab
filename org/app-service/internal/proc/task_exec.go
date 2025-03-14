package proc

import (
	"context"
	"fmt"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/dependencies"
	"github.com/org/2112-space-lab/org/app-service/internal/tasks"
	"github.com/org/2112-space-lab/org/app-service/internal/tasks/handlers"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xutils"
)

func TaskExec(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a task name")
		return
	}
	taskName := args[0]
	taskArgs := xutils.ResolveArgs(args[1:])

	deps, err := dependencies.NewDependencies(ctx, config.Env)
	if err != nil {
		log.Error(err.Error())
		return
	}
	monitor, err := tasks.NewTaskMonitor(ctx, deps)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = monitor.Process(ctx, handlers.TaskName(taskName), taskArgs)
	if err != nil {
		log.Error(err.Error())
		return
	}

	go deps.EventLoop.StartProcessing(ctx)
}
