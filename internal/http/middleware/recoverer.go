package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/apperr"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http/dto"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible.
//
// Recoverer prints a stack trace of the last function call.
func Recoverer(log *slog.Logger) func(http.Handler) http.Handler {
	res := &dto.ErrorResponse{
		Code:    apperr.InternalServerErr.Code(),
		Message: apperr.InternalServerErr.Msg(),
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						panic(rvr)
					}

					log.ErrorContext(r.Context(), "http request panic", slog.Any("recover", rvr),
						slog.String("stack", string(debug.Stack())))

					if r.Header.Get("Connection") != "Upgrade" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusInternalServerError)
						if err := json.NewEncoder(w).Encode(res); err != nil {
							log.ErrorContext(r.Context(), "error encoding response", slog.Any("error", err))
						}
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
