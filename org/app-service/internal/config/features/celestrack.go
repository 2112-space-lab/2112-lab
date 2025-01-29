package features

import xconstants "github.com/org/2112-space-lab/org/go-utils/pkg/fx/xconstants"

type CelestrackConfig struct {
	BaseUrl string `mapstructure:"CELESTRACK_URL"`
	Satcat  string `mapstructure:"CELESTRACK_SATCAT_URL"`
}

var celestrack = &Feature{
	Name:       xconstants.FEATURE_CELESTRACK,
	Config:     &CelestrackConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"BaseUrl",
		"SatcatUrl",
	},
}

func init() {
	Features.Add(celestrack)
}
