package database

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	log "github.com/org/2112-space-lab/org/testing/pkg/x-log"

	// import postgres driver
	_ "github.com/lib/pq"
	// import the postgres adapter
	"time"

	// Used for testing?
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

// PostgresConnection implements a Postgres database
type PostgresConnection struct {
	Driver *sql.DB
	DB     *goqu.Database
}

// PostgresConfig is a wrapper for configuration options used to configure a database
type PostgresConfig struct {
	Username    string
	Password    string
	Host        string
	Port        int
	DBname      string
	Timeout     time.Duration
	MaxConns    int
	MaxIdeConns int
	SSLConfig   SSLConfig
}

// NewConnection generates a new database from a configuration. This method sets up the
// connection string. This method will open a connection but will not actually attempt to
// talk to the database.
func NewConnection(cfg *PostgresConfig) (*PostgresConnection, error) {
	ctxLog := log.WithFields(log.Fields{"event": "database.New"})

	if cfg == nil {
		cfg = &PostgresConfig{}
	}
	ctxLog = ctxLog.WithFields(log.Fields{
		"username": cfg.Username,
		"host":     cfg.Host,
		"port":     cfg.Port,
		"dbname":   cfg.DBname,
		"timeout":  cfg.Timeout})

	logrusEntry, _ := log.GetLogrusEntry(ctxLog)
	// Database Connection
	driver := CreateDBConnection(logrusEntry, "postgres", PostgresConnDefault, DBConfig{
		User:         cfg.Username,
		Pass:         cfg.Password,
		Host:         cfg.Host,
		Port:         cfg.Port,
		Database:     cfg.DBname,
		Timeout:      cfg.Timeout,
		MaxConn:      cfg.MaxConns,
		MaxIdleConns: cfg.MaxIdeConns,
		SSLConfig:    cfg.SSLConfig,
	})

	return &PostgresConnection{
		Driver: driver,
		DB:     goqu.New("postgres", driver),
	}, nil
}

// Close will close the database connection.
func (d *PostgresConnection) Close() {
	ctxLog := log.WithFields(log.Fields{"event": "DB.close"})

	logrusEntry, _ := log.GetLogrusEntry(ctxLog)
	CloseDatabase(logrusEntry, d.Driver)
	ctxLog.WithFields(log.Fields{"status": "success"}).Info("closing database connection")
}
