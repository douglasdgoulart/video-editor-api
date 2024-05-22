package receiver

import (
	"context"

	"github.com/douglasdgoulart/video-editor-api/pkg/event"
)

type EventReceiver interface {
	Receive(ctx context.Context, handler func(event *event.Event) error)
}
