package features

import (
	"github.com/org/2112-space-lab/org/app-service/internal/config/constants"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	Timezone string `mapstructure:"DB_TIMEZONE"`
	Platform string `mapstructure:"DB_PLATFORM"`
	SslMode  string `mapstructure:"DB_SSL_MODE"`
}

var database = &Feature{
	Name:       constants.FEATURE_DATABASE,
	Config:     &DatabaseConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"Name",
		"Platform",
	},
}

func init() {
	Features.Add(database)
}
