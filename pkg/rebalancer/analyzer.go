package rebalancer

import (
	"context"
	"fmt"
	"sort"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// NodeInfo holds node information with resource usage
type NodeInfo struct {
	Name          string
	CPUPercent    float64
	MemoryPercent float64
	Score         float64 // Combined score for sorting
}

// NodeAnalyzer analyzes nodes and their resource usage
type NodeAnalyzer struct {
	k8sClient     kubernetes.Interface
	metricsClient versioned.Interface
}

// NewNodeAnalyzer creates a new NodeAnalyzer
func NewNodeAnalyzer(k8sClient kubernetes.Interface, metricsClient versioned.Interface) *NodeAnalyzer {
	return &NodeAnalyzer{
		k8sClient:     k8sClient,
		metricsClient: metricsClient,
	}
}

// GetOverloadedNodes returns the top N overloaded nodes based on the metric type
func (na *NodeAnalyzer) GetOverloadedNodes(ctx context.Context, count int, threshold float64, metricType string) ([]NodeInfo, error) {
	// Get all nodes
	nodes, err := na.k8sClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	// If using pod-count metric, use a different approach
	if metricType == "pod-count" {
		return na.getOverloadedNodesByPodCount(ctx, nodes.Items, count)
	}

	// Get node metrics
	nodeMetrics, err := na.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		// If metrics-server is not available, fall back to pod count
		fmt.Printf("Warning: metrics-server not available (%v), falling back to pod-count mode\n", err)
		return na.getOverloadedNodesByPodCount(ctx, nodes.Items, count)
	}

	// Create a map for quick metric lookup
	metricsMap := make(map[string]*metricsv1beta1.NodeMetrics)
	for i := range nodeMetrics.Items {
		metricsMap[nodeMetrics.Items[i].Name] = &nodeMetrics.Items[i]
	}

	// Calculate utilization for each node
	nodeInfos := make([]NodeInfo, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		// Skip nodes that are not ready
		if !isNodeReady(&node) {
			continue
		}

		metrics, ok := metricsMap[node.Name]
		if !ok {
			continue // Skip nodes without metrics
		}

		info := NodeInfo{Name: node.Name}

		// Calculate CPU percentage
		cpuCapacity := node.Status.Allocatable[corev1.ResourceCPU]
		cpuUsage := metrics.Usage[corev1.ResourceCPU]
		info.CPUPercent = float64(cpuUsage.MilliValue()) / float64(cpuCapacity.MilliValue()) * 100

		// Calculate memory percentage
		memCapacity := node.Status.Allocatable[corev1.ResourceMemory]
		memUsage := metrics.Usage[corev1.ResourceMemory]
		info.MemoryPercent = float64(memUsage.Value()) / float64(memCapacity.Value()) * 100

		// Calculate combined score based on metric type
		switch metricType {
		case "cpu":
			info.Score = info.CPUPercent
		case "memory":
			info.Score = info.MemoryPercent
		default: // "both"
			info.Score = (info.CPUPercent + info.MemoryPercent) / 2
		}

		// Only include nodes above threshold
		if info.Score >= threshold {
			nodeInfos = append(nodeInfos, info)
		}
	}

	// Sort by score (highest first)
	sort.Slice(nodeInfos, func(i, j int) bool {
		return nodeInfos[i].Score > nodeInfos[j].Score
	})

	// Return top N nodes
	if len(nodeInfos) > count {
		nodeInfos = nodeInfos[:count]
	}

	return nodeInfos, nil
}

// isNodeReady checks if a node is in Ready condition
func isNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

// getOverloadedNodesByPodCount returns nodes sorted by pod count
func (na *NodeAnalyzer) getOverloadedNodesByPodCount(ctx context.Context, nodes []corev1.Node, count int) ([]NodeInfo, error) {
	nodeInfos := make([]NodeInfo, 0, len(nodes))

	for _, node := range nodes {
		// Skip nodes that are not ready
		if !isNodeReady(&node) {
			continue
		}

		// Count pods on this node
		pods, err := na.k8sClient.CoreV1().Pods("").List(ctx, metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name),
		})
		if err != nil {
			fmt.Printf("Warning: failed to count pods on node %s: %v\n", node.Name, err)
			continue
		}

		// Count only running pods
		runningPods := 0
		for _, pod := range pods.Items {
			if pod.Status.Phase == corev1.PodRunning && pod.DeletionTimestamp == nil {
				runningPods++
			}
		}

		info := NodeInfo{
			Name:          node.Name,
			CPUPercent:    0,                    // Not available without metrics
			MemoryPercent: 0,                    // Not available without metrics
			Score:         float64(runningPods), // Use pod count as score
		}

		nodeInfos = append(nodeInfos, info)
	}

	// Sort by pod count (highest first)
	sort.Slice(nodeInfos, func(i, j int) bool {
		return nodeInfos[i].Score > nodeInfos[j].Score
	})

	// Return top N nodes
	if len(nodeInfos) > count {
		nodeInfos = nodeInfos[:count]
	}

	// Add pod count to the display
	for i := range nodeInfos {
		nodeInfos[i].CPUPercent = nodeInfos[i].Score    // Show pod count as "CPU%"
		nodeInfos[i].MemoryPercent = nodeInfos[i].Score // Show pod count as "Memory%"
	}

	return nodeInfos, nil
}
