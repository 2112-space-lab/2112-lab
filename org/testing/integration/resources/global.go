package resources

import (
	"context"
	"log/slog"
	"path/filepath"
	"sync"

	"github.com/org/2112-space-lab/org/testing/pkg/fx"
	testappdb "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-app-db"
	testserviceconfig "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-config"
	testservicecontainer "github.com/org/2112-space-lab/org/testing/pkg/testing/resources/test-service-container"
	models_common "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
	models_cont "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
	xtestdb "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db"
	models_db "github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-db/models"
)

var resourceManager *GlobalResourceManager
var resourceManagerInitLock sync.Mutex

type GlobalResourceManager struct {
	AppFolderRootPath       string
	RootDatabaseManager     *xtestdb.DatabaseManager
	AppDatabaseManager      *testappdb.AppDatabaseManager
	ServiceContainerManager *testservicecontainer.ServiceContainerManager
}

type CommonScenarioState interface {
	GetScenarioInfo() models_common.ScenarioInfo
	GetScenarioFolder() string
	GetLogger() *slog.Logger
}

func InitGlobalResourceManager(
	ctx context.Context,
	runLogger *slog.Logger,
	configDB models_db.DatabaseConnectionInfo,
	appRootPath string,
	appMigrationsFolder models_db.DatabaseMigrationPath,
	appServiceDockerImage models_cont.DockerContainerImage,
	propagatorServiceDockerImage models_cont.DockerContainerImage,
	dockerNetwork models_cont.NetworkName,
	testRunningEnv models_common.TestRunningEnv,
) (func(context.Context) error, error) {
	resourceManagerInitLock.Lock()
	defer resourceManagerInitLock.Unlock()
	if resourceManager != nil {
		panic("attempted to init already initialized GlobalResourceManager")
	}
	testserviceconfig.SetAppPaths(
		appRootPath,
		filepath.Join(appRootPath, "assets"),
		string(appMigrationsFolder),
	)
	rootDbm := xtestdb.GetOrInitDefaultDatabaseManager(ctx, runLogger, configDB)
	appDbm := testappdb.GetOrInitDefaultDatabaseManager(rootDbm, appMigrationsFolder, appRootPath)
	resourceManager = &GlobalResourceManager{
		AppFolderRootPath:   appRootPath,
		RootDatabaseManager: rootDbm,
		AppDatabaseManager:  appDbm,
		ServiceContainerManager: testservicecontainer.NewServiceContainerManager(
			models_cont.DockerContainerImage(appServiceDockerImage),
			models_cont.DockerContainerImage(propagatorServiceDockerImage),
			configDB,
			dockerNetwork,
			testRunningEnv,
		),
	}
	errRoles := resourceManager.RootDatabaseManager.CreateRolesDB(ctx, runLogger, "service_app")
	teardown := func(ctx context.Context) error {
		resourceManager.RootDatabaseManager.CloseManagementDatabaseConnectionPool(ctx, runLogger)
		return nil
	}
	return teardown, fx.FlattenErrorsIfAny(errRoles)
}

func GetGlobalResourceManager() *GlobalResourceManager {
	resourceManagerInitLock.Lock()
	defer resourceManagerInitLock.Unlock()
	if resourceManager == nil {
		panic("GlobalResourceManager expected to be initialized but got nil")
	}
	return resourceManager
}
