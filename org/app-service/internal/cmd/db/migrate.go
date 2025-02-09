package db

import (
	"github.com/org/2112-space-lab/org/app-service/internal/app"
	"github.com/org/2112-space-lab/org/app-service/internal/proc"
	"github.com/spf13/cobra"
)

// MigrateCmd creates the `migrate` command
func MigrateCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		Long:  "Run database migrations.",
		Run: func(cmd *cobra.Command, args []string) {
			proc.DBMigrate()
		},
	}
}
