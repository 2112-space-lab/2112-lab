package testservicecontainer

import (
	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-container/models"
)

func GetPropagatorDefaultEnv() models.EnvVarKeyValueMap {
	return map[string]string{
		"FLASK_ENV":         getEnv("FLASK_ENV", "development"),
		"SECRET_KEY":        getEnv("SECRET_KEY", "default-secret-key"),
		"SCHEMA_DIRECTORY":  getEnv("SCHEMA_DIRECTORY", "/app/graphql/schemas"),
		"REDIS_HOST":        getEnv("REDIS_HOST", "test-redis"),
		"REDIS_PORT":        getEnv("REDIS_PORT", "6379"),
		"REDIS_URL":         getEnv("REDIS_URL", "redis://test-redis:6379/0"),
		"RABBITMQ_HOST":     getEnv("RABBITMQ_HOST", "test-rabbitmq"),
		"RABBITMQ_PORT":     getEnv("RABBITMQ_PORT", "5672"),
		"RABBITMQ_USER":     getEnv("RABBITMQ_USER", "2112"),
		"RABBITMQ_PASSWORD": getEnv("RABBITMQ_PASSWORD", "2112"),
		"SERVER_NAME":       getEnv("SERVER_NAME", "localhost:5000"),
	}
}
