package adapters

import (
	"fmt"
	"reflect"

	"github.com/org/2112-space-lab/org/app-service/internal/config/constants"
	"github.com/org/2112-space-lab/org/app-service/internal/config/features"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	Adapters.AppendAdapter(constants.DB_PLATFORM_POSTGRES, &PostgresAdapter{
		requiredConfig: []string{"Host", "Port", "User", "Password", "Name", "SslMode", "Timezone"},
	})
}

// PostgresAdapter definition
type PostgresAdapter struct {
	IAdapter
	config         features.DatabaseConfig
	requiredConfig []string
}

func (a *PostgresAdapter) SetConfig(config features.DatabaseConfig) {
	a.config = config
}

// GetDriver get statement
func (a *PostgresAdapter) GetDriver() (gorm.Dialector, error) {
	dsn, _ := a.GetDSN()
	return postgres.Open(dsn), nil
}

// GetServerDriver get statement
func (a *PostgresAdapter) GetServerDriver() (gorm.Dialector, error) {
	dsn, _ := a.GetServerDSN()
	return postgres.Open(dsn), nil
}

// GetDSN get statement
func (a *PostgresAdapter) GetDSN() (string, error) {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		a.config.Host,
		a.config.User,
		a.config.Password,
		a.config.Name,
		a.config.Port,
		a.config.SslMode,
		a.config.Timezone), nil
}

// GetServerDSN get statement
func (a *PostgresAdapter) GetServerDSN() (string, error) {
	return fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=%s TimeZone=%s",
		a.config.Host,
		a.config.User,
		a.config.Password,
		a.config.Port,
		a.config.SslMode,
		a.config.Timezone), nil
}

// GetDbCreateStatement get statement
func (a *PostgresAdapter) GetDbCreateStatement() (string, error) {
	return "CREATE DATABASE IF NOT EXISTS ", nil
}

// GetDbDropStatement get drop statement
func (a *PostgresAdapter) GetDbDropStatement() (string, error) {
	return "DROP DATABASE IF EXISTS ", nil
}

// ValidateConfig validate config
func (a *PostgresAdapter) ValidateConfig() error {
	vOfConfig := reflect.ValueOf(a.config)
	if vOfConfig.Kind() == reflect.Ptr {
		vOfConfig = vOfConfig.Elem()
	}
	for _, requiredField := range a.requiredConfig {
		for i := 0; i < vOfConfig.NumField(); i++ {
			configField := vOfConfig.Field(i)
			if configField.Kind() == reflect.String {
				if vOfConfig.Type().Field(i).Name == requiredField &&
					vOfConfig.Field(i).Interface().(string) == "" {
					return fmt.Errorf("database adapter requirements not satisfied. missing required field: %s", vOfConfig.Type().Field(i).Name)
				}
			}
		}
	}
	return nil
}
