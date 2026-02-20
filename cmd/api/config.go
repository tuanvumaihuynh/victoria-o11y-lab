package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"

	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/http"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/log"
	"github.com/tuanvumaihuynh/victoria-o11y-lab/internal/postgres"

	_ "embed"
)

//go:embed config.yml
var defaultConfigBytes []byte

type Config struct {
	Log      log.Config      `yaml:"log"`
	Postgres postgres.Config `yaml:"postgres"`
	HTTP     http.Config     `yaml:"http"`
}

func (c *Config) Validate() error {
	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("log: %w", err)
	}

	if err := c.Postgres.Validate(); err != nil {
		return fmt.Errorf("postgres: %w", err)
	}

	if err := c.HTTP.Validate(); err != nil {
		return fmt.Errorf("http: %w", err)
	}

	return nil
}

func NewConfig() (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(rawbytes.Provider(defaultConfigBytes), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("load default config: %w", err)
	}

	if err := initFlags(k); err != nil {
		return nil, fmt.Errorf("init flags: %w", err)
	}

	configFilePath := k.String("config")
	if configFilePath != "" {
		if err := k.Load(file.Provider(configFilePath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("load config file: %w", err)
		}
	}

	// Load environment variables and merge into the loaded config.
	// API_foo__bar -> foo.bar (double underscore becomes dot for nested config)
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: "API_",
		TransformFunc: func(k, v string) (string, any) {
			key := strings.ToLower(strings.TrimPrefix(k, "API_"))
			key = strings.ReplaceAll(key, "__", ".")
			return key, v
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{
		Tag: "yaml",
	}); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func initFlags(k *koanf.Koanf) error {
	f := flag.NewFlagSet("api", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	f.String("config", "", "path to config file, if not provided, the default configuration will be used")

	if err := f.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		return fmt.Errorf("load flags: %w", err)
	}

	return nil
}
