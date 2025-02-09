package features

import "github.com/org/2112-space-lab/org/app-service/internal/config/constants"

type CelestrackConfig struct {
	BaseUrl string `mapstructure:"CELESTRACK_URL"`
	Satcat  string `mapstructure:"CELESTRACK_SATCAT_URL"`
}

var celestrack = &Feature{
	Name:       constants.FEATURE_CELESTRACK,
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
