package main

import (
	"context"
	"sync"

	"github.com/douglasdgoulart/video-editor-api/pkg/api"
	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
)

func main() {
	ctx := context.Background()
	cfg := configuration.NewConfiguration()

	wg := sync.WaitGroup{}
	if cfg.Api.Enabled {
		wg.Add(1)
		go api.NewApi(cfg).Run(ctx)
	}

	wg.Wait()
}
