package xtestdb

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"

	// import postgres driver
	_ "github.com/lib/pq"
	// used in testing for adding migrate file
	_ "github.com/golang-migrate/migrate/source/file"
)

const dbNameMaxLength = 63

var defaultGlobalDatabaseManager *DatabaseManager
var defaultGlobalDatabaseManagerLock sync.Mutex

type DatabaseManager struct {
	lock                   sync.RWMutex
	managementConnInfo     fx.Option[models.DatabaseConnectionInfo]
	managementDatabasePool *pgxpool.Pool
}

func GetOrInitDefaultDatabaseManager(ctx context.Context, logger *slog.Logger, connInfo models.DatabaseConnectionInfo) *DatabaseManager {
	defaultGlobalDatabaseManagerLock.Lock()
	defer defaultGlobalDatabaseManagerLock.Unlock()
	if defaultGlobalDatabaseManager == nil {
		dbm, err := NewDatabaseManager(ctx, logger, connInfo)
		if err != nil {
			logger.Error("failed to GetOrInitDefaultDatabaseManager",
				slog.Any("error", err),
				slog.Any("connInfo", connInfo),
			)
			panic(err)
		}
		defaultGlobalDatabaseManager = dbm
	}
	return defaultGlobalDatabaseManager
}

func NewDatabaseManager(ctx context.Context, logger *slog.Logger, connInfo models.DatabaseConnectionInfo) (*DatabaseManager, error) {
	m := &DatabaseManager{}
	err := m.connectToManagementDatabase(ctx, logger, connInfo)
	return m, err
}

func (m *DatabaseManager) connectToManagementDatabase(_ context.Context, logger *slog.Logger, connInfo models.DatabaseConnectionInfo) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	connectionString := connInfo.PreparePostgreConnectionString()
	logger = logger.With(
		slog.String("connnectionString", connectionString),
	)
	dbPool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		logger.With(
			slog.Any("error", err),
		).Error("failed to create management db connection pool")
		return fmt.Errorf("failed to create management db connection pool for [%s] [%w]", connectionString, err)
	}
	m.managementConnInfo = fx.NewValueOption(connInfo)
	m.managementDatabasePool = dbPool
	logger.Info("created management db connection pool")
	return nil
}

func (m *DatabaseManager) CloseManagementDatabaseConnectionPool(ctx context.Context, logger *slog.Logger) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.managementConnInfo = fx.NewEmptyOption[models.DatabaseConnectionInfo]()
	if m.managementDatabasePool == nil {
		return
	}
	m.managementDatabasePool.Close()
	m.managementDatabasePool = nil
	logger.Info("closed management database connection pool")
}
