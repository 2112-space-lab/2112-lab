package data

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
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
