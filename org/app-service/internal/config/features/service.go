package features

import "github.com/org/2112-space-lab/org/app-service/internal/config/constants"

type ServiceConfig struct {
	Host                   string `mapstructure:"HOST"`
	ProtectedApiPort       string `mapstructure:"PROTECTED_API_PORT"`
	PublicApiPort          string `mapstructure:"PUBLIC_API_PORT"`
	HiddenApiPort          string `mapstructure:"HIDDEN_API_PORT"`
	LogLevel               string `mapstructure:"LOG_LEVEL"`
	RequestTimeoutDuration string `mapstructure:"REQUEST_TIMEOUT_DURATION"`
	WatcherSleepInterval   string `mapstructure:"WATCHER_SLEEP_INTERVAL"`
}

var service = &Feature{
	Name:       constants.FEATURE_SERVICE,
	Config:     &ServiceConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"Host",
		"ProtectedApiPort",
		"PublicApiPort",
		"HiddenApiPort",
		"LogLevel",
		"RequestTimeoutDuration",
		"WatcherSleepInterval",
	},
}

func init() {
	Features.Add(service)
}
