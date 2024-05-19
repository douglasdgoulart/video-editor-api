package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
)

type Api struct {
	e      *echo.Echo
	port   string
	logger *slog.Logger
}

func NewApi(port string, logger *slog.Logger) *Api {
	e := echo.New()
	e.Use(slogecho.New(logger))
	registerHealthCheckRoute(e)

	return &Api{
		e:      e,
		port:   port,
		logger: logger,
	}
}

func registerHealthCheckRoute(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}

func (a Api) Run(parentCtx context.Context) {
	ctx, stop := signal.NotifyContext(parentCtx, os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := a.e.Start(a.port); err != nil && err != http.ErrServerClosed {
			a.e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		a.e.Logger.Fatal(err)
	}
	a.e.Logger.Fatal(a.e.Start(a.port))
}
