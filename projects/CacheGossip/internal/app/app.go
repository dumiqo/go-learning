package app

import (
	"CacheGossip/config"
	"CacheGossip/internal/api"
	"CacheGossip/internal/cache"
	"CacheGossip/pgk/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	logger, _ := logger.NewLogger(&cfg.App.Name)

	cache, _ := cache.NewCache()
	client, _ := api.NewClientApi(&cfg.App.Name, cache, logger)

	clientMux := http.NewServeMux()

	clientMux.HandleFunc("/health", client.Health)

	clientServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      clientMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("Starting client API server")
		if err := clientServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start client server")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientServer.Shutdown(ctx)

	logger.Info("Shutdown copleate")
}
