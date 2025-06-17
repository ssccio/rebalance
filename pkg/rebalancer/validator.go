package rebalancer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// Validator validates the pod selection before eviction
type Validator struct {
	k8sClient kubernetes.Interface
}

// NewValidator creates a new Validator
func NewValidator(k8sClient kubernetes.Interface) *Validator {
	return &Validator{k8sClient: k8sClient}
}

// Validate checks if the pod selection meets all criteria
func (v *Validator) Validate(pods []corev1.Pod, config Config) error {
	// Check if we found enough pods
	if len(pods) < config.Count {
		return fmt.Errorf("requested %d pods but only found %d matching pods on the most loaded nodes",
			config.Count, len(pods))
	}
	// Check minimum pods threshold
	if config.MinPods > 0 && len(pods) < config.MinPods {
		return fmt.Errorf("minimum threshold not met: required %d pods matching selector, but only found %d",
			config.MinPods, len(pods))
	}

	// TODO: Add more validation
	// - Check PodDisruptionBudgets
	// - Verify pods can be rescheduled (no conflicting nodeSelectors)
	// - Check if pods are part of critical system components

	return nil
}
