package testservicecontainer

import (
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
)

func GetAppDefaultEnv() models.EnvVarKeyValueMap {
	return map[string]string{
		"REDIS_HOST":        getEnv("REDIS_HOST", "test-redis"),
		"REDIS_PORT":        getEnv("REDIS_PORT", "6379"),
		"REDIS_URL":         getEnv("REDIS_URL", "redis://test-redis:6379/0"),
		"RABBITMQ_HOST":     getEnv("RABBITMQ_HOST", "test-rabbitmq"),
		"RABBITMQ_PORT":     getEnv("RABBITMQ_PORT", "5672"),
		"RABBITMQ_USER":     getEnv("RABBITMQ_USER", "2112"),
		"RABBITMQ_PASSWORD": getEnv("RABBITMQ_PASSWORD", "2112"),
		"RABBITMQ_QUEUE":    getEnv("RABBITMQ_QUEUE", "satellite_positions"),
	}
}
