package rebalancer

import (
	"time"
)

// Config holds all configuration for the rebalancer
type Config struct {
	Count         int
	MinPods       int
	FromNodes     int
	Threshold     float64
	Selector      string
	Namespace     string
	AllNamespaces bool
	Interval      time.Duration
	DryRun        bool
	MetricType    string
	OutputFormat  string
}