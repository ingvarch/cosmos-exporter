# Cosmos Exporter

A Prometheus exporter for collecting metrics from the Cosmos blockchain.

[![Go Report Card](https://goreportcard.com/badge/github.com/ingvarch/cosmos-exporter)](https://goreportcard.com/report/github.com/ingvarch/cosmos-exporter)
[![codecov](https://codecov.io/gh/ingvarch/cosmos-exporter/branch/main/graph/badge.svg)](https://codecov.io/gh/ingvarch/cosmos-exporter)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/ingvarch/cosmos-exporter)](https://github.com/ingvarch/cosmos-exporter/releases)
[![Build Status](https://github.com/ingvarch/cosmos-exporter/workflows/Run%20Tests%20and%20Coverage/badge.svg)](https://github.com/ingvarch/cosmos-exporter/actions)

## About

Cosmos Exporter is a tool for collecting and exporting metrics from the Cosmos blockchain in a Prometheus-compatible format. It allows monitoring of parameters such as block height, block time drift, number of connected peers, and their versions.

## Key Features

- Collection of block height metrics
- Monitoring of block time drift
- Tracking the number of connected peers
- Analysis of connected peer versions

## Usage

To run the exporter, execute:

```
cosmos-exporter
```

By default, metrics will be available at `http://localhost:8080/metrics`.

## Configuration

The exporter can be configured using environment variables:

- `COSMOS_NODE_ADDRESS`: The address of the Cosmos node (default: `http://localhost:26657`)
- `PORT`: The port on which the exporter serves metrics (default: `8080`)

## Building from Source

Ensure you have Go 1.23 or later installed, then run:

```
go get github.com/ingvarch/cosmos-exporter
cd $GOPATH/src/github.com/ingvarch/cosmos-exporter
go build
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
