// Package database provides support for creating connections to database servers.
// Package contains helper methods useful when dealing with DB structs.
// Supported database server: MySQL, PostgreSQL and MS SQL Server
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	mssqldb "github.com/denisenkom/go-mssqldb"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// SSLMode defines the SSL mode
type SSLMode string

// SSLConfig defines the SSL parameters to be use to connect to the server
type SSLConfig struct {
	// SSLMode specifies the SSL mode the client should use
	SSLMode SSLMode
	// SSLCertPath specifies the file name of the client SSL certificate; In postgre replaces the default ~/.postgresql/postgresql.crt
	SSLCertPath string
	// SSLKeyPath pecifies the location for the secret key used for the client certificate. In postgre replaces the default ~/.postgresql/postgresql.key
	SSLKeyPath string
	// SSLRootCert specifies the name of a file containing SSL certificate authority (CA) certificate(s). If the file exists,
	// the server's certificate will be verified to be signed by one of these authorities. In postgre the default is ~/.postgresql/root.crt.
	SSLRootCert string
}

const (
	// PostgreSSLModeRequire only try an SSL connection. If a root CA file is present, verify the certificate in the same way as if verify-ca was specified
	PostgreSSLModeRequire SSLMode = "require"
	// PostgreSSLModeDisable disables the SSL key validation.
	PostgreSSLModeDisable SSLMode = "disable"
	// PostgreSSLModeAllow first try a non-SSL connection; if that fails, try an SSL connection
	PostgreSSLModeAllow SSLMode = "allow"
	// PostgreSSLModeVerifyCA only try an SSL connection, and verify that the server certificate is issued by a trusted certificate authority (CA)
	PostgreSSLModeVerifyCA SSLMode = "verify-ca"
	// PostgreSSLModeVerifyFull only try an SSL connection, verify that the server certificate is issued by a trusted CA and that the server host name matches that in the certificate
	PostgreSSLModeVerifyFull SSLMode = "verify-full"
	defaultPostgreSSLMode            = PostgreSSLModeDisable
)

// DBConfig is configuration for database connection
type DBConfig struct {
	User         string
	Pass         string
	Host         string
	Port         int
	Database     string
	Timeout      time.Duration
	MaxConn      int
	MaxIdleConns int
	// SSLConfig defines the SSL connection settings to be used to connect to the server
	SSLConfig SSLConfig
	// SQLMode defines the SQL mode used for connection to the server
	SQLMode string
}

// DBConnFunc is any function that takes standard DBConfig and creates specific connection string from it
type DBConnFunc func(DBConfig) string

// MySQLConnDefault is used to construct correct mysql connection string to database
func MySQLConnDefault(c DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?sql_mode='%v'&collation=utf8mb4_unicode_ci&multiStatements=true&parseTime=true&timeout=%s",
		c.User, c.Pass, c.Host, c.Port, c.Database, c.SQLMode, c.Timeout)
}

// PostgresConnDefault is used to construct correct mysql connection string to database
func PostgresConnDefault(c DBConfig) string {
	sslMode := c.SSLConfig.SSLMode
	// avoid breaking existing applications by setting a default value
	if sslMode == "" {
		sslMode = defaultPostgreSSLMode
	}
	// `sslcert`, `sslkey` and `sslrootcert` will be ignored in case no SSL connection is configured
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
		c.Host, c.Port, c.User, c.Pass, c.Database, sslMode, c.SSLConfig.SSLCertPath, c.SSLConfig.SSLKeyPath, c.SSLConfig.SSLRootCert)
}

// SQLServerConnDefault is used to construct sqlserver connection string to database
func SQLServerConnDefault(c DBConfig) string {
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		c.Host, c.User, c.Pass, c.Port, c.Database)
}

// IsMySQLErrorCode checks if the error matches a specific MySQL error code
func IsMySQLErrorCode(err error, num uint16) bool {

	if sqlErr, ok := err.(*mysql.MySQLError); ok {
		return sqlErr.Number == num
	}

	return false
}

// IsPostgresErrorCode checks if the error matches a specific Postgres error code
func IsPostgresErrorCode(err error, code pq.ErrorCode) bool {

	if sqlErr, ok := err.(*pq.Error); ok {
		return sqlErr.Code == code
	}

	return false
}

// IsSQLServerErrorCode checks if the error matches a specific Postgres error code
func IsSQLServerErrorCode(err error, code int32) bool {

	if sqlErr, ok := err.(*mssqldb.Error); ok {
		return sqlErr.Number == code
	}

	return false
}

// CreateDBConnection creates database connection with given parameters and tries to ping target database.
// This is a blocking version of method, it will try to backoff till db connection is established
func CreateDBConnection(logger *logrus.Entry, driverType string, connFunc DBConnFunc, config DBConfig) *sql.DB {
	b := createBackoff()

	var db *sql.DB

	err := backoff.RetryNotify(func() error {
		var e error
		db, e = createDBConnection(driverType, connFunc, config)
		return e
	}, b, func(e error, d time.Duration) {
		nextInterval := math.Round(float64(d / time.Second))
		logger.WithFields(logrus.Fields{"status": "failed", "error": e}).Warnf("unable to establish database connection, retrying after %.1f seconds", nextInterval)
	})

	if err != nil {
		logger.WithFields(logrus.Fields{"status": "failed", "error": err}).Error("unable to create backoff for establishing database connection")
		return nil
	}

	// if connection is established, return true
	logger.WithFields(logrus.Fields{"status": "success"}).Infof("connected to database %v@%v:%v", config.Database, config.Host, config.Port)
	return db
}

// CreateDBConnectionOrDie creates database connection with given parameters and tries to ping target database.
// This is a simple version that will fail if connection cannot be established
func CreateDBConnectionOrDie(logger *logrus.Entry, driverType string, connFunc DBConnFunc, config DBConfig) *sql.DB {
	db, err := createDBConnection(driverType, connFunc, config)
	if err != nil {
		logger.WithFields(logrus.Fields{"status": "failed", "error": err}).Error("unable to establish database connection")
		return nil
	}

	// if connection is established, return true
	logger.WithFields(logrus.Fields{"status": "success"}).Infof("connected to database %v@%v:%v", config.Database, config.Host, config.Port)
	return db
}

// CloseDatabase helper function that can be used as a part of defer in main
func CloseDatabase(logger *logrus.Entry, db *sql.DB) {
	if db == nil {
		logger.WithFields(logrus.Fields{"status": "failed"}).Warn("database already closed!")
		return
	}
	err := db.Close()
	if err != nil {
		logger.WithFields(logrus.Fields{"status": "failed", "error": err}).Warn("unable to close database cleanly")
	}
}

// createDBConnection create database connection with given parameters and tries to ping target database.
func createDBConnection(driverType string, connFunc DBConnFunc, config DBConfig) (*sql.DB, error) {
	db, err := sql.Open(driverType, connFunc(config))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxConn)
	db.SetMaxIdleConns(config.MaxIdleConns)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createBackoff() *backoff.ExponentialBackOff {

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 0
	b.InitialInterval = 1 * time.Second
	b.MaxInterval = 30 * time.Second
	b.RandomizationFactor = .25

	b.Reset()
	return b
}

// ResourceNameToID extracts ID from given string name at given index. Index starts at 0. Number of parts should be at
// least 2. If given empty string, 0 is returned.
// E.g. for resourceName remotes/1 run ResourceNameToID("remotes/123",1,2) to return id 123
// E.g. for resourceName remotes/1/antennas/2 run ResourceNameToID("remotes/123/antennas/234",3,4) to return id 234
func ResourceNameToID(name string, targetFieldIndex uint, maxNumParts int) (int64, error) {
	if name == "" {
		return 0, nil
	}
	parts := strings.Split(name, "/")
	if len(parts) != maxNumParts {
		return -1, fmt.Errorf("expected number of fields is different from actual number (%v)", len(parts))
	}
	if int(targetFieldIndex) >= len(parts) {
		return -1, errors.New("target field index is bigger than number of fields")
	}
	id, err := strconv.ParseInt(parts[targetFieldIndex], 10, 64)
	if err != nil {
		return -1, errors.New("error converting ID to int64")
	}
	return id, nil
}
