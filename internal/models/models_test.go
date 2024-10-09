package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBlockInfo(t *testing.T) {
	jsonData := `{
		"result": {
			"block": {
				"header": {
					"height": "100",
					"time": "2023-04-15T12:00:00Z"
				}
			}
		}
	}`

	var blockInfo BlockInfo
	err := json.Unmarshal([]byte(jsonData), &blockInfo)
	if err != nil {
		t.Fatalf("Failed to unmarshal BlockInfo: %v", err)
	}

	if blockInfo.Result.Block.Header.Height != "100" {
		t.Errorf("Expected height to be 100, got %s", blockInfo.Result.Block.Header.Height)
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2023-04-15T12:00:00Z")
	if !blockInfo.Result.Block.Header.Time.Equal(expectedTime) {
		t.Errorf("Expected time to be %v, got %v", expectedTime, blockInfo.Result.Block.Header.Time)
	}
}

func TestNetInfo(t *testing.T) {
	jsonData := `{
		"result": {
			"n_peers": "3",
			"peers": [
				{"node_info": {"version": "1.0.0"}},
				{"node_info": {"version": "1.1.0"}},
				{"node_info": {"version": "1.0.0"}}
			]
		}
	}`

	var netInfo NetInfo
	err := json.Unmarshal([]byte(jsonData), &netInfo)
	if err != nil {
		t.Fatalf("Failed to unmarshal NetInfo: %v", err)
	}

	if netInfo.Result.NPeers != "3" {
		t.Errorf("Expected n_peers to be 3, got %s", netInfo.Result.NPeers)
	}

	if len(netInfo.Result.Peers) != 3 {
		t.Errorf("Expected 3 peers, got %d", len(netInfo.Result.Peers))
	}

	versionCounts := make(map[string]int)
	for _, peer := range netInfo.Result.Peers {
		versionCounts[peer.NodeInfo.Version]++
	}

	if versionCounts["1.0.0"] != 2 || versionCounts["1.1.0"] != 1 {
		t.Errorf("Expected 2 peers with version 1.0.0 and 1 peer with version 1.1.0, got %v", versionCounts)
	}
}
