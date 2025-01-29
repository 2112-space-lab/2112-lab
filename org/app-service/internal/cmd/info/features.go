package info

import (
	"github.com/Elbujito/2112/src/app-service/internal/app"
	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/spf13/cobra"
)

// FeaturesCmd creates the `features` subcommand
func FeaturesCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "features",
		Short: "List enabled features",
		Long:  "Display the list of enabled features for the service.",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Listing enabled features...")
		},
	}
}
