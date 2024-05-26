package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
)

func TestApi_Run(t *testing.T) {
	t.Run("Given nothing, when the api is running it should return OK at health route", func(t *testing.T) {
		cfg := &configuration.Configuration{
			Logger: slog.Default(),
			Api: configuration.ApiConfig{
				Port: ":0",
			},
		}
		api := NewApi(cfg)

		go api.Run(context.Background())
		time.Sleep(1 * time.Second)

		port := api.(*Api).e.Listener.Addr().String()[strings.LastIndex(api.(*Api).e.Listener.Addr().String(), ":"):]
		resp, err := http.Get(fmt.Sprintf("http://localhost%s/health", port))
		if err != nil {
			t.Fatalf("Failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}

		if string(body) != "OK" {
			t.Errorf("Expected body 'OK'; got %v", string(body))
		}
	})

}
