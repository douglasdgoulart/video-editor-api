package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/douglasdgoulart/video-editor-api/pkg/api"
	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/job"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cfg := configuration.NewConfiguration()

	wg := sync.WaitGroup{}
	wg.Add(1)
	cfg.Logger.Info("Starting API server")
	go func() {
		defer wg.Done()
		api.NewServer(cfg).Run(ctx)
	}()

	if cfg.Job.Enabled {
		for jobId := range cfg.Job.Workers {
			wg.Add(1)
			cfg.Logger.Info("Starting job", "job_id", jobId)
			go func(jobId int) {
				defer wg.Done()
				job.NewJob(cfg, jobId).Run(ctx)
			}(jobId)
		}
	}

	wg.Wait()
}
