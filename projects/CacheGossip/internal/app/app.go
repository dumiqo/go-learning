package app

import (
	"CacheGossip/config"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func Run(cfg *config.Config) {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	l.WithLevel(zerolog.InfoLevel).Msg(fmt.Sprintf("Server:%s start init", cfg.App.Name))
}
