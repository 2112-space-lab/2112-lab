package service

import "github.com/org/2112-space-lab/org/app-service/internal/config/constants"

var client *ServiceClient

func init() {
	client = &ServiceClient{
		name: constants.FEATURE_SERVICE,
	}
}

// GetClient getters
func GetClient() *ServiceClient {
	return client
}
