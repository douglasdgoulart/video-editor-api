package api

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/douglasdgoulart/video-editor-api/pkg/configuration"
)

func TestApi_Run(t *testing.T) {
	t.Run("Given nothing, when the api is running it should return OK at health route", func(t *testing.T) {
		cfg := &configuration.Configuration{
			Logger: slog.Default(),
		}
		api := NewApi(cfg)

		server := httptest.NewServer(api.GetHandler())
		defer server.Close()

		resp, err := http.Get(fmt.Sprintf("%s/health", server.URL))
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
