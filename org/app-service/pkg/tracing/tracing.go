package tracing

import (
	"context"
	"net/http"

	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
)

// Tracer interface wraps the different types of tracer with a common interface
type Tracer interface {
	NewSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewClientSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewServerSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *Span)
	NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *Span)
	NewSpanHTTP(r *http.Request, caller string) (*http.Request, *Span)
	Sampler() sdk_trace.Sampler
}

// T top level tracer instance
var T Tracer

// Sampler returns the otel sampler
func Sampler() sdk_trace.Sampler {
	if T == nil {
		return sdk_trace.NeverSample()
	}
	return T.Sampler()
}

// NewSpan creates a new span from a given context
func NewSpan(ctx context.Context, caller string) (context.Context, *Span) {
	if T == nil {
		return ctx, SpanWithEndOpts(nil)
	}
	return T.NewSpan(ctx, caller)
}

// NewClientSpan creates a new span from a given context with information about the caller/client.
func NewClientSpan(ctx context.Context, caller string) (context.Context, *Span) {
	if T == nil {
		return ctx, SpanWithEndOpts(nil)
	}
	return T.NewClientSpan(ctx, caller)
}

// NewServerSpan creates a new span from a given context with information about the server.
func NewServerSpan(ctx context.Context, caller string) (context.Context, *Span) {
	if T == nil {
		return ctx, SpanWithEndOpts(nil)
	}
	return T.NewServerSpan(ctx, caller)
}

// NewClientInterceptorSpan creates a new span from a given context with information about the client interceptor.
func NewClientInterceptorSpan(ctx context.Context, caller string) (context.Context, *Span) {
	if T == nil {
		return ctx, SpanWithEndOpts(nil)
	}
	return T.NewClientInterceptorSpan(ctx, caller)
}

// NewServerInterceptorSpan creates a new span from a given context with information about the server interceptor.
func NewServerInterceptorSpan(ctx context.Context, caller string) (context.Context, *Span) {
	if T == nil {
		return ctx, SpanWithEndOpts(nil)
	}
	return T.NewServerInterceptorSpan(ctx, caller)
}

// NewSpanHTTP creates a new span from a http Request
func NewSpanHTTP(r *http.Request, caller string) (*http.Request, *Span) {
	if T == nil {
		return r, SpanWithEndOpts(nil)
	}
	return T.NewSpanHTTP(r, caller)
}

// func GetTraceIDFromSpan(ctx context.Context) string {
// 	return api_trace.SpanFromContext(ctx).SpanContext().TraceID().String()
// }

// func GetSpanIDFromSpan(ctx context.Context) string {
// 	return api_trace.SpanFromContext(ctx).SpanContext().SpanID().String()
// }
