package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ssccio/rebalance/pkg/rebalancer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	// Import all auth plugins (for cloud providers)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	config rebalancer.Config

	// Version information (set by goreleaser)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Detect if running as kubectl plugin
	use := "kubectl-rebalance"
	if strings.HasPrefix(filepath.Base(os.Args[0]), "kubectl-") {
		use = "kubectl rebalance"
	}

	rootCmd := &cobra.Command{
		Use:     use,
		Short:   "Rebalance pods from overloaded nodes",
		Version: version,
		Long: `A kubectl plugin to evict pods from heavily loaded nodes,
allowing them to be rescheduled on less loaded nodes.

REQUIRED FLAGS:
  --selector, -l    Label selector to identify pods to evict

EXAMPLES:
  # Evict 10 pods with label 'name=php' from the 3 most loaded nodes
  kubectl rebalance --count 10 --selector name=php --from-nodes 3

  # Dry-run to see what would be evicted
  kubectl rebalance --count 10 --selector name=php --dry-run

  # Use pod-count for clusters without metrics-server
  kubectl rebalance --count 10 --selector name=php --metric pod-count

  # Target nodes over 90% utilization, with minimum 5 pods
  kubectl rebalance --count 20 --selector app=web --threshold 90 --min-pods 5`,
		RunE: run,
	}

	// If no args provided, show help
	if len(os.Args) == 1 {
		rootCmd.SetArgs([]string{"--help"})
	}
	flags := rootCmd.Flags()
	flags.IntVarP(&config.Count, "count", "c", 10, "Number of pods to evict")
	flags.IntVar(&config.MinPods, "min-pods", 0, "Minimum pods required to proceed")
	flags.IntVar(&config.FromNodes, "from-nodes", 3, "Number of most loaded nodes to target")
	flags.Float64Var(&config.Threshold, "threshold", 80.0, "CPU/Memory threshold percentage")
	flags.StringVarP(&config.Selector, "selector", "l", "", "Label selector for pods (REQUIRED)")
	flags.StringVarP(&config.Namespace, "namespace", "n", "", "Namespace (empty for all)")
	flags.BoolVarP(&config.AllNamespaces, "all-namespaces", "A", false, "Target all namespaces (same as empty namespace)")
	flags.DurationVar(&config.Interval, "interval", 10*time.Second, "Interval between evictions")
	flags.BoolVar(&config.DryRun, "dry-run", false, "Show what would be evicted")
	flags.StringVar(&config.MetricType, "metric", "both", "Metric type: cpu, memory, both, or pod-count")
	flags.StringVarP(&config.OutputFormat, "output", "o", "table", "Output format: table, json, yaml")

	// Mark required flags
	rootCmd.MarkFlagRequired("selector")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Handle --all-namespaces flag
	if config.AllNamespaces {
		config.Namespace = ""
	}

	// Load kubeconfig
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Create Kubernetes client
	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Create metrics client
	metricsClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to create metrics client: %w", err)
	}
	// Create rebalancer
	rb := rebalancer.New(k8sClient, metricsClient, config)

	// Execute rebalancing
	return rb.Execute()
}
