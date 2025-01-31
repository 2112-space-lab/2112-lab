package testserviceconfig

import "sync"

var assetsLock sync.Mutex

var appWorkspaceRootPath string
var appAssetsPath string
var appMigrationPath string

func SetAppPaths(
	appWorkspaceRoot string,
	appAssets string,
	appMigration string,
) {
	assetsLock.Lock()
	defer assetsLock.Unlock()

	appWorkspaceRootPath = appWorkspaceRoot
	appAssetsPath = appAssets
	appMigrationPath = appMigration
}

func GetAppWorkspaceRootPath() string {
	assetsLock.Lock()
	defer assetsLock.Unlock()
	return appWorkspaceRootPath
}

func GetAppAssetsPath() string {
	assetsLock.Lock()
	defer assetsLock.Unlock()
	return appAssetsPath
}

func GetAppMigrationPath() string {
	assetsLock.Lock()
	defer assetsLock.Unlock()
	return appMigrationPath
}
