package job

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/editor"
	"github.com/douglasdgoulart/video-editor-api/pkg/event"
	"github.com/douglasdgoulart/video-editor-api/pkg/event/receiver"
)

type JobInterface interface {
	Run(ctx context.Context)
}

type Job struct {
	eventReceiver receiver.EventReceiver
	editor        editor.EditorInterface
	logger        *slog.Logger
	apiHost       string
	apiPort       string
	outputPath    string
}

func NewJob(cfg *configuration.Configuration, jobId int) JobInterface {
	var eventReceiver receiver.EventReceiver
	if cfg.Kafka.Enabled {
		eventReceiver = receiver.NewKafkaEventReceiver(cfg)
	} else {
		eventReceiver = receiver.NewInternalQueueEventReceiver(cfg)
	}

	editor := editor.NewFFMpegEditor(cfg)
	logger := cfg.Logger.WithGroup(fmt.Sprintf("job_%d", jobId))
	return &Job{
		eventReceiver: eventReceiver,
		editor:        editor,
		logger:        logger,
		apiHost:       cfg.Api.Host,
		apiPort:       cfg.Api.Port,
		outputPath:    cfg.OutputPath,
	}
}

func (j *Job) Run(ctx context.Context) {
	j.eventReceiver.Receive(ctx, j.handleEvent(ctx))
	j.logger.Info("job stoped")
}

func (j *Job) handleEvent(ctx context.Context) func(event *event.Event) error {
	return func(event *event.Event) error {
		otputFileLocation, err := j.editor.HandleRequest(ctx, event.EditorRequest)
		if err != nil {
			j.logger.Error("error handling event", "error", err)
			err := j.callWebhook(event, otputFileLocation, err)
			if err != nil {
				j.logger.Error("error calling webhook", "error", err)
			}
			return err
		}
		return j.callWebhook(event, otputFileLocation, err)
	}
}

type WebhookResponse struct {
	Status        string   `json:"status"`
	Id            string   `json:"id"`
	FileLocations []string `json:"file_location,omitempty"`
	ErrorMsg      string   `json:"error_msg,omitempty"`
}

func (j *Job) getFileLocationURL(fileLocations []string, host string, port string) []string {
	var urls []string
	if port == "" {
		port = ":80"
	}
	method := "http"
	if host != "localhost" {
		method = "https"
	}
	for _, fileLocation := range fileLocations {
		fileLocation = strings.Replace(fileLocation, fmt.Sprintf("%s/", j.outputPath), "", 1)
		urls = append(urls, fmt.Sprintf("%s://%s%s/files/%s", method, host, port, fileLocation))
	}
	return urls
}

func (j *Job) callWebhook(event *event.Event, outputFilesLocation []string, inputErr error) error {
	outputFileLocationsURL := j.getFileLocationURL(outputFilesLocation, j.apiHost, j.apiPort)
	if event.EditorRequest.Output.WebhookURL == "" {
		return nil
	}
	url := event.EditorRequest.Output.WebhookURL

	status := "success"
	errMsg := ""
	if inputErr != nil {
		status = "error"
		errMsg = inputErr.Error()
	}
	payload := WebhookResponse{
		Status:        status,
		Id:            event.Id,
		FileLocations: outputFileLocationsURL,
		ErrorMsg:      errMsg,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing request:", err)
		return err
	}
	defer resp.Body.Close()

	return nil
}
