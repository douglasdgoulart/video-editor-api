package main

import (
	"context"

	"github.com/douglasdgoulart/video-editor-api/pkg/api"
	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
)

func main() {
	ctx := context.Background()
	cfg := configuration.NewConfiguration()
	api.NewApi(cfg).Run(ctx)
}
