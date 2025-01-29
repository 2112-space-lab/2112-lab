package proc

import (
	"fmt"
	"io"
	"os"

	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
	log "github.com/Elbujito/2112/src/app-service/pkg/log"
)

func DBCreate() {

	var logWriter io.Writer
	logWriter = os.Stdout
	logger, err := log.NewLogger(logWriter, log.DebugLevel, log.LoggerTypes.Logrus())
	if err != nil {
		panic(err)
	}
	log.SetDefaultLogger(logger)
	dbClient := dbc.GetDBClient()

	dbClient.InitServerConnection()

	if err := dbClient.CreateDatabase(); err != nil {
		panic(fmt.Errorf("failed to create database: %w", err))
	}

	logger.Info("Database created successfully.")

}
