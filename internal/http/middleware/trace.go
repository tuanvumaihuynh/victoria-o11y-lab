package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http/metrics"
)

func Trace(tracer trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skipTracingPaths(r) {
				next.ServeHTTP(w, r)
				return
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			ctx := otel.GetTextMapPropagator().Extract(
				r.Context(),
				propagation.HeaderCarrier(r.Header),
			)

			ctx, span := tracer.Start(ctx, "http.server",
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPRequestMethodKey.String(r.Method),
					semconv.URLPathKey.String(r.URL.Path),
					semconv.ServerAddressKey.String(r.Host),
					semconv.UserAgentOriginalKey.String(r.UserAgent()),
				),
			)
			defer span.End()

			r = r.WithContext(ctx)
			next.ServeHTTP(ww, r)

			// Safe route extraction
			routePattern := "<unknown>"
			if rctx := chi.RouteContext(r.Context()); rctx != nil {
				if rp := rctx.RoutePattern(); rp != "" {
					routePattern = rp
				}
			}

			span.SetName(r.Method + " " + routePattern)
			span.SetAttributes(
				semconv.HTTPRouteKey.String(routePattern),
			)

			status := ww.Status()
			span.SetAttributes(
				semconv.HTTPResponseStatusCodeKey.Int(status),
			)

			if status >= 500 {
				span.SetStatus(codes.Error, http.StatusText(status))
			}
		})
	}
}

var skipPaths = map[string]struct{}{
	metrics.Path:        {},
	"/docs":             {},
	"/docs/openapi.yml": {},
}

func skipTracingPaths(r *http.Request) bool {
	_, ok := skipPaths[r.URL.Path]
	return ok
}
