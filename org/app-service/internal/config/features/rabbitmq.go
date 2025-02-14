package features

import "github.com/org/2112-space-lab/org/app-service/internal/config/constants"

// RabbitMQConfig defines the structure for RabbitMQ settings.
type RabbitMQConfig struct {
	Host        string `mapstructure:"RABBITMQ_HOST"`
	Port        string `mapstructure:"RABBITMQ_PORT"`
	User        string `mapstructure:"RABBITMQ_USER"`
	Password    string `mapstructure:"RABBITMQ_PASSWORD"`
	VHost       string `mapstructure:"RABBITMQ_VHOST"`
	InputQueue  string `mapstructure:"RABBITMQ_INPUT_QUEUE"`
	OutputQueue string `mapstructure:"RABBITMQ_OUTPUT_QUEUE"`
}

// GetAddr returns the full RabbitMQ connection address.
func (r *RabbitMQConfig) GetAddr() string {
	return "amqp://" + r.User + ":" + r.Password + "@" + r.Host + ":" + r.Port + "/" + r.VHost
}

// Feature definition for RabbitMQ.
var rabbitmq = &Feature{
	Name:       constants.FEATURE_RABBITMQ,
	Config:     &RabbitMQConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"Host",
		"Port",
		"User",
		"Password",
		"InputQueue",
		"OutputQueue",
	},
}

func init() {
	Features.Add(rabbitmq)
}
