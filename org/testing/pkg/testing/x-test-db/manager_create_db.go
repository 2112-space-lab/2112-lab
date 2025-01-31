package xtestdb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (m *DatabaseManager) CreateDatabase(
	ctx context.Context,
	logger *slog.Logger,
	pDbName models.PotentialDatabaseName,
) (DatabaseConnection, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	dbName := normalizeDatabaseName(logger, pDbName)
	logger = logger.With(
		slog.String("dbName", string(dbName)),
	)

	if m.managementDatabasePool == nil {
		logger.Error("cannot create database - database server management connection pool not initialized")
		return DatabaseConnection{}, fmt.Errorf("cannot create database [%s] - database server management connection pool not initialized", dbName)
	}

	conn, err := m.managementDatabasePool.Acquire(ctx)
	if err != nil {
		logger.With(
			slog.Any("error", err),
		).Error("cannot create database - failed to acquire management connection from pool")
		return DatabaseConnection{}, fmt.Errorf("cannot create database [%s] - failed to acquire management connection from pool [%s]", dbName, err)
	}
	defer conn.Release()
	t, err := conn.Exec(ctx, "CREATE DATABASE "+string(dbName)+";")
	if err != nil {
		logger.With(
			slog.String("pgCommandTag", t.String()),
			slog.Any("error", err),
		).Error("failed to create database")
		return DatabaseConnection{}, fmt.Errorf("failed to create database [%s] - exec failed [%w]", dbName, err)
	}

	connInfo := models.NewDatabaseConnectionInfo(
		m.managementConnInfo.Value.HostName,
		m.managementConnInfo.Value.Port,
		m.managementConnInfo.Value.HostNameDocker,
		m.managementConnInfo.Value.PortDocker,
		dbName,
		m.managementConnInfo.Value.PoolMaxConn,
		m.managementConnInfo.Value.OwnerUser,
		m.managementConnInfo.Value.OwnerPassword,
		m.managementConnInfo.Value.TLSConfig,
	)
	connectionString := connInfo.PreparePostgreConnectionString()
	logger = logger.With(
		slog.String("connnectionString", connectionString),
	)
	dbPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.With(
			slog.Any("error", err),
		).Error("failed to create db connection pool")
		return DatabaseConnection{}, fmt.Errorf("failed to create db connection pool for [%s] - [%w]", connectionString, err)
	}
	logger.Info("created db connection pool")

	return NewDatabaseConnection(connInfo, dbPool), nil
}
