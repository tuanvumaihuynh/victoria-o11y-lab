package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http/metrics"
)

func Metrics(m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == metrics.Path {
				next.ServeHTTP(w, r)
				return
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			m.InflightRequests.Inc()
			defer m.InflightRequests.Dec()

			next.ServeHTTP(ww, r)

			duration := time.Since(t1).Seconds()
			labels := []string{r.Method, r.URL.Path}

			m.RequestsTotal.WithLabelValues(labels...).Inc()
			m.RequestDuration.WithLabelValues(labels...).Observe(duration)
		})
	}
}
