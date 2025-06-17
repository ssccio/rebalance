package rebalancer

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	config := Config{
		Count:         10,
		MinPods:       5,
		FromNodes:     3,
		Threshold:     80.0,
		Selector:      "app=test",
		Namespace:     "default",
		AllNamespaces: false,
		Interval:      10 * time.Second,
		DryRun:        true,
		MetricType:    "both",
		OutputFormat:  "table",
	}

	if config.Count != 10 {
		t.Errorf("expected Count to be 10, got %d", config.Count)
	}

	if config.Selector != "app=test" {
		t.Errorf("expected Selector to be 'app=test', got %s", config.Selector)
	}
}

func TestIsNodeReady(t *testing.T) {
	// This is a placeholder for more comprehensive tests
	// In a real implementation, you would test with mock nodes
	t.Log("Node readiness check test placeholder")
}
