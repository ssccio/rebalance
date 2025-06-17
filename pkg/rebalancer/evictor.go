package rebalancer

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Evictor handles pod evictions
type Evictor struct {
	k8sClient kubernetes.Interface
}

// NewEvictor creates a new Evictor
func NewEvictor(k8sClient kubernetes.Interface) *Evictor {
	return &Evictor{k8sClient: k8sClient}
}

// EvictPods evicts the given pods with the specified interval between evictions
func (e *Evictor) EvictPods(ctx context.Context, pods []corev1.Pod, interval time.Duration) error {
	for i, pod := range pods {
		fmt.Printf("Evicting pod %d/%d: %s/%s\n", i+1, len(pods), pod.Namespace, pod.Name)

		eviction := &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
		}

		err := e.k8sClient.PolicyV1().Evictions(pod.Namespace).Evict(ctx, eviction)
		if err != nil {
			// Log error but continue with other pods
			fmt.Printf("  WARNING: Failed to evict pod %s/%s: %v\n", pod.Namespace, pod.Name, err)
			continue
		}

		fmt.Printf("  Successfully evicted pod %s/%s\n", pod.Namespace, pod.Name)

		// Wait before next eviction (except for the last one)
		if i < len(pods)-1 {
			fmt.Printf("  Waiting %v before next eviction...\n", interval)
			time.Sleep(interval)
		}
	}

	return nil
}
