package app

import (
	"CacheGossip/config"
	"fmt"
)

func Run(cfg *config.Config) {
	fmt.Println(cfg.App.Name)
}
