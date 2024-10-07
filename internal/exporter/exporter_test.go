package exporter

import (
	"testing"

	"github.com/ingvarch/cosmos-exporter/internal/config"
)

func TestNewCosmosExporter(t *testing.T) {
	cfg := &config.Config{
		CosmosNodeAddress: "http://localhost:26657",
	}
	exporter := NewCosmosExporter(cfg.CosmosNodeAddress)
	if exporter == nil {
		t.Error("NewCosmosExporter returned nil")
	}
	if exporter.cosmosNodeAddress != cfg.CosmosNodeAddress {
		t.Errorf("Expected CosmosNodeAddress to be %s, but got %s", cfg.CosmosNodeAddress, exporter.cosmosNodeAddress)
	}
}
