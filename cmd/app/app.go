package main

import (
	"context"
	"sync"

	"github.com/douglasdgoulart/video-editor-api/pkg/api"
	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/job"
)

func main() {
	ctx := context.Background()
	cfg := configuration.NewConfiguration()

	wg := sync.WaitGroup{}
	if cfg.Api.Enabled {
		wg.Add(1)
		go api.NewApi(cfg).Run(ctx)
	}

	if cfg.Job.Enabled {
		for jobId := range cfg.Job.Workers {
			wg.Add(1)
			cfg.Logger.Info("Starting job", "job_id", jobId)
			go job.NewJob(cfg, jobId).Run(ctx)
		}
	}

	wg.Wait()
}
