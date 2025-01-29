package service

import xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

var client *ServiceClient

func init() {
	client = &ServiceClient{
		name: xconstants.FEATURE_SERVICE,
	}
}

// GetClient getters
func GetClient() *ServiceClient {
	return client
}
