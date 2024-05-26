package editor

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/douglasdgoulart/video-editor-api/pkg/request"
)

var ffmpegLocation = "../../bin/ffmpeg/ffmpeg"

func TestFfmpegEditor_buildCommand(t *testing.T) {
	t.Run("buildCommand should return a valid command", func(t *testing.T) {
		editor := NewFFMpegEditor(ffmpegLocation)

		req := request.EditorRequest{
			Input: request.Input{
				UploadedFilePath: "../../input.mp4",
			},
			Output: request.Output{
				FilePattern: "thumbnail.jpg",
			},
			ExtraOptions: "-vf \"thumbnail,scale=640:480\" -frames:v 1",
		}

		cmd, err := editor.(*FfmpegEditor).buildCommand(req)
		if err != nil {
			t.Fatalf("Failed to build command: %v", err)
		}

		expectedArgs := []string{
			"-y",
			"-i", "../../input.mp4",
			"-vf", "\"thumbnail,scale=640:480\"",
			"-frames:v", "1",
			"thumbnail.jpg",
		}

		if cmd.Path != ffmpegLocation {
			t.Errorf("Expected path '%s'; got %v", ffmpegLocation, cmd.Path)
		}

		for i, arg := range cmd.Args[1:] {
			if arg != expectedArgs[i] {
				t.Errorf("Expected arg '%v'; got %v", expectedArgs[i], arg)
			}
		}
	})
}

func TestFfmpegEditor_extractThumbnail(t *testing.T) {
	t.Run("extractThumbnail should return a valid command", func(t *testing.T) {
		editor := NewFFMpegEditor(ffmpegLocation)
		outputFile := fmt.Sprintf("/tmp/thumbnail_%d.jpg", rand.Int())

		filters := map[string]string{
			"scale": "-1:100",
		}

		req := request.EditorRequest{
			Input: request.Input{
				UploadedFilePath: "../../internal/testdata/testsrc.mp4",
			},
			Output: request.Output{
				FilePattern: outputFile,
			},
			StartTime:    "00:00:05.0",
			Filters:      filters,
			Frames:       "1",
			ExtraOptions: "",
		}

		_, err := editor.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatalf("Failed to extract thumbnail: %v", err)
		}

		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Fatalf("Expected output file '%s' to exist", outputFile)
		}
	})
}
