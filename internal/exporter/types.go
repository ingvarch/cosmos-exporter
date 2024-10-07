package exporter

import "github.com/prometheus/client_golang/prometheus"

// CosmosExporter represents the structure of the Cosmos exporter
type CosmosExporter struct {
	highestBlock      prometheus.Gauge
	blockTimeDrift    prometheus.Gauge
	connectedPeers    prometheus.Gauge
	peersByVersion    *prometheus.GaugeVec
	cosmosNodeAddress string
}
