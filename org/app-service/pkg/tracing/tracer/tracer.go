package server

import (
	"context"
	"net/http"

	"github.com/org/2112-space-lab/org/app-service/pkg/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	api_trace "go.opentelemetry.io/otel/trace"
)

// Tracer holds the exporter and sampler
type Tracer struct {
	Exporter api_trace.Tracer
	sampler  sdk_trace.Sampler
}

// NewTracer creates a new tracer given the sampler and exporter
func NewTracer(serviceName string, sampler sdk_trace.Sampler, exporter sdk_trace.SpanExporter) (*Tracer, error) {
	defaultResource := resource.Default()
	resource, err := resource.Merge(
		defaultResource,
		resource.NewWithAttributes(
			defaultResource.SchemaURL(),
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	spanProcessor := sdk_trace.NewBatchSpanProcessor(exporter)
	tp := sdk_trace.NewTracerProvider(
		sdk_trace.WithSampler(sampler),
		sdk_trace.WithBatcher(exporter),
		sdk_trace.WithSpanProcessor(spanProcessor),
		sdk_trace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)
	tc := propagation.TraceContext{}
	otel.SetTextMapPropagator(tc)

	return &Tracer{Exporter: tp.Tracer(""), sampler: sampler}, nil
}

// Sampler returns the sampler
func (t *Tracer) Sampler() sdk_trace.Sampler {
	return t.sampler
}

// NewServerInterceptorSpan creates a new server interceptor span
func (t *Tracer) NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, name, api_trace.WithSpanKind(api_trace.SpanKindServer))
}

// NewServerSpan creates a new server span
func (t *Tracer) NewServerSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, api_trace.WithSpanKind(api_trace.SpanKindServer))
}

// NewClientInterceptorSpan creates a new client interceptor span
func (t *Tracer) NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, name, api_trace.WithSpanKind(api_trace.SpanKindClient))
}

// NewClientSpan creates a new client span
func (t *Tracer) NewClientSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, api_trace.WithSpanKind(api_trace.SpanKindClient))
}

// NewSpan creates a new span
func (t *Tracer) NewSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller)
}

// NewSpanHTTP creates a new span from http request
func (t *Tracer) NewSpanHTTP(r *http.Request, caller string) (*http.Request, *tracing.Span) {
	ctx, span := t.NewSpan(r.Context(), caller)
	r = r.WithContext(ctx)
	return r, span
}

func (t *Tracer) newSpan(ctx context.Context, caller string, options ...api_trace.SpanStartOption) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, caller, options...)
}

func (t *Tracer) newSpanFromName(ctx context.Context, name string, options ...api_trace.SpanStartOption) (context.Context, *tracing.Span) {
	ctx, span := t.Exporter.Start(ctx, name, options...)
	return ctx, tracing.SpanWithEndOpts(span)
}
