package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// Try to connect to the Cosmos node
	if err := exporter.ConnectToCosmosNode(cfg); err != nil {
		slog.Error("Failed to connect to Cosmos node", "error", err)
		os.Exit(1)
	}

	// Create a new exporter with the Cosmos node address
	cosmosExporter := exporter.NewCosmosExporter(cfg.CosmosNodeAddress)

	// Register the exporter with Prometheus
	prometheus.MustRegister(cosmosExporter)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Set up handler for the /metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Create a new server with the ServeMux and configured timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Start HTTP server in a goroutine
	go func() {
		slog.Info("Starting Cosmos exporter", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}
