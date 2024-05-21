package editor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slog"
)

type EditorInterface interface {
	HandleRequest(ctx context.Context, req EditorRequest, callback func(outputFile string, err error))
}

type FfmpegEditor struct {
	BinaryPath string
}

func NewFFMpegEditor(binaryPath string) EditorInterface {
	return &FfmpegEditor{BinaryPath: binaryPath}
}

func (f *FfmpegEditor) HandleRequest(ctx context.Context, req EditorRequest, done func(outputFile string, err error)) {
	cmd, err := f.buildCommand(req)
	if err != nil {
		done("", err)
		return
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	result := make(chan error)

	go func(resultChannel chan error) {
		err = cmd.Run()
		resultChannel <- err
	}(result)

	select {
	case <-ctx.Done():
		err := cmd.Process.Kill()
		if err != nil {
			slog.Error("Failed to kill process", "error", err)
		}
		done("", fmt.Errorf("process killed"))
		return
	case err = <-result:
		if err != nil {
			done("", err)
			return
		}
	}

	done("output_1.mp4", nil)
}

func (f *FfmpegEditor) buildCommand(req EditorRequest) (*exec.Cmd, error) {
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
