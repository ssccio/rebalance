package rebalancer

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// PodSelector selects pods from nodes based on criteria
type PodSelector struct {
	k8sClient kubernetes.Interface
}

// NewPodSelector creates a new PodSelector
func NewPodSelector(k8sClient kubernetes.Interface) *PodSelector {
	return &PodSelector{k8sClient: k8sClient}
}

// SelectPods selects pods from the given nodes matching the selector
func (ps *PodSelector) SelectPods(ctx context.Context, nodes []NodeInfo, selector string, namespace string, count int) ([]corev1.Pod, error) {
	// Parse the label selector
	labelSelector, err := labels.Parse(selector)
	if err != nil {
		return nil, fmt.Errorf("invalid label selector: %w", err)
	}

	selectedPods := make([]corev1.Pod, 0, count)

	// Iterate through nodes and collect matching pods
	for _, node := range nodes {
		if len(selectedPods) >= count {
			break
		}

		// List pods on this node
		listOptions := metav1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name),
		}

		if namespace != "" {
			pods, err := ps.k8sClient.CoreV1().Pods(namespace).List(ctx, listOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
			}
			selectedPods = append(selectedPods, ps.filterPods(pods.Items, labelSelector, count-len(selectedPods))...)
		} else {
			// List pods in all namespaces
			pods, err := ps.k8sClient.CoreV1().Pods("").List(ctx, listOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list pods: %w", err)
			}
			selectedPods = append(selectedPods, ps.filterPods(pods.Items, labelSelector, count-len(selectedPods))...)
		}
	}

	return selectedPods, nil
}

// filterPods filters pods based on label selector and returns up to maxCount pods
func (ps *PodSelector) filterPods(pods []corev1.Pod, selector labels.Selector, maxCount int) []corev1.Pod {
	filtered := make([]corev1.Pod, 0, maxCount)

	for _, pod := range pods {
		if len(filtered) >= maxCount {
			break
		}

		// Skip pods that are not running or are being deleted
		if pod.Status.Phase != corev1.PodRunning || pod.DeletionTimestamp != nil {
			continue
		}

		// Check if pod matches selector
		if selector.Matches(labels.Set(pod.Labels)) {
			filtered = append(filtered, pod)
		}
	}

	return filtered
}
