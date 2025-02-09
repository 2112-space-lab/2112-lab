package db

import (
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/app"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	"github.com/org/2112-space-lab/org/app-service/internal/proc"
	"github.com/spf13/cobra"
)

// DropCmd creates the `drop` command
func DropCmd(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drop",
		Short: "Drop database",
		Long:  "Drop the database.",
		Run: func(cmd *cobra.Command, args []string) {
			if !config.ConfirmFlag {
				fmt.Println("This is a destructive action and it is irreversible.")
				fmt.Println("To delete, please run again using the `--confirm` flag.")
				return
			}
			proc.DBDrop()
		},
	}

	// Add flags
	cmd.PersistentFlags().BoolVar(&config.ConfirmFlag, "confirm", false, "Confirm deletion of database")

	return cmd
}
