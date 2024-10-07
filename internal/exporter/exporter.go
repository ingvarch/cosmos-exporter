package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ingvarch/cosmos-exporter/internal/models"
	"github.com/prometheus/client_golang/prometheus"
)

// NewCosmosExporter creates a new instance of CosmosExporter
func NewCosmosExporter(cosmosNodeAddress string) *CosmosExporter {
	return &CosmosExporter{
		highestBlock: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cosmos_highest_block",
			Help: "Highest block number in the Cosmos blockchain",
		}),
		blockTimeDrift: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cosmos_block_time_drift",
			Help: "Current block time drift in seconds",
		}),
		connectedPeers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cosmos_connected_peers",
			Help: "Number of connected peers to the Cosmos node",
		}),
		peersByVersion: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "cosmos_peers_by_version",
				Help: "Number of peers by version",
			},
			[]string{"version"},
		),
		cosmosNodeAddress: cosmosNodeAddress,
	}
}

// Describe sends the descriptions of all metrics to the provided channel
func (e *CosmosExporter) Describe(ch chan<- *prometheus.Desc) {
	e.highestBlock.Describe(ch)
	e.blockTimeDrift.Describe(ch)
	e.connectedPeers.Describe(ch)
	e.peersByVersion.Describe(ch)
}

// Collect fetches the values for all metrics and sends them to the provided channel
func (e *CosmosExporter) Collect(ch chan<- prometheus.Metric) {
	e.updateMetrics()
	e.highestBlock.Collect(ch)
	e.blockTimeDrift.Collect(ch)
	e.connectedPeers.Collect(ch)
	e.peersByVersion.Collect(ch)
}

// updateMetrics updates all metrics by fetching current data from the Cosmos node
func (e *CosmosExporter) updateMetrics() {
	blockInfo, err := e.getLatestBlockInfo()
	if err != nil {
		slog.Error("Error getting latest block info", "error", err)
	} else {
		height, err := strconv.ParseFloat(blockInfo.Result.Block.Header.Height, 64)
		if err != nil {
			slog.Error("Error parsing block height", "error", err)
		} else {
			e.highestBlock.Set(height)
		}

		e.blockTimeDrift.Set(time.Since(blockInfo.Result.Block.Header.Time).Seconds())
	}

	netInfo, err := e.getNetInfo()
	if err != nil {
		slog.Error("Error getting net info", "error", err)
	} else {
		nPeers, _ := strconv.Atoi(netInfo.Result.NPeers)
		e.connectedPeers.Set(float64(nPeers))

		versionCounts := make(map[string]int)
		for _, peer := range netInfo.Result.Peers {
			version := peer.NodeInfo.Version
			versionCounts[version]++
		}
		for version, count := range versionCounts {
			e.peersByVersion.WithLabelValues(version).Set(float64(count))
		}
	}
}

// getLatestBlockInfo fetches the latest block information from the Cosmos node
func (e *CosmosExporter) getLatestBlockInfo() (models.BlockInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/block", e.cosmosNodeAddress))
	if err != nil {
		return models.BlockInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.BlockInfo{}, err
	}

	var blockInfo models.BlockInfo
	err = json.Unmarshal(body, &blockInfo)
	if err != nil {
		return models.BlockInfo{}, err
	}

	return blockInfo, nil
}

// getNetInfo fetches network information from the Cosmos node
func (e *CosmosExporter) getNetInfo() (models.NetInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/net_info", e.cosmosNodeAddress))
	if err != nil {
		return models.NetInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.NetInfo{}, err
	}

	var netInfo models.NetInfo
	err = json.Unmarshal(body, &netInfo)
	if err != nil {
		return models.NetInfo{}, err
	}

	return netInfo, nil
}
