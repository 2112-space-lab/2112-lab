package proc

import (
	"io"
	"os"

	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
	"github.com/Elbujito/2112/src/app-service/internal/data/migrations"
	log "github.com/Elbujito/2112/src/app-service/pkg/log"
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
