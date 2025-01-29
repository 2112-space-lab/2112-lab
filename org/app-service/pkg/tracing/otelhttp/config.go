package otelhttp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Elbujito/2112/src/app-service/pkg/tracing"
	server "github.com/Elbujito/2112/src/app-service/pkg/tracing/tracer"

	otlphttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
)

// Config holds config flags used to create a tracer
type Config struct {
	ServiceName string
	Fraction    float64
	Endpoint    string
}

// NewTracerFromConfig creates a new tracer with exports via http and sets it as the tracer for the application
func NewTracerFromConfig(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.Endpoint, _ = rawConfig["endpoint"].(string)
	c.Fraction, err = FractionFromConfig(rawConfig["fraction"])
	if err != nil {
		return err
	}
	c.ServiceName, _ = rawConfig["service"].(string)
	return c.NewTracer(c.ServiceName)
}

// FractionFromConfig parses the sampling fraction rate from the config provided
func FractionFromConfig(i interface{}) (float64, error) {
	if i == nil {
		return 0, nil
	}
	switch fraction := i.(type) {
	case float64:
		return fraction, nil
	case int:
		return float64(fraction), nil
	case string:
		f, err := strconv.ParseFloat(fraction, 64)
		if err != nil {
			return 0, fmt.Errorf("could not parse sample fraction: [%w]", err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("could not parse sample fraction: [unknown format]")
	}
}

// NewTracer creates a new tracer with sampler and http exporter
func (c *Config) NewTracer(serviceName string) error {
	sampler := sdk_trace.ParentBased(sdk_trace.TraceIDRatioBased(c.Fraction))
	exporter, err := otlphttp.New(context.Background(),
		otlphttp.WithEndpoint(c.Endpoint),
		otlphttp.WithCompression(otlphttp.GzipCompression),
	)
	if err != nil {
		return err
	}

	tracing.T, err = server.NewTracer(serviceName, sampler, exporter)
	return err
}
