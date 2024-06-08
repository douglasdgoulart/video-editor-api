package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
)

type Server struct {
	server *http.Server
	logger *slog.Logger
}

func NewServer(cfg *configuration.Configuration) *Server {
	api := NewApi(cfg)

	// Declare Server config
	server := &http.Server{
		Addr:         cfg.Api.Port,
		Handler:      api.GetHandler(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{
		server: server,
		logger: cfg.Logger.WithGroup("server"),
	}
}

func (s Server) Run(parentCtx context.Context) {
	ctx, stop := signal.NotifyContext(parentCtx, os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("shutting down the server", "error", err)
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	s.logger.Info("shutting down the api server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("shutting down the server", "error", err)
		panic(err)
	}
	s.logger.Info("api server stopped")
}
