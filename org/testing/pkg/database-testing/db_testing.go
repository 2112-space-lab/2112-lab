package databasetesting

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/namsral/flag"

	"github.com/doug-martin/goqu/v9"

	"github.com/golang-migrate/migrate"
	"github.com/google/uuid"

	"fmt"
	"time"

	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/lib/pq"
	"github.com/org/2112-space-lab/org/testing/pkg/database"

	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	// DuplicatedKey postgres duplicated key code
	DuplicatedKey = "23505"
	// PostgresDatabaseNameMaxLength is the maximum length possible for a Postgres DB
	PostgresDatabaseNameMaxLength = 63
)

var testCfg struct {
	migrationSubFolder string
	dbUser             string
	dbPass             string
	dbHost             string
	dbPort             int
	dbName             string
}

// InitTestDatabaseConfig prepare flags for DB unit testing
func InitTestDatabaseConfig() {
	// rand.Seed(time.Now().UTC().UnixNano()) //! deprecated and apparently not needed
	flag.StringVar(&testCfg.migrationSubFolder, "test-migration-subfolder", "", "folder containing migrations for test DB")
	flag.StringVar(&testCfg.dbUser, "test-db-user", "2112_app", "database user for testing")
	flag.StringVar(&testCfg.dbPass, "test-db-pass", "2112_app", "database password for testing")
	flag.StringVar(&testCfg.dbHost, "test-db-host", "localhost", "database IP address for testing")
	flag.IntVar(&testCfg.dbPort, "test-db-port", 5432, "database port for testing")
	flag.StringVar(&testCfg.dbName, "test-db-name", "postgres", "database name for testing")
	flag.Parse()
}

// CreateRolesDB creates roles on DB if they do not exist yet
func CreateRolesDB(connFunc func(database.DBConfig) string, rolesDB []string) {
	// create connection string
	conn0 := connFunc(database.DBConfig{
		User: testCfg.dbUser, Pass: testCfg.dbPass, Host: testCfg.dbHost, Port: testCfg.dbPort, Database: testCfg.dbName, Timeout: 90 * time.Second},
	)

	// create the database to run tests against
	db0, err := sql.Open("postgres", conn0)
	if err != nil {
		log.Fatalf("unable to connect to db %q: %s", conn0, err)
	}
	for _, role := range rolesDB {
		if _, err := db0.Exec("CREATE ROLE " + role + ";"); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				log.Fatal("unable to create role:", err)
			}
		}
	}
}

// NewTestDBHandler inits a database connection and a handler for testing purpose
func NewTestDBHandler(t *testing.T, connFunc func(database.DBConfig) string, migrationSubFolder string) (db *database.PostgresConnection, dbh *database.HandlerDB, drop func()) {
	db, drop = NewTestDB(t, connFunc, migrationSubFolder)
	dbh = database.NewHandlerDB(db.DB)
	return db, dbh, drop
}

// NewTestDB is used to create tmp database on Postgres server to be used in unit tests
func NewTestDB(t *testing.T, connFunc func(database.DBConfig) string, subfolder string) (*database.PostgresConnection, func()) {
	// displayName the database something unique
	unixTs := fmt.Sprintf("%d", time.Now().Unix())
	randStr := strings.Replace(uuid.New().String(), "-", "_", -1)[:6] // randomness to prevent collisions
	testName := strings.ToLower(t.Name())
	if len(testName) > PostgresDatabaseNameMaxLength-len(unixTs+randStr) {
		testName = testName[:PostgresDatabaseNameMaxLength-len(unixTs+randStr)]
	}
	name := fmt.Sprintf("%s_%v_%s",
		testName,
		unixTs, // time (for lingering tables)
		randStr,
	)

	// subtests name is not valid name for database (e.g. TestImportCallRates/valid_file). "/" should be removed.
	name = strings.Replace(name, "/", "", -1)
	// trim the name to fit max 64 characters
	if len(name) > 64 {
		name = name[:64]
	}
	// create connection string
	conn0 := connFunc(database.DBConfig{
		User: testCfg.dbUser, Pass: testCfg.dbPass, Host: testCfg.dbHost, Port: testCfg.dbPort, Database: testCfg.dbName, Timeout: 90 * time.Second},
	)

	// test DB connection
	conn := connFunc(database.DBConfig{
		User: testCfg.dbUser, Pass: testCfg.dbPass, Host: testCfg.dbHost, Port: testCfg.dbPort, Database: name, Timeout: 90 * time.Second},
	)

	// create the database to run tests against
	db0, err := sql.Open("postgres", conn0)
	if err != nil {
		t.Fatalf("unable to connect to db %q: %s", conn0, err)
	}
	if _, err := db0.Exec("CREATE DATABASE " + name + ";"); err != nil {
		t.Fatal("unable to create database:", err)
	}

	// define a func to cleanup the created db
	cleanup := func() {
		if _, err := db0.Exec("REVOKE CONNECT ON DATABASE " + name + " FROM public;"); err != nil {
			t.Error("unable to revoke connect database:", err)
		}
		if _, err := db0.Exec("SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname='" + name + "';"); err != nil {
			t.Error("unable to terminate all connections to database:", err)
		}
		if _, err := db0.Exec("DROP DATABASE " + name + ";"); err != nil {
			t.Error("unable to drop database:", err)
		}
		if err := db0.Close(); err != nil {
			t.Error("unable to close db0:", err)
		}
	}

	db, err := migrateDB(conn, subfolder)
	if err != nil {
		cleanup()
		t.Fatal(err)
	}
	gdb := goqu.New("postgres", db)

	// return a function to do the cleanup
	return &database.PostgresConnection{
			Driver: db,
			DB:     gdb,
		}, func() {
			if err := db.Close(); err != nil {
				panic(fmt.Errorf("unable to close db connection: %v", err))
			}
			cleanup()
		}
}

// migrateDB migrates the DB. Subfolder is used to address differents migrations
func migrateDB(conn, subfolder string) (*sql.DB, error) {
	// Create db connection for tests to run against.
	if !strings.HasPrefix(subfolder, "/") {
		subfolder = "/" + subfolder
	}
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db %q: %w", conn, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping db: %w", err)
	}

	// Run migrations against newly created database.
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed creating postgres driver: %w", err)
	}

	curDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed determining current directory: %w", err)
	}

	subStrIndex := strings.Index(curDir, "internal")
	if subStrIndex == -1 {
		return nil, errors.New("unable to determine path to migrations")
	}

	curDir = curDir[0:subStrIndex]

	migrationPath := fmt.Sprintf("file://%s/assets%s", curDir, subfolder)
	if strings.Contains(migrationPath, "\\") {
		migrationPath = strings.ReplaceAll(migrationPath, "\\", "/")
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to open migration: %w", err)
	}

	if err := m.Up(); err != nil {
		return nil, fmt.Errorf("unable to migrate database up: %w", err)
	}

	return db, nil
}
