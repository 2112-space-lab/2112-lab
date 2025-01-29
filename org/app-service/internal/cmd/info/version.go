package info

import (
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/app"
	"github.com/spf13/cobra"
)

// VersionCmd creates the `version` subcommand
func VersionCmd(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display service version",
		Long:  "Print the current version of the service.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("v%s\n", app.Version)
		},
	}
}
