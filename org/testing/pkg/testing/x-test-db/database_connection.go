package xtestdb

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

type DatabaseConnection struct {
	Info     models.DatabaseConnectionInfo
	ConnPool *pgxpool.Pool
}

func NewDatabaseConnection(
	info models.DatabaseConnectionInfo,
	connPool *pgxpool.Pool,
) DatabaseConnection {
	return DatabaseConnection{
		Info:     info,
		ConnPool: connPool,
	}
}

func (c *DatabaseConnection) Close() {
	c.ConnPool.Close()
}

func (c *DatabaseConnection) ExecuteSqlScript(ctx context.Context, logger *slog.Logger, sqlFilePath string, placeholders map[string]string) error {
	fileContent, err := getSQLFileFromTemplate(sqlFilePath, placeholders)
	if err != nil {
		return err
	}
	sqlScript := string(fileContent)

	conn, err := c.ConnPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf(
			"cannot seed database [%s] - faile to acquire conn from pool [%w]",
			c.Info.DatabaseName,
			err,
		)
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, sqlScript)
	if err != nil {
		return fmt.Errorf(
			"error while applying seed [%s] on database [%s] - [%w]",
			sqlFilePath,
			c.Info.DatabaseName,
			err,
		)
	}
	return nil
}

func (c *DatabaseConnection) MigrateUp(ctx context.Context, logger *slog.Logger, migrationPath models.DatabaseMigrationPath) error {
	logger = logger.With(
		slog.String("dbName", string(c.Info.DatabaseName)),
		slog.String("migrationsPath", string(migrationPath)),
	)
	migrator, err := PrepareMigrator(ctx, logger, c, migrationPath)
	if err != nil {
		logger.Error("cannot migrate database - failed to prepare migrator",
			slog.Any("error", err),
		)
		return fmt.Errorf("cannot migrate database - failed to prepare migrator [%w]", err)
	}
	defer func() {
		sErr, dbErr := migrator.Close()
		err = fx.FlattenErrorsIfAny(err, sErr, dbErr)
	}()
	err = migrator.Up()
	if err != nil {
		logger.Error("cannot migrate database - failed to run migration up",
			slog.Any("error", err),
		)
		return fmt.Errorf("cannot migrate database - failed to run migration up [%w]", err)
	}
	// Create copy DB as template
	return nil
}
