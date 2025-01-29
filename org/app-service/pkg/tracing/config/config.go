package config

import (
	"fmt"

	otelHttpExporter "github.com/Elbujito/2112/src/app-service/pkg/tracing/otelhttp"
	otelStdoutExporter "github.com/Elbujito/2112/src/app-service/pkg/tracing/otelstdout"
)

// Config holds config flags used to create a tracer
type Config struct {
	Type   string
	Config map[string]interface{} `mapstructure:",remain"`
}

// NewTracer creates a tracer from a given config
func (c *Config) NewTracer() error {
	t, ok := tracer[c.Type]
	if !ok {
		return fmt.Errorf("tracer type %s not supported", c.Type)
	}

	return t(c.Config)
}

var tracer = map[string]func(map[string]interface{}) error{
	"otelhttp": otelHttpExporter.NewTracerFromConfig,
	"stdout":   otelStdoutExporter.NewTracerFromConfig,
	"":         NoTracer,
}

// NoTracer return a noop tracer
func NoTracer(_ map[string]interface{}) error {
	return nil
}
