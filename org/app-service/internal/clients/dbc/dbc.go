package dbc

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc/adapters"
	"github.com/Elbujito/2112/src/app-service/internal/config/features"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xutils"

	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

// DBClient definition
type DBClient struct {
	name       string
	config     features.DatabaseConfig
	silent     bool
	DB         *gorm.DB
	gormConfig *gorm.Config
	adapter    *adapters.Adapter
	driver     gorm.Dialector
	dsn        string
}

// Name definition
func (dbc *DBClient) Name() string {
	return dbc.name
}

// SetSilent definition
func (dbc *DBClient) SetSilent(val bool) {
	dbc.silent = val
	if !dbc.IsSilent() {
		dbc.gormConfig = &gorm.Config{
			Logger: gLogger.Default.LogMode(gLogger.Info),
		}
	}
}

// IsSilent definition
func (dbc *DBClient) IsSilent() bool {
	return dbc.silent
}

// clean up all this
// InitServerConnection definition
func (dbc *DBClient) InitServerConnection() {
	dbc.ResolveServerDriver()
	dbc.Connect()
}

// InitDBConnection definition
func (dbc *DBClient) InitDBConnection() {
	dbc.ResolveDriver()
	dbc.Connect()
}

// Configure definition
func (dbc *DBClient) Configure(v any) {
	dbc.config = v.(reflect.Value).Interface().(features.DatabaseConfig)
	dbc.adapter.SetConfig(dbc.config)
	if err := dbc.adapter.ValidateConfig(); err != nil {
		logger.Error(err.Error())
		xutils.Exit(1)
	}
}

// ResolveDriver definition
func (dbc *DBClient) ResolveDriver() {
	dsn, err := dbc.adapter.GetDSN()
	if err != nil {
		logger.Error(err.Error())
		xutils.Exit(1)
	}
	dbc.dsn = dsn

	driver, err := dbc.adapter.GetDriver()
	if err != nil {
		logger.Error(err.Error())
		xutils.Exit(1)
	}
	dbc.driver = driver
}

// ResolveServerDriver definition
func (dbc *DBClient) ResolveServerDriver() {
	dsn, err := dbc.adapter.GetServerDSN()
	if err != nil {
		logger.Error(err.Error())
		xutils.Exit(1)
	}
	dbc.dsn = dsn

	driver, err := dbc.adapter.GetServerDriver()
	if err != nil {
		logger.Error(err.Error())
		xutils.Exit(1)
	}
	dbc.driver = driver
}

// Connect definition
func (dbc *DBClient) Connect() {
	var err error
	dbc.DB, err = gorm.Open(dbc.driver, dbc.gormConfig)
	if err != nil {
		fmt.Println(err)
		panic("failed to establish database connection")
	}
}

// Ping definition
func (dbc *DBClient) Ping() error {
	var err error
	var d *sql.DB
	d, err = dbc.DB.DB()
	if err != nil {
		return err
		// fmt.Println(err)
		// panic("gorm error")
	}
	err = d.Ping()
	if err != nil {
		return err
		// fmt.Println(err)
		// panic("failed to ping database")
	}
	return nil
}

// CreateDatabase definition
func (dbc *DBClient) CreateDatabase() error {
	sql, err := dbc.adapter.GetDbCreateStatement()
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	err = dbc.DB.Exec(sql + dbc.config.Name + ";").Error
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	return nil
}

// DropDatabase definition
func (dbc *DBClient) DropDatabase() error {
	sql, err := dbc.adapter.GetDbDropStatement()
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}
	err = dbc.DB.Exec(sql + dbc.config.Name + ";").Error
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}
	return nil
}

// type IBaseModel interface {
// 	BeforeCreate(tx *gorm.DB) error
// }

// func (dbc *DBClient) embedsBaseModel(v interface{}) bool {
// 	rt := reflect.TypeOf(v)
// 	if rt.Kind() != reflect.Struct {
// 		return false
// 	}

// 	base := reflect.TypeOf(dbc.baseModel)
// 	for i := 0; i < rt.NumField(); i++ {
// 		if sf := rt.Field(i); sf.Type == base && sf.Anonymous {
// 			return true
// 		}
// 	}
// 	return false
// }
