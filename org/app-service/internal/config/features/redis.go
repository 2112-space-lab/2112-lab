package features

import "github.com/org/2112-space-lab/org/app-service/internal/config/constants"

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	CacheTTL string `mapstructure:"REDIS_CACHE_TTL"`
}

var redis = &Feature{
	Name:       constants.FEATURE_REDIS,
	Config:     &RedisConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"Host",
		"Port",
		"Password",
	},
}

func init() {
	Features.Add(redis)
}

func (r *RedisConfig) GetAddr() string {
	return r.Host + ":" + r.Port
}
