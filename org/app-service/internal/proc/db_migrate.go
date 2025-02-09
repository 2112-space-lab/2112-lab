package proc

import (
	"io"
	"os"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/dbc"
	"github.com/org/2112-space-lab/org/app-service/internal/data/migrations"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

func DBMigrate() {

	var logWriter io.Writer
	logWriter = os.Stdout
	logger, err := log.NewLogger(logWriter, log.DebugLevel, log.LoggerTypes.Logrus())
	if err != nil {
		panic(err)
	}
	log.SetDefaultLogger(logger)
	dbClient := dbc.GetDBClient()

	dbClient.InitDBConnection()

	migrations.Init(dbClient.DB)

	if err := migrations.Migrate(); err != nil {
		logger.Errorf("Failed to apply migrations: %s", err)
	}

	logger.Info("Migrations applied successfully")

}
