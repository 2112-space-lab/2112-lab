package cmd

import (
	"github.com/Elbujito/2112/src/app-service/internal/app"
	"github.com/Elbujito/2112/src/app-service/internal/cmd/db"
	"github.com/Elbujito/2112/src/app-service/internal/proc"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/spf13/cobra"
)

// DbCmd creates the `db` command with its subcommands
func DbCmd(app *app.App) *cobra.Command {
	dbCmd := &cobra.Command{
		Use:   "db <option>",
		Short: "Start db-related operations",
		Long: `Start a database operation.
Please key in an option to start. Type 'db -h' for more information.

Popular options are:
- db migrate
- db rollback`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			rootCmd.PersistentPreRun(cmd, args)
			execDbPersistentPreRun(app, cmd)
		},
	}

	// Register subcommands
	dbCmd.AddCommand(db.CreateCmd(app))
	dbCmd.AddCommand(db.DropCmd(app))
	dbCmd.AddCommand(db.MigrateCmd(app))
	dbCmd.AddCommand(db.RollbackCmd(app))
	dbCmd.AddCommand(db.SeedCmd(app))

	return dbCmd
}

// execDbPersistentPreRun handles shared setup logic before any db subcommand
func execDbPersistentPreRun(app *app.App, cmd *cobra.Command) {
	logger.Debug("Executing db persistent pre run ...")

	proc.InitDbClient()
	proc.ConfigureClients()

	ca := cmd.CalledAs()

	// Establish connection only for specific commands
	switch ca {
	case xconstants.NAME_CMD_DB_MIGRATE,
		xconstants.NAME_CMD_DB_ROLLBACK,
		xconstants.NAME_CMD_DB_SEED:
		proc.InitDbConnection()
		proc.InitModels()
	}

	// Additional global initializations can be added here if needed
}
