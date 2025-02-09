package cmd

import (
	"github.com/org/2112-space-lab/org/app-service/internal/app"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/proc"
	logger "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/spf13/cobra"
)

// StartCmd creates the `start` command
func StartCmd(app *app.App) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start public and protected API services",
		Long:  "Start both public and protected API services, as well as optional daemons.",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Starting public and protected API services...")
			proc.StartPublicApi(app.Services)
			proc.StartProtectedApi(app.Services)
		},
	}

	// Add subcommands
	startCmd.AddCommand(PublicApiCmd(app))
	startCmd.AddCommand(ProtectedApiCmd(app))
	startCmd.AddCommand(WatcherCmd(app))

	// Set global flags
	startCmd.PersistentFlags().BoolVar(&config.StartWatcherFlag, "watcher", false, "Start watcher daemon in background")
	startCmd.PersistentFlags().StringVarP(&config.HostFlag, "host", "H", "", "Service host")
	startCmd.PersistentFlags().StringVar(&config.ProtectedPortFlag, "protected-api-port", "", "Protected API Service port")
	startCmd.PersistentFlags().StringVar(&config.PublicPortFlag, "public-api-port", "", "Public API Service port")

	return startCmd
}

// PublicApiCmd creates the `publicApi` subcommand
func PublicApiCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "publicApi",
		Short: "Start public API service",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Starting public API service...")
			proc.StartPublicApi(app.Services)
		},
	}
}

// ProtectedApiCmd creates the `protectedApi` subcommand
func ProtectedApiCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "protectedApi",
		Short: "Start protected API service",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Starting protected API service...")
			proc.StartProtectedApi(app.Services)
		},
	}
}

// WatcherCmd creates the `watcher` subcommand
func WatcherCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "watcher",
		Short: "Start the watcher daemon",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("Starting watcher daemon...")
		},
	}
}
