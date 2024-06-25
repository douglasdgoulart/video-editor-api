package handler

import (
	"log/slog"
	"net/http"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/labstack/echo/v4"
	"github.com/twmb/franz-go/pkg/kgo"
)

type HealthHandler struct {
	cl     *kgo.Client
	logger *slog.Logger
}

func NewHealthHandler(cfg *configuration.Configuration) *HealthHandler {
	logger := cfg.Logger.WithGroup("health-handler")
	var cl *kgo.Client
	var err error
	if cfg.Kafka.Enabled {
		brokers := []string{}
		if cfg.Api.Enabled {
			brokers = append(brokers, cfg.Kafka.KafkaProducerConfig.Brokers...)
		}
		if cfg.Job.Enabled {
			brokers = append(brokers, cfg.Kafka.KafkaConsumerConfig.Brokers...)
		}
		cl, err = kgo.NewClient(
			kgo.SeedBrokers(brokers...),
		)
		if err != nil {
			logger.Error("error creating kafka client", "error", err)
			panic(err)
		}
	}

	return &HealthHandler{
		cl:     cl,
		logger: logger,
	}
}

func (h *HealthHandler) HealthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (h *HealthHandler) ReadyHandler(c echo.Context) error {
	if h.cl != nil {
		err := h.cl.Ping(c.Request().Context())
		if err != nil {
			h.logger.Error("error pinging kafka", "error", err)
			return c.String(http.StatusInternalServerError, "error")
		}
	}

	return c.String(http.StatusOK, "OK")
}
