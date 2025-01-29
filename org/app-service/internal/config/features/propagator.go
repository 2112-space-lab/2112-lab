package features

import xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

type PropagatorConfig struct {
	BaseUrl string `mapstructure:"PROPAGATOR_URL"`
}

var propagator = &Feature{
	Name:       xconstants.FEATURE_PROPAGATOR,
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
