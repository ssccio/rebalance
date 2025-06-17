package rebalancer

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// Rebalancer orchestrates the pod rebalancing process
type Rebalancer struct {
	k8sClient     kubernetes.Interface
	metricsClient versioned.Interface
	config        Config
	analyzer      *NodeAnalyzer
	selector      *PodSelector
	validator     *Validator
	evictor       *Evictor
}

// New creates a new Rebalancer instance
func New(k8sClient kubernetes.Interface, metricsClient versioned.Interface, config Config) *Rebalancer {
	return &Rebalancer{
		k8sClient:     k8sClient,
		metricsClient: metricsClient, config: config,
		analyzer:  NewNodeAnalyzer(k8sClient, metricsClient),
		selector:  NewPodSelector(k8sClient),
		validator: NewValidator(k8sClient),
		evictor:   NewEvictor(k8sClient),
	}
}

// Execute runs the rebalancing process
func (r *Rebalancer) Execute() error {
	ctx := context.Background()

	// Step 1: Analyze nodes and find overloaded ones
	fmt.Println("Analyzing node resource usage...")
	overloadedNodes, err := r.analyzer.GetOverloadedNodes(ctx, r.config.FromNodes, r.config.Threshold, r.config.MetricType)
	if err != nil {
		return fmt.Errorf("failed to analyze nodes: %w", err)
	}

	if len(overloadedNodes) == 0 {
		fmt.Printf("No nodes found with utilization above %.0f%%\n", r.config.Threshold)
		return nil
	}

	fmt.Printf("Found %d overloaded nodes\n", len(overloadedNodes))
	for _, node := range overloadedNodes {
		if r.config.MetricType == "pod-count" {
			fmt.Printf("  - %s (Pods: %d)\n", node.Name, int(node.Score))
		} else {
			fmt.Printf("  - %s (CPU: %.1f%%, Memory: %.1f%%)\n",
				node.Name, node.CPUPercent, node.MemoryPercent)
		}
	}

	// Step 2: Select pods from overloaded nodes
	fmt.Printf("\nSelecting pods with selector '%s'...\n", r.config.Selector)
	pods, err := r.selector.SelectPods(ctx, overloadedNodes, r.config.Selector, r.config.Namespace, r.config.Count)
	if err != nil {
		return fmt.Errorf("failed to select pods: %w", err)
	}

	// Step 3: Validate the selection
	if err := r.validator.Validate(pods, r.config); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("Found %d pods to evict\n", len(pods))

	// Step 4: Show what will be done (dry-run or execute)
	if r.config.DryRun {
		fmt.Println("\nDRY-RUN MODE - Would evict the following pods:")
		for i, pod := range pods {
			fmt.Printf("  %d. %s/%s (node: %s)\n", i+1, pod.Namespace, pod.Name, pod.Spec.NodeName)
		}
		return nil
	}

	// Step 5: Execute evictions
	fmt.Println("\nStarting pod evictions...")
	if err := r.evictor.EvictPods(ctx, pods, r.config.Interval); err != nil {
		return fmt.Errorf("eviction failed: %w", err)
	}

	fmt.Println("\nRebalancing complete!")
	return nil
}
