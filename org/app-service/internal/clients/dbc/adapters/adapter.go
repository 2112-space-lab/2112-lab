package adapters

import (
	"github.com/Elbujito/2112/src/app-service/internal/config/features"
	xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

	"gorm.io/gorm"
)

// Adapter constructor
var Adapters = &Adapter{
	defaultPlatform: xconstants.DEFAULT_DB_PLATFORM,
	currentPlatform: xconstants.DEFAULT_DB_PLATFORM,
	adapters:        make(map[string]IAdapter),
}

// IAdapter interface
type IAdapter interface {
	SetConfig(config features.DatabaseConfig)
	GetDriver() (gorm.Dialector, error)
	GetServerDriver() (gorm.Dialector, error)
	GetDSN() (string, error)
	GetServerDSN() (string, error)
	GetDbCreateStatement() (string, error)
	GetDbDropStatement() (string, error)
	ValidateConfig() error
}

// Adapter definition
type Adapter struct {
	IAdapter
	adapters        map[string]IAdapter
	defaultPlatform string
	currentPlatform string
	config          features.DatabaseConfig
}

// SetConfig sets configs
func (a *Adapter) SetConfig(config features.DatabaseConfig) {
	a.config = config
	a.currentPlatform = config.Platform
	for _, adapter := range a.adapters {
		adapter.SetConfig(a.config)
	}
}

// GetDriver sets configs
func (a *Adapter) GetDriver() (gorm.Dialector, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDriver()
	}
	return nil, xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

// GetServerDriver return server driver
func (a *Adapter) GetServerDriver() (gorm.Dialector, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDriver()
	}
	return nil, xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

// GetDSN returns DSN
func (a *Adapter) GetDSN() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDSN()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

// GetServerDSN returns DSN
func (a *Adapter) GetServerDSN() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetServerDSN()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

// AppendAdapter add adapter
func (a *Adapter) AppendAdapter(name string, adapter IAdapter) {
	a.adapters[name] = adapter
}

// AppendAdapter get db create statement
func (a *Adapter) GetDbCreateStatement() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDbCreateStatement()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

// AppendAdapter get db statement
func (a *Adapter) GetDbDropStatement() (string, error) {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.GetDbDropStatement()
	}
	return "", xconstants.ERROR_UNKNOWN_DB_PLATFORM
}

// ValidateConfig validates config
func (a *Adapter) ValidateConfig() error {
	if adapter, ok := a.adapters[a.currentPlatform]; ok {
		return adapter.ValidateConfig()
	}
	return xconstants.ERROR_UNKNOWN_DB_PLATFORM
}
