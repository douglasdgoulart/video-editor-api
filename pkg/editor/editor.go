package editor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/request"
	"github.com/google/uuid"
)

type EditorInterface interface {
	HandleRequest(ctx context.Context, req request.EditorRequest) (string, error)
}

type FfmpegEditor struct {
	BinaryPath string
	logger     *slog.Logger
}

func NewFFMpegEditor(cfg *configuration.Configuration) EditorInterface {
	return &FfmpegEditor{
		BinaryPath: cfg.Ffmpeg.Path,
		logger:     cfg.Logger.WithGroup("ffmpeg_editor"),
	}

}

func (f *FfmpegEditor) HandleRequest(ctx context.Context, req request.EditorRequest) (string, error) {
	cmd, err := f.buildCommand(req)
	if err != nil {
		return "", err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	result := make(chan error)

	go func(resultChannel chan<- error) {
		f.logger.Info("Running command", "command", strings.Join(cmd.Args, " "))
		err := cmd.Run()
		f.logger.Info("Command finished", "error", err)
		resultChannel <- err
	}(result)

	f.logger.Info("Waiting for command to finish")
	select {
	case <-ctx.Done():
		err := cmd.Process.Kill()
		if err != nil {
			f.logger.Error("Failed to kill process", "error", err)
		}
		return "", fmt.Errorf("process killed")
	case err = <-result:
		if err != nil {
			return "", err
		}
	}
	f.logger.Info("Command finished successfully")

	return req.Output.FilePattern, nil
}

func (f *FfmpegEditor) buildCommand(req request.EditorRequest) (*exec.Cmd, error) {
	var inputFilePath string

	if req.Input.FileURL != "" {
		inputFilePath = req.Input.FileURL
	} else if req.Input.UploadedFilePath != "" {
		inputFilePath = req.Input.UploadedFilePath
	} else {
		return nil, fmt.Errorf("no valid input file provided")
	}

	outputFilePattern := req.Output.FilePattern
	inputExtention := strings.ToLower(inputFilePath[strings.LastIndex(inputFilePath, ".")+1:])
	if outputFilePattern == "" {
		outputFilePattern = fmt.Sprintf("%s/%s.%s", os.TempDir(), uuid.New().String(), inputExtention)
	}

	args := []string{"-y"}

	if req.StartTime != "" {
		args = append(args, "-ss", req.StartTime)
	}

	args = append(args, "-i", inputFilePath)

	if len(req.Filters) > 0 {
		var filterStrings []string

		for name, options := range req.Filters {
			if options == "" {
				filterStrings = append(filterStrings, name)
				continue
			}

			filterString := fmt.Sprintf("%s=%s", name, options)
			filterStrings = append(filterStrings, filterString)
		}

		filterGraph := strings.Join(filterStrings, ",")
		args = append(args, "-vf", filterGraph)
	}
	if req.Frames != "" {
		args = append(args, "-frames:v", req.Frames)
	}

	if req.ExtraOptions != "" {
		extraArgs := strings.Split(req.ExtraOptions, " ")
		args = append(args, extraArgs...)
	}

	args = append(args, outputFilePattern)

	cmd := exec.Command(f.BinaryPath, args...)

	f.logger.Info("Running command", "command", strings.Join(cmd.Args, " "))
	return cmd, nil
}
