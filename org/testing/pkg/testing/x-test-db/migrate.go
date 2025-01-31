package xtestdb

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

func PrepareMigrator(
	ctx context.Context,
	logger *slog.Logger,
	conn *DatabaseConnection,
	migrateFolder models.DatabaseMigrationPath,
) (*migrate.Migrate, error) {
	var err error

	migrationPath := fmt.Sprintf("file://%s", migrateFolder)
	if strings.Contains(migrationPath, "\\") {
		migrationPath = strings.ReplaceAll(migrationPath, "\\", "/")
	}
	sqlDB := stdlib.OpenDBFromPool(conn.ConnPool)
	if sqlDB == nil {
		return nil, fmt.Errorf("failed to open std connection from database pool")
	}
	defer func() {
		errC := sqlDB.Close()
		err = fx.FlattenErrorsIfAny(err, errC)
	}()
	if err := sqlDB.Ping(); err != nil {
		logger.Error("unable to ping DB",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("unable to ping DB [%w]", err)
	}
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		logger.Error("failed to get DB driver",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to get DB driver [%w]", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("unable to open migration: %w", err)
	}
	return migrator, nil
}
