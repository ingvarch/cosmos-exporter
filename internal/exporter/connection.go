package exporter

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ingvarch/cosmos-exporter/internal/config"
)

// ConnectToCosmosNode attempts to connect to the Cosmos node with retries
func ConnectToCosmosNode(cfg *config.Config) error {
	var lastErr error
	for i := 0; i < cfg.MaxRetries; i++ {
		resp, err := http.Get(cfg.CosmosNodeAddress + "/status")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				slog.Info("Successfully connected to Cosmos node", "address", cfg.CosmosNodeAddress)
				return nil
			}
			lastErr = fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		slog.Warn("Failed to connect to Cosmos node, retrying...", "attempt", i+1, "error", lastErr)
		time.Sleep(cfg.RetryDelay)
	}
	return fmt.Errorf("failed to connect to Cosmos node after %d attempts: %v", cfg.MaxRetries, lastErr)
}
