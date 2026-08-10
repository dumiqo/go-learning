package server

import (
	"RaftLike/config"
	"RaftLike/pkg/logger"
	myMiddleware "RaftLike/pkg/middleware"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	logger *logger.Logger
	srv    *http.Server
}

func NewServer(cfg config.Config, logger *logger.Logger) *Server {

	r := chi.NewRouter()

	r.Use(myMiddleware.Logger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	httpServer := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &Server{logger, &httpServer}
}

func (s *Server) Start() {
	s.logger.Info("Starting client API server")
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("Failed to start client server")
	}
}
