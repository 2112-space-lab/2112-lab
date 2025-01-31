package xtestdb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
)

func (m *DatabaseManager) DropAllDatabases(ctx context.Context, logger *slog.Logger, databases ...DatabaseConnection) error {
	errs := []error{}
	for _, d := range databases {
		err := m.DropDatabase(ctx, logger, d)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return fx.FlattenErrorsIfAny(errs...)
}

func (m *DatabaseManager) DropDatabase(ctx context.Context, logger *slog.Logger, database DatabaseConnection) (err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	dbName := database.Info.DatabaseName

	logger = logger.With(
		slog.String("dbName", string(dbName)),
	)
	defer func() {
		if err != nil {
			logger.With(
				slog.Any("error", err),
			).Error("failed to drop database")
		} else {
			logger.Info("database dropped successfully")
		}
	}()

	if m.managementDatabasePool == nil {
		return fmt.Errorf("cannot drop database [%s] - database server management connection pool not initialized", dbName)
	}

	conn, err := m.managementDatabasePool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("cannot drop database [%s] - failed to acquire management connection from pool [%s]", dbName, err)
	}
	defer conn.Release()
	if database.ConnPool != nil {
		database.ConnPool.Close()
	}

	if _, err := conn.Exec(ctx, "REVOKE CONNECT ON DATABASE "+string(dbName)+" FROM public;"); err != nil {
		return fmt.Errorf("unable to revoke connect database: [%w]", err)
	}
	if _, err := conn.Exec(ctx,
		"SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname=$1;",
		string(dbName)); err != nil {
		return fmt.Errorf("unable to terminate all connections to database: [%w]", err)
	}
	if _, err := conn.Exec(ctx, "DROP DATABASE "+string(dbName)+";"); err != nil {
		return fmt.Errorf("unable to drop database: [%w]", err)
	}
	return nil
}
