package app

import (
	"RaftLike/config"
	"RaftLike/internal/api"
	"RaftLike/internal/server"
	"RaftLike/pkg/logger"
	myMiddleware "RaftLike/pkg/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	config *config.Config
	logger *logger.Logger
	server *server.Server
}

func NewApp() (*App, error) {
	cfg, err := config.NewConfig()

	if err != nil {
		return nil, fmt.Errorf("Error in config loading")
	}

	logger, err := logger.NewLogger(cfg.AppName.Name)
	if err != nil {
		return nil, fmt.Errorf("Error in logger loading")
	}
	return &App{cfg, logger, server.NewServer(*cfg, logger)}, nil
}

func (a *App) Start() {
	a.logger.Info("Starting...")
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	go a.server.Start()

	a.logger.Info("Start")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	a.logger.Info("Shut down...")
	cancel()
}

func (a *App) Stop() {

}

func (a *App) newServer() *http.Server {
	r := chi.NewRouter()

	r.Use(myMiddleware.Logger(a.logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	raft := api.NewRaftApi(a.config.AppName.Name, a.logger)
	r.Get("/health", raft.Health)

	clientServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.config.Http.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		a.logger.Info("Starting client API server")
		if err := clientServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Failed to start client server")
		}
	}()

	return clientServer
}
