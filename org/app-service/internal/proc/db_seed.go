package proc

import (
	"io"
	"os"

	"github.com/org/2112-space-lab/org/app-service/internal/clients/dbc"
	"github.com/org/2112-space-lab/org/app-service/internal/data/seeds"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

func DBSeed() {

	var logWriter io.Writer
	logWriter = os.Stdout
	logger, err := log.NewLogger(logWriter, log.DebugLevel, log.LoggerTypes.Logrus())
	if err != nil {
		panic(err)
	}
	log.SetDefaultLogger(logger)
	dbClient := dbc.GetDBClient()

	dbClient.InitDBConnection()

	seeds.Init(dbClient.DB)

	if err := seeds.Apply(); err != nil {
		logger.Errorf("Failed to apply seeds: %s", err)
	}

	logger.Info("Seeds applied successfully")

}
