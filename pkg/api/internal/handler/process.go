package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
	"github.com/douglasdgoulart/video-editor-api/pkg/event/emitter"
	"github.com/douglasdgoulart/video-editor-api/pkg/request"
	"github.com/douglasdgoulart/video-editor-api/pkg/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProcessHandler struct {
	logger    *slog.Logger
	emitter   emitter.EventEmitter
	inputPath string
}

func NewProcessHandler(cfg *configuration.Configuration) *ProcessHandler {
	var eventEmitter emitter.EventEmitter
	if cfg.Kafka.Enabled {
		eventEmitter = emitter.NewKafkaEmitter(&cfg.Kafka.KafkaProducerConfig)
	} else {
		eventEmitter = emitter.NewInternalQueueEmitter(cfg)
	}

	return &ProcessHandler{
		logger:    cfg.Logger.WithGroup("process_handler"),
		emitter:   eventEmitter,
		inputPath: cfg.InputPath,
	}
}

func (ph *ProcessHandler) Handler(c echo.Context) error {
	var request request.EditorRequest
	eventJson := c.FormValue("event")
	err := json.Unmarshal([]byte(eventJson), &request)
	if err != nil {
		ph.logger.Error("Failed to decode request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err = validator.ValidateRequiredFields(request)
	if err != nil {
		ph.logger.Error("Failed to validate request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	file, err := c.FormFile("file")
	if err != nil {
		ph.logger.Error("Failed to get file from request", "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	fileLocation, err := ph.downloadFile(file)
	if err != nil {
		ph.logger.Error("Failed to download file", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	request.Input.UploadedFilePath = fileLocation

	eventId := uuid.New().String()
	err = ph.emitter.Send(c.Request().Context(), event.Event{
		Id:            eventId,
		EditorRequest: request,
	})
	if err != nil {
		ph.logger.Error("Failed to send event", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "processing request", "id": eventId})
}

func (ph *ProcessHandler) downloadFile(f *multipart.FileHeader) (string, error) {
	src, err := f.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileExtention := f.Filename[strings.LastIndex(f.Filename, ".")+1:]
	dst, err := os.Create(fmt.Sprintf("%s/%s.%s", ph.inputPath, uuid.New().String(), fileExtention))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dst.Name(), nil
}
