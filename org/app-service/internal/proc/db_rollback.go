package proc

import (
	"io"
	"os"

	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
	"github.com/Elbujito/2112/src/app-service/internal/data/migrations"
	log "github.com/Elbujito/2112/src/app-service/pkg/log"
)

func DBRollback() {

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

	if err := migrations.Rollback(); err != nil {
		logger.Errorf("Failed to rollback migrations: %s", err)
	}

}
