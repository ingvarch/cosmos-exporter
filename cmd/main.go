package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ingvarch/cosmos-exporter/internal/config"
	"github.com/ingvarch/cosmos-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Set up the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.New()

	// Create a new exporter with the Cosmos node address
	cosmosExporter := exporter.NewCosmosExporter(cfg.CosmosNodeAddress)

	// Register the exporter with Prometheus
	prometheus.MustRegister(cosmosExporter)

	// Set up handler for the /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start HTTP server on the specified port
	slog.Info("Starting Cosmos exporter", "port", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
