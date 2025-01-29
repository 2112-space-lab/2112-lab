package db

import (
	"github.com/Elbujito/2112/src/app-service/internal/app"
	"github.com/Elbujito/2112/src/app-service/internal/proc"
	"github.com/spf13/cobra"
)

// SeedCmd creates the `seed` command
func SeedCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed database",
		Long:  "Backfill database with seed data.",
		Run: func(cmd *cobra.Command, args []string) {
			proc.DBSeed()
		},
	}
}
