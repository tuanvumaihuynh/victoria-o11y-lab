package http

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http/dto"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http/middleware"
)

type Config struct {
	Port          uint `yaml:"port"`
	SwaggerEnabled bool `yaml:"swagger_enabled"`
}

func (h *Config) Validate() error {
	if h.Port == 0 {
		return fmt.Errorf("port is required")
	}

	return nil
}

type Service struct {
	cfg    Config
	logger *slog.Logger
}

type CleanupFunc func(ctx context.Context) error

func New(cfg Config, logger *slog.Logger) *Service {
	return &Service{
		cfg:    cfg,
		logger: logger.With(slog.String("service", "http")),
	}
}

func (s *Service) Run(ctx context.Context) (CleanupFunc, error) {
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer(s.logger),
		middleware.Logger(s.logger),
		middleware.Cors(),
	)

	api := s.newHumaAPI(r)

	s.RegisterRoutes(api)

	if s.cfg.SwaggerEnabled {
		if err := s.registerDocs(r, api); err != nil {
			return nil, fmt.Errorf("register docs: %w", err)
		}
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.Port),
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 16, // 64 KB
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("http service error listening and serving: %w", err))
		}
	}()

	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	}, nil
}

func (s *Service) newHumaAPI(r *chi.Mux) huma.API {
	huma.NewError = newHumaError(s.logger)
	huma.NewErrorWithContext = newHumaErrorWithContext(s.logger)

	cfg := huma.DefaultConfig("Victoria O11y Lab API", "1.0.0")

	// Remove $schema from all schemas
	// https://github.com/danielgtaylor/huma/issues/428
	cfg.CreateHooks = nil
	cfg.OpenAPIPath = ""
	cfg.DocsPath = ""

	api := humachi.New(r, cfg)

	return api
}

func newHumaError(logger *slog.Logger) func(status int, message string, errs ...error) huma.StatusError {
	return func(status int, message string, errs ...error) huma.StatusError {
		if len(errs) == 0 {
			// Avoid logging when huma register route
			if status != 0 {
				logger.Error(
					"no error provided in huma error handler",
					slog.Int("status", status),
					slog.String("message", message),
				)
			}
			return dto.InternalServerErrResponse
		}

		err := errs[0]
		errResp := dto.NewErrorResponse(err)
		if errResp.GetStatus() >= 500 {
			logger.Error("handler error", slog.Any("error", err))
		}

		return errResp
	}
}

func newHumaErrorWithContext(logger *slog.Logger) func(hctx huma.Context, status int, message string, errs ...error) huma.StatusError {
	return func(hctx huma.Context, status int, message string, errs ...error) huma.StatusError {
		ctx := hctx.Context()

		if len(errs) == 0 {
			// Avoid logging when huma register route
			if status != 0 {
				logger.ErrorContext(
					ctx,
					"no error provided in huma error handler",
					slog.Int("status", status),
					slog.String("message", message),
				)
			}
			return dto.InternalServerErrResponse
		}

		err := errs[0]

		errResp := dto.NewErrorResponse(err)
		if errResp.GetStatus() >= 500 {
			logger.ErrorContext(ctx, "handler error", slog.Any("error", err))
		}

		return errResp
	}
}
