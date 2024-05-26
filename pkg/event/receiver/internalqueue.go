package receiver

import (
	"context"
	"log/slog"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
)

type InternalQueueEventReceiver struct {
	queue  <-chan event.Event
	logger *slog.Logger
}

func NewInternalQueueEventReceiver(cfg *configuration.Configuration) EventReceiver {
	return &InternalQueueEventReceiver{
		queue:  cfg.InternalQueue,
		logger: cfg.Logger.WithGroup("internal_queue_event_receiver"),
	}
}

func (i *InternalQueueEventReceiver) Receive(ctx context.Context, handler func(event *event.Event) error) {
	for {
		select {
		case e := <-i.queue:
			err := handler(&e)
			if err != nil {
				i.logger.Error("error handling event", "error", err, "event", e)
			}
		case <-ctx.Done():
			return
		}
	}
}
