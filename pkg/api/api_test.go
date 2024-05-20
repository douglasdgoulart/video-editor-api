package api

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func TestApi_Run(t *testing.T) {
	t.Run("Given nothing, when the api is running it should return OK at health route", func(t *testing.T) {
		api := NewApi(":8081", slog.Default())

		go api.Run(context.Background())
		time.Sleep(1 * time.Second)

		resp, err := http.Get("http://localhost:8081/health")
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
