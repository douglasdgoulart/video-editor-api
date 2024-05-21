package emmiter

import (
	"context"

	"github.com/douglasdgoulart/video-editor-api/pkg/event"
)

type EventEmitter interface {
	Send(ctx context.Context, event event.Event) error
}
