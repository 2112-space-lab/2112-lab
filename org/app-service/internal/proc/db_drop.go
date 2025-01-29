package proc

import (
	"fmt"
	"io"
	"os"

	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	log "github.com/Elbujito/2112/src/app-service/pkg/log"
)

func DBDrop() {
	var logWriter io.Writer
	logWriter = os.Stdout
	logger, err := log.NewLogger(logWriter, log.DebugLevel, log.LoggerTypes.Logrus())
	if err != nil {
		panic(err)
	}
	log.SetDefaultLogger(logger)
	dbClient := dbc.GetDBClient()

	dbClient.InitServerConnection()

	if err := dbClient.DropDatabase(); err != nil {
		panic(fmt.Errorf("failed to drop database: %w", err))
	}

	logger.Infof("Database '" + config.Env.ServiceName + "' dropped.")

}
