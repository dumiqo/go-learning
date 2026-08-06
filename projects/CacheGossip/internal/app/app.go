package app

import (
	"CacheGossip/config"
	"CacheGossip/internal/api"
	"CacheGossip/internal/cache"
	"CacheGossip/internal/gossip"
	"CacheGossip/internal/membership"
	"CacheGossip/pkg/logger"
	myMiddleware "CacheGossip/pkg/middleware"
	"context"
	"encoding/json"
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
type services struct {
	logger *logger.Logger
	cache  *cache.Cache
	gossip *gossip.Gossip
}

func initServices(cfg config.Config) services {
	logger, _ := logger.NewLogger(cfg.App.Name)
	cache, _ := cache.NewCache(logger)
	var nodes map[string]string
	json.Unmarshal([]byte(cfg.SeedNodes.NodesRaw), &nodes)
	membersip := membership.NewMembership(nodes)
	gossip := gossip.NewGossip(cfg.App.Name, cfg.App.Address, membersip, cache, logger)
	return services{logger, cache, gossip}
}

func initServers(src services, cfg config.Config) servers {
	return servers{initClientApi(src, cfg), initGossip(src, cfg)}
}

func initGossip(src services, cfg config.Config) *http.Server {
	client := api.NewGossipApi(cfg.App.Name, src.gossip, src.logger)
	r := chi.NewRouter()

	r.Use(myMiddleware.Logger(src.logger))
	r.Use(middleware.Recoverer)
	r.Use(myMiddleware.JSON)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", client.Health)
	r.Post("/gossip", client.Gossip)
	r.Get("/status", client.Status)

	clientServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Gossip.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		src.logger.Info("Starting gossip API server")
		if err := clientServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			src.logger.Error("Failed to start gossip server. %s", err)
		}
	}()

	return clientServer
}

func initClientApi(src services, cfg config.Config) *http.Server {
	client, _ := api.NewClientApi(cfg.App.Name, src.cache, src.logger)

	r := chi.NewRouter()

	r.Use(myMiddleware.Logger(src.logger))
	r.Use(middleware.Recoverer)
	r.Use(myMiddleware.JSON)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", client.Health)
	r.Route("/cache", func(r chi.Router) {
		r.Post("/", client.Post)
		r.Delete("/{key}", client.Delete)
		r.Get("/{key}", client.Get)
		r.Get("/status", client.Status)
	})

	clientServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		src.logger.Info("Starting client API server")
		if err := clientServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			src.logger.Error("Failed to start client server")
		}
	}()

	return clientServer
}

func Run(cfg *config.Config) {

	services := initServices(*cfg)
	servers := initServers(services, *cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go services.gossip.Start(ctx)
	go services.cache.Start(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	services.logger.Info("Shutting down...")

	cancel()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	servers.client.Shutdown(ctx)
	servers.gossip.Shutdown(ctx)

	services.logger.Info("Shutdown copleate")
}
