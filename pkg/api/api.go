package api

import (
	"log/slog"
	"net/http"

	"github.com/douglasdgoulart/video-editor-api/pkg/api/internal/handler"
	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
)

type ApiInterface interface {
	GetHandler() http.Handler
}

type Api struct {
	e          *echo.Echo
	logger     *slog.Logger
	outputPath string
}

func NewApi(cfg *configuration.Configuration) ApiInterface {
	e := echo.New()
	e.Use(slogecho.New(cfg.Logger))

	logger := cfg.Logger.WithGroup("api")
	api := &Api{
		e:          e,
		logger:     logger,
		outputPath: cfg.OutputPath,
	}

	processHandler := handler.NewProcessHandler(cfg)
	healthHandler := handler.NewHealthHandler()

	api.e.GET("/health", healthHandler.Handler)
	api.e.POST("/process", processHandler.Handler)
	api.e.Static("/files", api.outputPath)

	return api
}

func (a *Api) GetHandler() http.Handler {
	return a.e
}
