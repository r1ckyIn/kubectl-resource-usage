package collector

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// PodCollector fetches Pod specs from the Kubernetes API
type PodCollector struct {
	client *kubernetes.Clientset
}

// NewPodCollector creates a new PodCollector
func NewPodCollector(config *rest.Config) (*PodCollector, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &PodCollector{
		client: client,
	}, nil
}

// GetPods fetches pods for the specified namespace with optional label selector
// If namespace is empty, it fetches pods from all namespaces
func (c *PodCollector) GetPods(ctx context.Context, namespace, selector string) (*corev1.PodList, error) {
	opts := metav1.ListOptions{}
	if selector != "" {
		opts.LabelSelector = selector
	}

	pods, err := c.client.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	return pods, nil
}
