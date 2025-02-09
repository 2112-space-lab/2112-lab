package cmd

import (
	"context"

	"github.com/org/2112-space-lab/org/app-service/internal/app"

	logger "github.com/org/2112-space-lab/org/app-service/pkg/log"
	"github.com/spf13/cobra"
)

// Global variable to hold the application version
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "CLI for managing services",
	Long:  "CLI for managing services, databases, and tasks.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug("Root command persistent pre-run executed.")
	},
}

// Execute runs the root command and its subcommands
func Execute(ctx context.Context) {
	app, err := app.NewApp(ctx, rootCmd.Use, Version)
	if err != nil {
		logger.Errorf("Command execution failed: %v", err)
		return
	}

	// Register subcommands with the App instance
	registerSubcommands(&app)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("Command execution failed: %v", err)
	}
}

// registerSubcommands dynamically registers all subcommands
func registerSubcommands(app *app.App) {
	rootCmd.AddCommand(StartCmd(app))
	rootCmd.AddCommand(DbCmd(app))
	rootCmd.AddCommand(InfoCmd(app))
	rootCmd.AddCommand(TaskCmd(app))
}
