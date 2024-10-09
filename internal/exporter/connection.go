package exporter

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ingvarch/cosmos-exporter/internal/config"
)

// ConnectToCosmosNode attempts to connect to the Cosmos node with retries
func ConnectToCosmosNode(cfg *config.Config) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for i := 0; i < cfg.MaxRetries; i++ {
		resp, err := client.Get(cfg.CosmosNodeAddress + "/status")
		if err != nil {
			slog.Warn("Failed to connect to Cosmos node", "attempt", i+1, "error", err)
			time.Sleep(cfg.RetryDelay)
			continue
		}

		// Ensure the body is always closed
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				slog.Warn("Error closing response body", "error", cerr)
			}
		}()

		// Read and discard the body to reuse the connection
		_, err = io.Copy(io.Discard, resp.Body)
		if err != nil {
			slog.Warn("Error reading response body", "error", err)
		}

		if resp.StatusCode == http.StatusOK {
			slog.Info("Successfully connected to Cosmos node", "address", cfg.CosmosNodeAddress)
			return nil
		}

		slog.Warn("Received non-OK status code", "attempt", i+1, "statusCode", resp.StatusCode)
		time.Sleep(cfg.RetryDelay)
	}

	return fmt.Errorf("failed to connect to Cosmos node after %d attempts", cfg.MaxRetries)
}
