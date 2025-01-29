package db

import (
	"github.com/Elbujito/2112/src/app-service/internal/app"
	"github.com/Elbujito/2112/src/app-service/internal/proc"
	"github.com/spf13/cobra"
)

// RollbackCmd creates the `rollback` command
func RollbackCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback database",
		Long:  "Rollback one database migration.",
		Run: func(cmd *cobra.Command, args []string) {
			proc.DBRollback()
		},
	}
}
