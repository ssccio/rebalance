package rebalancer

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	config := Config{
		Count:        10,
		MinPods:      5,
		FromNodes:    3,
		Threshold:    80.0,
		Selector:     "app=test",
		Namespace:    "default",
		Interval:     10 * time.Second,
		DryRun:       true,
		MetricType:   "both",
		OutputFormat: "table",
	}

	if config.Count != 10 {
		t.Errorf("Expected Count to be 10, got %d", config.Count)
	}

	if config.Selector != "app=test" {
		t.Errorf("Expected Selector to be 'app=test', got %s", config.Selector)
	}
}