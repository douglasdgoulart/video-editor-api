package editor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/request"
	"golang.org/x/exp/slog"
)

type EditorInterface interface {
	HandleRequest(ctx context.Context, req request.EditorRequest) (string, error)
}

type FfmpegEditor struct {
	BinaryPath string
}

func NewFFMpegEditor(cfg *configuration.Configuration) EditorInterface {
	return &FfmpegEditor{BinaryPath: cfg.Ffmpeg.Path}
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
		slog.Info("Running command", "command", strings.Join(cmd.Args, " "))
		err := cmd.Run()
		slog.Info("Command finished", "error", err)
		resultChannel <- err
	}(result)

	slog.Info("Waiting for command to finish")
	select {
	case <-ctx.Done():
		err := cmd.Process.Kill()
		if err != nil {
			slog.Error("Failed to kill process", "error", err)
		}
		return "", fmt.Errorf("process killed")
	case err = <-result:
		if err != nil {
			return "", err
		}
	}
	slog.Info("Command finished successfully")

	return req.Output.FilePattern, nil
}

func (f *FfmpegEditor) buildCommand(req request.EditorRequest) (*exec.Cmd, error) {
	var inputFilePath string

	// Determine input file path
	if req.Input.FileURL != "" {
		inputFilePath = req.Input.FileURL // Use URL directly
	} else if req.Input.UploadedFilePath != "" {
		inputFilePath = req.Input.UploadedFilePath
	} else {
		return nil, fmt.Errorf("no valid input file provided")
	}

	// Build the FFmpeg command
	outputFilePattern := req.Output.FilePattern
	if outputFilePattern == "" {
		outputFilePattern = "output.jpg"
	}

	// Construct the FFmpeg command arguments
	args := []string{"-y"} // Overwrite output files without asking

	if req.StartTime != "" {
		args = append(args, "-ss", req.StartTime)
	}

	args = append(args, "-i", inputFilePath)

	if len(req.Filters) > 0 {
		var filterStrings []string

		for name, options := range req.Filters {
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

	slog.Info("Running command", "command", strings.Join(cmd.Args, " "))
	return cmd, nil
}
