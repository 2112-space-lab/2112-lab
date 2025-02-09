package features

import "github.com/org/2112-space-lab/org/app-service/internal/config/constants"

type PropagatorConfig struct {
	BaseUrl string `mapstructure:"PROPAGATOR_URL"`
}

var propagator = &Feature{
	Name:       constants.FEATURE_PROPAGATOR,
	Config:     &PropagatorConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"BaseUrl",
	},
}

func init() {
	Features.Add(propagator)
}
