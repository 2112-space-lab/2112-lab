package xtestdb

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

// CreateRolesDB creates roles on DB if they do not exist yet
func (m *DatabaseManager) CreateRolesDB(ctx context.Context, logger *slog.Logger, rolesDB ...string) error {
	conn, err := m.managementDatabasePool.Acquire(ctx)
	if err != nil {
		logger.Error("unable to connect to db", slog.Any("error", err))
		return fmt.Errorf("unable to connect to db [%w]", err)
	}
	for _, role := range rolesDB {
		if _, err := conn.Exec(ctx, "CREATE ROLE "+role+";"); err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				logger.Error("unable to create role",
					slog.Any("error", err),
					slog.String("role", role),
				)
				return fmt.Errorf("unable to create role [%s] - [%w]", role, err)
			}
		}
	}
	return nil
}
