package postgres

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	//nolint:gosec
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
	SSLMode  string `yaml:"ssl_mode"`

	MaxConns        int32         `yaml:"max_conns"`
	MinConns        int32         `yaml:"min_conns"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
}

func (p *Config) Validate() error {
	if p.Host == "" {
		return fmt.Errorf("host is required")
	}
	if p.Port == 0 {
		return fmt.Errorf("port is required")
	}
	if p.User == "" {
		return fmt.Errorf("user is required")
	}
	if p.Password == "" {
		return fmt.Errorf("password is required")
	}
	if p.DB == "" {
		return fmt.Errorf("db is required")
	}

	allowedSSLModes := []string{"disable", "require", "verify-full"}
	if p.SSLMode != "" && !slices.Contains(allowedSSLModes, p.SSLMode) {
		return fmt.Errorf("ssl mode must be one of the following values: %s", strings.Join(allowedSSLModes, ", "))
	}
	if p.MaxConns <= 0 {
		return fmt.Errorf("max conns must be greater than 0")
	}
	if p.MinConns <= 0 {
		return fmt.Errorf("min conns must be greater than 0")
	}
	if p.MaxConnLifetime <= 0 {
		return fmt.Errorf("max conn lifetime must be greater than 0")
	}
	if p.MaxConnIdleTime <= 0 {
		return fmt.Errorf("max conn idle time must be greater than 0")
	}

	return nil
}

// NewPgxPool creates a new pgxpool.Pool with the given configuration.
func NewPgxPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	connStr := connectionString(cfg)
	pgConf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pgConf.MaxConns = cfg.MaxConns
	pgConf.MinConns = cfg.MinConns
	pgConf.MaxConnLifetime = cfg.MaxConnLifetime
	pgConf.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, pgConf)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// Create a context with timeout for ping
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	if err := pool.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// connectionString constructs the connection string for the database.
func connectionString(cfg Config) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Path:   cfg.DB,
	}
	q := u.Query()
	q.Add("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}
