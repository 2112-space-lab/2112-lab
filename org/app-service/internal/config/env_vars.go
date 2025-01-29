package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Elbujito/2112/src/app-service/internal/config/features"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"
	"github.com/org/2112-space-lab/org/go-utils/pkg/fx/xutils"

	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
	"github.com/spf13/viper"
)

type EnvVars struct {
	DisableFeatures []string                  `mapstructure:"DISABLE_FEATURES"`
	Service         features.ServiceConfig    `mapstructure:",squash"`
	Database        features.DatabaseConfig   `mapstructure:",squash"`
	Redis           features.RedisConfig      `mapstructure:",squash"`
	Celestrack      features.CelestrackConfig `mapstructure:",squash"`
	Propagator      features.PropagatorConfig `mapstructure:",squash"`
	Clerk           features.ClerkConfig      `mapstructure:",squash"`
}

func (c *EnvVars) Init() {
	logger.Info("Initializing environment ...")
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	c.registerFieldsMapstructure(reflect.ValueOf(c))
	c.setDefaults()

	viper.AutomaticEnv()

	logger.Info("Reading env vars ...")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("file '.env' not found")
		fmt.Println("attempting to read env vars ...")
	}
	if err := viper.Unmarshal(c); err != nil {
		fmt.Println("error decoding config file: ", err.Error())
		os.Exit(1)
	}
}

func (c *EnvVars) registerFieldsMapstructure(vOfConfig reflect.Value) {
	if vOfConfig.Kind() == reflect.Ptr {
		vOfConfig = vOfConfig.Elem()
	}
	for i := 0; i < vOfConfig.NumField(); i++ {
		configField := vOfConfig.Field(i)
		switch configField.Kind() {
		case reflect.Struct:
			c.registerFieldsMapstructure(configField)
		case reflect.Slice:
			for j := 0; j < configField.Len(); j++ {
				c.registerFieldsMapstructure(configField.Index(i))
			}
		case reflect.String:
			viper.BindEnv(vOfConfig.Type().Field(i).Tag.Get("mapstructure"))
		}
	}
}

func (c *EnvVars) setDefaults() {
	logger.Info("Configuring default settings ...")
	viper.SetDefault("HOST", xconstants.DEFAULT_HOST)
	viper.SetDefault("PROTECTED_API_PORT", xconstants.DEFAULT_PROTECTED_API_PORT)
	viper.SetDefault("PUBLIC_API_PORT", xconstants.DEFAULT_PUBLIC_API_PORT)
	viper.SetDefault("LOG_LEVEL", xconstants.DEFAULT_LOG_LEVEL)
	viper.SetDefault("REQUEST_TIMEOUT_DURATION", strconv.Itoa(xconstants.DEFAULT_REQUEST_TIMEOUT_DURATION))
	viper.SetDefault("WATCHER_SLEEP_INTERVAL", strconv.Itoa(xconstants.DEFAULT_WATCHER_SLEEP_INTERVAL))

	viper.SetDefault("DB_PLATFORM", xconstants.DEFAULT_DB_PLATFORM)
	viper.SetDefault("DB_NAME", xconstants.DEFAULT_SQLITE_DB_NAME)
	viper.SetDefault("DB_SSL_MODE", xconstants.DEFAULT_DB_SSL_MODE)
	viper.SetDefault("DB_TIMEZONE", xconstants.DEFAULT_DB_TIMEZONE)

	viper.SetDefault("CORS_ALLOW_ORIGINS", xconstants.DEFAULT_CORS_ALLOW_ORIGINS)

	viper.SetDefault("CELESTRACK_URL", xconstants.DEFAULT_PUBLIC_CESLESTRACK_URL)
	viper.SetDefault("PROPAGATOR_URL", xconstants.DEFAULT_PRIVATE_PROPAGATOR_URL)
	viper.SetDefault("CELESTRACK_SATCAT_URL", xconstants.DEFAULT_PUBLIC_CESLESTRACK_SATCAT_URL)
}

func (c *EnvVars) OverrideUsingFlags() {
	if HostFlag != "" {
		c.Service.Host = HostFlag
	}
	if ProtectedPortFlag != "" {
		c.Service.ProtectedApiPort = ProtectedPortFlag
	}
	if PublicPortFlag != "" {
		c.Service.PublicApiPort = PublicPortFlag
	}
}

func (c *EnvVars) OverrideLoggerUsingFlags() {
	if LogLevelFlag == "" {
		return
	}
	if !xutils.StrInArr(LogLevelFlag, xconstants.LOG_LEVELS) {
		panic("Invalid log level Valid options")
	}
	c.Service.LogLevel = LogLevelFlag
}

func (c *EnvVars) SetDevMode() {
	c.Service.Host = xconstants.DEFAULT_DEV_HOST
	c.Service.LogLevel = xconstants.DEFAULT_DEV_LOG_LEVEL
}

func (c *EnvVars) GetConfigByName(name string) (string, error) {
	v := reflect.ValueOf(*c)
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if typeOfS.Field(i).Name == name {
			return v.Field(i).Interface().(string), nil
		}
	}
	return "", fmt.Errorf("config not found: %s", name)
}

func (c *EnvVars) FeatureInDisabledList(name string) bool {
	for _, v := range c.DisableFeatures {
		if v == name {
			return true
		}
	}
	return false
}
