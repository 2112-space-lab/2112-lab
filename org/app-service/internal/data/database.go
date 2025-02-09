package data

import (
	"github.com/org/2112-space-lab/org/app-service/internal/clients/dbc"
	"gorm.io/gorm"
)

type Database struct {
	DbHandler *gorm.DB
}

func NewDatabase() Database {
	return Database{
		DbHandler: dbc.GetDBClient().DB,
	}
}
