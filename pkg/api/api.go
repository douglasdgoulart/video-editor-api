package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
	"github.com/douglasdgoulart/video-editor-api/pkg/event/emitter"
	"github.com/douglasdgoulart/video-editor-api/pkg/request"
	"github.com/douglasdgoulart/video-editor-api/pkg/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	slogecho "github.com/samber/slog-echo"
)

type ApiInterface interface {
	Run(ctx context.Context)
}

type Api struct {
	e          *echo.Echo
	port       string
	logger     *slog.Logger
	emitter    emitter.EventEmitter
	outputPath string
	inputPath  string
}

func NewApi(cfg *configuration.Configuration) ApiInterface {
	e := echo.New()
	e.Use(slogecho.New(cfg.Logger))

	var eventEmitter emitter.EventEmitter
	if cfg.Kafka.Enabled {
		eventEmitter = emitter.NewKafkaEmitter(&cfg.Kafka.KafkaProducerConfig)
	} else {
		eventEmitter = emitter.NewInternalQueueEmitter(cfg)
	}

	logger := cfg.Logger.WithGroup("api")
	api := &Api{
		e:          e,
		port:       cfg.Api.Port,
		logger:     logger,
		emitter:    eventEmitter,
		outputPath: cfg.OutputPath,
		inputPath:  cfg.InputPath,
	}

	api.registerHealthCheckRoute()
	api.registerProcessRoute()
	api.registerStaticFiles()

	return api
}

func (a *Api) registerStaticFiles() {
	a.e.Static("/files", a.outputPath)
}

func (a *Api) registerHealthCheckRoute() {
	a.e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}

func (a *Api) registerProcessRoute() {
	a.e.POST("/process", func(c echo.Context) error {
		var request request.EditorRequest
		eventJson := c.FormValue("event")
		err := json.Unmarshal([]byte(eventJson), &request)
		if err != nil {
			a.logger.Error("Failed to decode request", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		err = validator.ValidateRequiredFields(request)
		if err != nil {
			a.logger.Error("Failed to validate request", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		file, err := c.FormFile("file")
		if err != nil {
			a.logger.Error("Failed to get file from request", "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		fileLocation, err := a.downloadFile(file)
		if err != nil {
			a.logger.Error("Failed to download file", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
		request.Input.UploadedFilePath = fileLocation

		eventId := uuid.New().String()
		err = a.emitter.Send(c.Request().Context(), event.Event{
			Id:            eventId,
			EditorRequest: request,
		})
		if err != nil {
			a.logger.Error("Failed to send event", "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "processing request", "id": eventId})
	})
}

func (a *Api) downloadFile(f *multipart.FileHeader) (string, error) {
	src, err := f.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileExtention := f.Filename[strings.LastIndex(f.Filename, ".")+1:]
	dst, err := os.Create(fmt.Sprintf("%s/%s.%s", a.inputPath, uuid.New().String(), fileExtention))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dst.Name(), nil
}

func (a Api) Run(parentCtx context.Context) {
	ctx, stop := signal.NotifyContext(parentCtx, os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := a.e.Start(a.port); err != nil && err != http.ErrServerClosed {
			a.logger.Error("shutting down the server", "error", err)
			panic(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		a.logger.Error("shutting down the server", "error", err)
		panic(err)
	}
	a.e.Logger.Fatal(a.e.Start(a.port))
}
