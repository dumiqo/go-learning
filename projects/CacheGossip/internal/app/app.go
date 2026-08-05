package app

import (
	"CacheGossip/config"
	"CacheGossip/internal/api"
	"CacheGossip/internal/cache"
	"CacheGossip/pkg/logger"
	myMiddleware "CacheGossip/pkg/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type servers struct {
	client *http.Server
	gossip *http.Server
}

func Run(cfg *config.Config) {
	logger, _ := logger.NewLogger(cfg.App.Name)
	cache, _ := cache.NewCache()
	client, _ := api.NewClientApi(cfg.App.Name, cache, logger)

	r := chi.NewRouter()

	r.Use(myMiddleware.LoggerMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(myMiddleware.JSONMiddleware)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", client.Health)
	r.Route("/objects", func(r chi.Router) {
		r.Post("/", client.Post)
		r.Delete("/{key}", client.Delete)
		r.Get("/{key}", client.Get)
	})

	clientServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      r,
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
