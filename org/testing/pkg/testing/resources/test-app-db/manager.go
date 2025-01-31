package testappdb

import (
	"sync"

	xtestdb "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db"
	models_db "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

var defaultAppDatabaseManager *AppDatabaseManager
var defaultAppDatabaseManagerLock sync.Mutex

func GetOrInitDefaultDatabaseManager(
	rootDatabaseManager *xtestdb.DatabaseManager,
	migrationsPath models_db.DatabaseMigrationPath,
	appWorkspaceRootPath string,
) *AppDatabaseManager {
	defaultAppDatabaseManagerLock.Lock()
	defer defaultAppDatabaseManagerLock.Unlock()
	if defaultAppDatabaseManager == nil {
		defaultAppDatabaseManager = NewAppDatabaseManager(rootDatabaseManager, migrationsPath, appWorkspaceRootPath)
	}
	return defaultAppDatabaseManager
}

type AppDatabaseManager struct {
	rootDatabaseManager  *xtestdb.DatabaseManager
	migrationsPath       models_db.DatabaseMigrationPath
	appWorkspaceRootPath string
}

func NewAppDatabaseManager(
	rootDatabaseManager *xtestdb.DatabaseManager,
	migrationsPath models_db.DatabaseMigrationPath,
	appWorkspaceRootPath string,
) *AppDatabaseManager {
	m := &AppDatabaseManager{
		rootDatabaseManager:  rootDatabaseManager,
		migrationsPath:       migrationsPath,
		appWorkspaceRootPath: appWorkspaceRootPath,
	}
	return m
}
