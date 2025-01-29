package cmd

import (
	"github.com/Elbujito/2112/src/app-service/internal/app"
	"github.com/Elbujito/2112/src/app-service/internal/cmd/info"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/spf13/cobra"
)

// InfoCmd creates the `info` command
func InfoCmd(app *app.App) *cobra.Command {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Print service env and config info",
		Long: `Print information related to the service environment and feature configuration. 
This command is a helper to get you started in your debugging journey.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			execInfoPersistentPreRun(app)
		},
	}

	// Add flags
	infoCmd.PersistentFlags().BoolVarP(&config.NoBorderFlag, "no-border", "N", false, "Print tables without border")

	// Add subcommands
	infoCmd.AddCommand(info.EnvCmd(app))
	infoCmd.AddCommand(info.FeaturesCmd(app))
	infoCmd.AddCommand(info.VersionCmd(app))

	return infoCmd
}

// execInfoPersistentPreRun handles shared setup logic before running any info subcommand
func execInfoPersistentPreRun(app *app.App) {
	logger.Debug("Executing info persistent pre run ...")
}
