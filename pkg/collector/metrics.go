package collector

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

// MetricsCollector fetches pod metrics from the Kubernetes Metrics API
type MetricsCollector struct {
	client metricsclient.Interface
}

// NewMetricsCollector creates a new MetricsCollector
func NewMetricsCollector(config *rest.Config) (*MetricsCollector, error) {
	client, err := metricsclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	return &MetricsCollector{
		client: client,
	}, nil
}

// GetPodMetrics fetches pod metrics for the specified namespace
// If namespace is empty, it fetches metrics for all namespaces
func (c *MetricsCollector) GetPodMetrics(ctx context.Context, namespace string) (*metricsv1beta1.PodMetricsList, error) {
	podMetrics, err := c.client.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("metrics API not available: please install metrics-server")
		}
		if errors.IsForbidden(err) {
			return nil, fmt.Errorf("insufficient permissions to access metrics API")
		}
		return nil, fmt.Errorf("failed to list pod metrics")
	}

	return podMetrics, nil
}
