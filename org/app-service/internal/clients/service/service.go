package service

import (
	"reflect"

	"github.com/org/2112-space-lab/org/app-service/internal/config/features"
)

// ServiceClient definition
type ServiceClient struct {
	name   string
	config features.ServiceConfig
}

// Name getters
func (c *ServiceClient) Name() string {
	return c.name
}

// Configure sets configuration
func (c *ServiceClient) Configure(v any) {
	c.config = v.(reflect.Value).Interface().(features.ServiceConfig)
}

// GetConfig getters
func (c *ServiceClient) GetConfig() features.ServiceConfig {
	return c.config
}
