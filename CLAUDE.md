# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go HTTP service demonstrating the VictoriaMetrics observability stack (metrics, traces, logs). The app exports OpenTelemetry traces via gRPC to Vector, which forwards to VictoriaTraces. Docker container logs go to VictoriaLogs, and host metrics go to VictoriaMetrics. Grafana serves as the visualization layer.

## Commands

### Run locally (requires PostgreSQL and `.env` file — copy from `.env.example`)
```bash
make run-api          # Start the API server
make run-migrate      # Run database migrations
```

### Docker Compose
```bash
make dc-up            # App + PostgreSQL
make dc-o11y-up       # Observability stack (VictoriaMetrics, VictoriaLogs, VictoriaTraces, Grafana, Vector)
make dc-all-up        # Everything
```

### Testing & Linting
```bash
make test             # go test -v --failfast ./...
make test-cov         # Generate coverage report → bin/coverage.html
make lint             # golangci-lint v2 with .golangci.yml config
```

### Database Migrations (goose)
```bash
make migrate-up
make migrate-down
make migrate-status
make migrate-create name=<migration_name>   # Creates SQL migration in internal/postgres/migration/
```

## Architecture

### Two binaries in `cmd/`
- **`cmd/api`** — HTTP API server (chi router + huma framework for OpenAPI). Initializes logger, OTel tracer, and HTTP service.
- **`cmd/migrate`** — Runs goose migrations against PostgreSQL, then exits.

### Configuration system
Each binary embeds a `config.yml` as defaults, then merges: flags → config file → environment variables (via koanf). Env var format uses prefix + double-underscore nesting: `API_POSTGRES__HOST=localhost` maps to `postgres.host`. The migrate binary uses `MIGRATE_` prefix.

### Key internal packages
- **`internal/http`** — HTTP service built on chi/v5 + huma/v2. Routes registered via `RegisterRoutes()` using a generic `registerHandler` helper. Middleware chain: recoverer → correlation ID → trace → metrics → logger → CORS.
- **`internal/telemetry`** — OTel tracer initialization. Exports traces via gRPC to a configurable collector. No-ops gracefully when `CollectorURL` is empty.
- **`internal/postgres`** — pgxpool setup with OTel tracing (otelpgx), goose migrations via embedded SQL files.
- **`internal/log`** — slog-based structured logging with tint for colored output.
- **`internal/apperr`** — Predefined application errors using zerror.

### Reusable packages in `pkg/`
- **`pkg/zerror`** — Domain error type (`ZError`) with status codes, error codes, and messages. Use `WithParent`/`WithMsg` to wrap or customize predefined errors.
- **`pkg/correlationid`** — Correlation ID context helpers.
- **`pkg/cmdutil`** — OS signal handling.

## Coding Conventions

- **Structured logging**: Uses `slog` with the `sloglint` linter enforced — use `slog.Attr` key-value pairs (not mixed args), static lowercased messages, `snake_case` keys, args on separate lines.
- **Import ordering** (enforced by goimports): stdlib → external → `github.com/tuanvumaihuynh/victoria-o11y-lab/...`
- **Error handling**: Domain errors use `pkg/zerror.ZError`. HTTP layer maps `ZError` status to HTTP status via huma's error model.
- **API routes**: All under `/api/v1`. Use `registerHandler` generic helper + separate `*Docs()` functions for OpenAPI metadata.
