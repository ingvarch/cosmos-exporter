package models

import "time"

// BlockInfo represents the structure of block information
type BlockInfo struct {
	Result struct {
		Block struct {
			Header struct {
				Height string    `json:"height"`
				Time   time.Time `json:"time"`
			} `json:"header"`
		} `json:"block"`
	} `json:"result"`
}

// NetInfo represents the structure of network information
type NetInfo struct {
	Result struct {
		NPeers string `json:"n_peers"`
		Peers  []struct {
			NodeInfo struct {
				Version string `json:"version"`
			} `json:"node_info"`
		} `json:"peers"`
	} `json:"result"`
}
