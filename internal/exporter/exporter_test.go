package exporter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestUpdateMetrics(t *testing.T) {
	// Create a mock server to simulate Cosmos node responses
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/block":
			_, err := w.Write([]byte(`
                {
                    "result": {
                        "block": {
                            "header": {
                                "height": "23921295",
                                "time": "2023-04-15T12:00:00Z"
                            }
                        }
                    }
                }
            `))
			if err != nil {
				t.Errorf("Error writing mock block response: %v", err)
			}
		case "/net_info":
			_, err := w.Write([]byte(`
                {
                    "result": {
                        "n_peers": "10",
                        "peers": [
                            {"node_info": {"version": "1.0.0"}},
                            {"node_info": {"version": "1.0.0"}},
                            {"node_info": {"version": "1.1.0"}}
                        ]
                    }
                }
            `))
			if err != nil {
				t.Errorf("Error writing mock net_info response: %v", err)
			}
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Create a new exporter with the mock server URL
	exporter := NewCosmosExporter(mockServer.URL)

	// Call updateMetrics
	exporter.updateMetrics()

	// Check if the highestBlock metric has the correct value
	expectedHeight := 23921295
	if actual := testutil.ToFloat64(exporter.highestBlock); int(actual) != expectedHeight {
		t.Errorf("Unexpected value for highestBlock: got %v, want %v", int(actual), expectedHeight)
	}

	// Check if the connectedPeers metric has the correct value
	expectedPeers := 10
	if actual := testutil.ToFloat64(exporter.connectedPeers); int(actual) != expectedPeers {
		t.Errorf("Unexpected value for connectedPeers: got %v, want %v", int(actual), expectedPeers)
	}

	// Check if the peersByVersion metric has the correct values
	expectedVersions := map[string]int{"1.0.0": 2, "1.1.0": 1}
	for version, count := range expectedVersions {
		if actual := testutil.ToFloat64(exporter.peersByVersion.WithLabelValues(version)); int(actual) != count {
			t.Errorf("Unexpected value for peersByVersion[%s]: got %v, want %v", version, int(actual), count)
		}
	}
}

func TestNewCosmosExporter(t *testing.T) {
	cosmosNodeAddress := "http://localhost:26657"
	exporter := NewCosmosExporter(cosmosNodeAddress)
	if exporter == nil {
		t.Fatal("NewCosmosExporter returned nil")
	}
	if exporter.cosmosNodeAddress != cosmosNodeAddress {
		t.Errorf("Expected CosmosNodeAddress to be %s, but got %s", cosmosNodeAddress, exporter.cosmosNodeAddress)
	}
	// Add additional checks for other fields if necessary
	if exporter.highestBlock == nil {
		t.Error("highestBlock metric is nil")
	}
	if exporter.blockTimeDrift == nil {
		t.Error("blockTimeDrift metric is nil")
	}
	if exporter.connectedPeers == nil {
		t.Error("connectedPeers metric is nil")
	}
	if exporter.peersByVersion == nil {
		t.Error("peersByVersion metric is nil")
	}
}
