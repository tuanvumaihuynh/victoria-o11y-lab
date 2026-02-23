package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/log"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/telemetry"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/pkg/cmdutil"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error running application: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := NewConfig()
	if err != nil {
		return fmt.Errorf("new config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	logger, err := log.NewLogger(cfg.Log)
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}

	cleanupTracer, err := telemetry.InitTracer(ctx, cfg.Otel)
	if err != nil {
		return fmt.Errorf("init tracer: %w", err)
	}
	defer func() {
		if err := cleanupTracer(ctx); err != nil {
			logger.ErrorContext(ctx, "error cleaning up tracer", slog.Any("error", err))
		}
	}()

	interruptChan := cmdutil.InterruptChan()

	svc := http.New(cfg.HTTP, logger)
	cleanup, err := svc.Run(ctx)
	if err != nil {
		panic(fmt.Errorf("error running http service: %w", err))
	}

	logger.InfoContext(ctx, "http service started", slog.String("addr", fmt.Sprintf(":%d", cfg.HTTP.Port)))

	<-interruptChan

	logger.InfoContext(ctx, "http service is shutting down")
	if err := cleanup(ctx); err != nil {
		logger.ErrorContext(ctx, "error cleaning up http service", slog.Any("error", err))
	}

	logger.InfoContext(ctx, "http service is stopped")

	return nil
}
