package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		Http      http
		Gossip    gossip
		SeedNodes seedNodes
		App       app
	}
	app struct {
		Name    string `env:"APP_NAME,required"`
		Address string `env:"APP_ADDRESS,required"`
	}
	http struct {
		Port int `env:"HTTP_PORT,required"`
	}
	gossip struct {
		Port int `env:"GOSSIP_PORT,required"`
	}
	seedNodes struct {
		NodesRaw string `env:"SEEDNODES_NODES,required"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
