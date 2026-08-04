package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		Http      http
		SeedNodes seedNodes
		App       app
	}
	app struct {
		Name string `env:"APP_NAME,required"`
	}
	http struct {
		Port string `env:"HTTP_PORT,required"`
	}
	seedNodes struct {
		Urls []string `env:"SEEDNODES_URLS,required"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
