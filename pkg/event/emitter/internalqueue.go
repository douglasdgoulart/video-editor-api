package emitter

import (
	"context"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
)

type InternalQueueEventEmitter struct {
	queue chan<- event.Event
}

func NewInternalQueueEmitter(cfg *configuration.Configuration) EventEmitter {
	return &InternalQueueEventEmitter{
		queue: cfg.InternalQueue,
	}
}

func (i *InternalQueueEventEmitter) Send(ctx context.Context, e event.Event) error {
	go func(queue chan<- event.Event, e event.Event) {
		i.queue <- e
	}(i.queue, e)
	return nil
}
