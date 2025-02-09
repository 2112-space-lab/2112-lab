package db

import (
	"github.com/org/2112-space-lab/org/app-service/internal/app"
	"github.com/org/2112-space-lab/org/app-service/internal/proc"
	"github.com/spf13/cobra"
)

// CreateCmd creates the `create` command
func CreateCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create database",
		Long:  "Create the database.",
		Run: func(cmd *cobra.Command, args []string) {
			proc.DBCreate()
		},
	}
}
