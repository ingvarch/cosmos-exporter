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
		if r.URL.Path == "/block" {
			w.Write([]byte(`
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
		}
	}))
	defer mockServer.Close()

	// Create a new exporter with the mock server URL
	exporter := NewCosmosExporter(mockServer.URL)

	// Call updateMetrics
	exporter.updateMetrics()

	// Check if the highestBlock metric has the correct value
	expected := 23921295
	if actual := testutil.ToFloat64(exporter.highestBlock); int(actual) != expected {
		t.Errorf("Unexpected value for highestBlock: got %v, want %v", int(actual), expected)
	}
}

func TestNewCosmosExporter(t *testing.T) {
	cosmosNodeAddress := "http://localhost:26657"
	exporter := NewCosmosExporter(cosmosNodeAddress)
	if exporter == nil {
		t.Error("NewCosmosExporter returned nil")
	}
	if exporter.cosmosNodeAddress != cosmosNodeAddress {
		t.Errorf("Expected CosmosNodeAddress to be %s, but got %s", cosmosNodeAddress, exporter.cosmosNodeAddress)
	}
}
