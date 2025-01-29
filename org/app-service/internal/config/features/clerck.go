package features

type ClerkConfig struct {
	CLERK_API_KEY string `mapstructure:"CLERK_API_KEY"`
}

var clerk = &Feature{
	Name:       "clerk",
	Config:     &ClerkConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"ClerkAPIKey",
	},
}

func init() {
	Features.Add(clerk)
}
