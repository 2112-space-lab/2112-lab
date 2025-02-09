package dbc

import (
	"github.com/org/2112-space-lab/org/app-service/internal/clients/dbc/adapters"
	"github.com/org/2112-space-lab/org/app-service/internal/config/constants"

	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

var dbClient *DBClient

func init() {
	dbClient = &DBClient{
		name:    constants.FEATURE_DATABASE,
		adapter: adapters.Adapters,
		silent:  true,
		gormConfig: &gorm.Config{
			Logger: gLogger.Default.LogMode(gLogger.Silent),
		},
	}
}

// GetDBClient definition
func GetDBClient() *DBClient {
	return dbClient
}
