package exporter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ingvarch/cosmos-exporter/internal/config"
)

func TestConnectToCosmosNode(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/status" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a test configuration
	cfg := &config.Config{
		CosmosNodeAddress: server.URL,
		MaxRetries:        3,
		RetryDelay:        100 * time.Millisecond,
	}

	// Test successful connection
	err := ConnectToCosmosNode(cfg)
	if err != nil {
		t.Errorf("Expected successful connection, got error: %v", err)
	}

	// Test failed connection
	cfg.CosmosNodeAddress = "http://invalid-address"
	err = ConnectToCosmosNode(cfg)
	if err == nil {
		t.Error("Expected error for invalid address, got nil")
	}
}
