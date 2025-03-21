package info

import (
	"github.com/org/2112-space-lab/org/app-service/internal/app"
	logger "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/spf13/cobra"
)

// EnvCmd creates the `env` subcommand
func EnvCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Print environment variables",
		Long:  "Display the environment variables for the service.",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Printing environment variables...")
			// Logic to display environment variables
		},
	}
}
