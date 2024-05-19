package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/douglasdgoulart/video-editor-api/pkg/api"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.Background()
	api.NewApi(":8080", logger).Run(ctx)
}
