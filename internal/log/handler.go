package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/pkg/correlationid"
)

var _ slog.Handler = (*traceHandler)(nil)

// traceHandler enriches logs with trace and correlation data
type traceHandler struct {
	h slog.Handler
}

func newTraceHandler(h slog.Handler) traceHandler {
	return traceHandler{h: h}
}

func (eh traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return eh.h.Enabled(ctx, level)
}

func (eh traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if correlationID, ok := correlationid.FromContext(ctx); ok {
		r.Add("correlation_id", slog.StringValue(correlationID))
	}

	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		r.Add("trace_id", slog.StringValue(span.SpanContext().TraceID().String()))
		r.Add("span_id", slog.StringValue(span.SpanContext().SpanID().String()))
	}

	return eh.h.Handle(ctx, r)
}

func (eh traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newTraceHandler(eh.h.WithAttrs(attrs))
}

func (eh traceHandler) WithGroup(name string) slog.Handler {
	return newTraceHandler(eh.h.WithGroup(name))
}
