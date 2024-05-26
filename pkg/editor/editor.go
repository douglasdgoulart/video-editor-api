package editor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"regexp"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
	"github.com/douglasdgoulart/video-editor-api/pkg/request"
	"github.com/google/uuid"
)

type EditorInterface interface {
	HandleRequest(ctx context.Context, req request.EditorRequest) ([]string, error)
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

func (f *FfmpegEditor) HandleRequest(ctx context.Context, req request.EditorRequest) (output []string, err error) {
	outputPattern := getOutputPath(req.Output.FilePattern)
	outputPath := filepath.Dir(outputPattern)
	req.Output.FilePattern = outputPattern

	cmd, err := f.buildCommand(req)
	if err != nil {
		return
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
		err = cmd.Process.Kill()
		if err != nil {
			f.logger.Error("Failed to kill process", "error", err)
		}
		err = fmt.Errorf("process killed")
		return
	case err = <-result:
		if err != nil {
			return
		}
	}
	f.logger.Info("Command finished successfully")

	output, err = getFilesInDirectory(outputPath)
	return
}

func getFilesInDirectory(directory string) ([]string, error) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func getOutputPath(outputFilePattern string) string {
	outputPattern := "output"
	outputFileExtention := strings.ToLower(outputFilePattern[strings.LastIndex(outputFilePattern, ".")+1:])
	placeholderRegex := regexp.MustCompile(`%[0-9]{2}d`)
	if placeholderRegex.MatchString(outputFilePattern) {
		outputPattern = placeholderRegex.FindString(outputFilePattern)
	}

	outputFilePattern = fmt.Sprintf("%s/%s/%s.%s", os.TempDir(), uuid.New().String(), outputPattern, outputFileExtention)
	_ = os.MkdirAll(filepath.Dir(outputFilePattern), os.ModePerm)

	return outputFilePattern
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

	args = append(args, req.Output.FilePattern)

	cmd := exec.Command(f.BinaryPath, args...)

	f.logger.Info("Running command", "command", strings.Join(cmd.Args, " "))
	return cmd, nil
}
