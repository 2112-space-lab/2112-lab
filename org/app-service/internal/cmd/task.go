package cmd

import (
	"github.com/Elbujito/2112/src/app-service/internal/app"
	"github.com/Elbujito/2112/src/app-service/internal/cmd/task"
	"github.com/Elbujito/2112/src/app-service/internal/proc"

	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/spf13/cobra"
)

// TaskCmd creates the `task` command with its subcommands
func TaskCmd(app *app.App) *cobra.Command {
	taskCmd := &cobra.Command{
		Use:   "task <option>",
		Short: "Start task",
		Long: `Run a one-time only task.
Please key in an option to start. Type 'task -h' for more information.

Popular options are:
- task list
- task init
- task cleanup

Use -d or --dev to start in development mode.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			execTaskPersistentPreRun(app)
		},
	}

	// Register subcommands dynamically
	taskCmd.AddCommand(task.ExecCmd(app))

	return taskCmd
}

// execTaskPersistentPreRun handles shared setup logic for all task subcommands
func execTaskPersistentPreRun(app *app.App) {
	logger.Debug("Executing task persistent pre run ...")

	// Initialize required services and dependencies
	proc.InitClients()
	proc.ConfigureClients()
	proc.InitDbConnection()
	proc.InitModels()

	// Add additional task-specific initializations if needed
}
