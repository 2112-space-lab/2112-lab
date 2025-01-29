package tracing

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Span struct holds a telemetry span and handlers which are applied when the span ends
type Span struct {
	span trace.Span
	opts []trace.SpanEndOption
}

// SpanWithEndOpts creates a Span with optional custom end handlers
func SpanWithEndOpts(span trace.Span) *Span {
	return &Span{span: span, opts: []trace.SpanEndOption{}}
}

// End ends a span if not nil
func (s *Span) End() {
	if s.span == nil {
		return
	}

	s.span.End(s.opts...)
}

// EndWithError ends the span and adds the error to the span attributes
func (s *Span) EndWithError(err error) {
	s.SetStatusByError(err)
	s.End()
}

// SetStatusByError records the error and error code into the attribute of the span
func (s *Span) SetStatusByError(err error) {
	if s.span == nil {
		return
	}
	if err != nil {
		s.span.RecordError(err)
		s.span.SetAttributes(
			attribute.KeyValue{},
		)
	}
}

// SetAttributes sets a list of key/value attributes on the span
func (s *Span) SetAttributes(kv ...attribute.KeyValue) {
	if s.span == nil {
		return
	}
	s.span.SetAttributes(kv...)
}

// GetTraceID returns the TraceID associated with a span, if found.
func (s *Span) GetTraceID() (string, bool) {
	if s.span == nil {
		return "", false
	}
	if !s.span.SpanContext().HasTraceID() {
		return "", false
	}
	return s.span.SpanContext().TraceID().String(), true
}

// GetSpanID returns the SpanID associated with a span, if found.
func (s *Span) GetSpanID() (string, bool) {
	if s.span == nil {
		return "", false
	}
	if !s.span.SpanContext().HasSpanID() {
		return "", false
	}
	return s.span.SpanContext().SpanID().String(), true
}
