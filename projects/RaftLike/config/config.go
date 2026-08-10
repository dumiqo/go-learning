package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppName app
	Http    http
}

type http struct {
	Port int `env:"HTTP_PORT,required" envDefault:"8080"`
}

type app struct {
	Name string `env:"APP_NAME,required" envDefault:"App"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
