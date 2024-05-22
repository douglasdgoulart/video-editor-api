package event

import (
	"context"

	"github.com/douglasdgoulart/video-editor-api/pkg/editor"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KgoClient interface {
	ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults
	PollFetches(ctx context.Context) kgo.Fetches
}

type Event struct {
	Id            string               `json:"id"`
	EditorRequest editor.EditorRequest `json:"editor_request"`
}
