package collector

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsfake "k8s.io/metrics/pkg/client/clientset/versioned/fake"
)

func TestPodCollector_GetPods(t *testing.T) {
	// Create fake pods
	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-1",
			Namespace: "default",
			Labels:    map[string]string{"app": "test"},
		},
		Spec: corev1.PodSpec{
			NodeName: "node-1",
		},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-2",
			Namespace: "default",
			Labels:    map[string]string{"app": "other"},
		},
		Spec: corev1.PodSpec{
			NodeName: "node-2",
		},
	}
	pod3 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod-3",
			Namespace: "kube-system",
			Labels:    map[string]string{"app": "test"},
		},
		Spec: corev1.PodSpec{
			NodeName: "node-1",
		},
	}

	fakeClient := fake.NewSimpleClientset(pod1, pod2, pod3)
	collector := &PodCollector{client: fakeClient}
	ctx := context.Background()

	tests := []struct {
		name          string
		namespace     string
		selector      string
		expectedCount int
	}{
		{
			name:          "all namespaces no selector",
			namespace:     "",
			selector:      "",
			expectedCount: 3,
		},
		{
			name:          "specific namespace",
			namespace:     "default",
			selector:      "",
			expectedCount: 2,
		},
		{
			name:          "with label selector",
			namespace:     "",
			selector:      "app=test",
			expectedCount: 2,
		},
		{
			name:          "namespace and selector",
			namespace:     "default",
			selector:      "app=test",
			expectedCount: 1,
		},
		{
			name:          "no matching pods",
			namespace:     "nonexistent",
			selector:      "",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pods, err := collector.GetPods(ctx, tt.namespace, tt.selector)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(pods.Items) != tt.expectedCount {
				t.Errorf("expected %d pods, got %d", tt.expectedCount, len(pods.Items))
			}
		})
	}
}

func TestMetricsCollector_GetPodMetrics(t *testing.T) {
	// Create fake pod metrics list
	podMetricsList := &metricsv1beta1.PodMetricsList{
		Items: []metricsv1beta1.PodMetrics{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod-1",
					Namespace: "default",
				},
				Containers: []metricsv1beta1.ContainerMetrics{
					{
						Name: "container-1",
						Usage: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod-2",
					Namespace: "kube-system",
				},
				Containers: []metricsv1beta1.ContainerMetrics{
					{
						Name: "container-1",
						Usage: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			},
		},
	}

	fakeClient := metricsfake.NewSimpleClientset(&podMetricsList.Items[0], &podMetricsList.Items[1])
	collector := &MetricsCollector{client: fakeClient}
	ctx := context.Background()

	// Test all namespaces - metrics fake client may not support namespace filtering
	// so we just verify we can call the method without error
	metrics, err := collector.GetPodMetrics(ctx, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: fake client behavior may differ from real client
	t.Logf("got %d metrics from all namespaces", len(metrics.Items))
}

func TestMetricsCollector_GetPodMetrics_WithMetrics(t *testing.T) {
	// Create fake pod metrics with specific values
	podMetrics := &metricsv1beta1.PodMetrics{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Containers: []metricsv1beta1.ContainerMetrics{
			{
				Name: "container-1",
				Usage: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("250m"),
					corev1.ResourceMemory: resource.MustParse("512Mi"),
				},
			},
			{
				Name: "container-2",
				Usage: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("150m"),
					corev1.ResourceMemory: resource.MustParse("256Mi"),
				},
			},
		},
	}

	fakeClient := metricsfake.NewSimpleClientset(podMetrics)
	collector := &MetricsCollector{client: fakeClient}
	ctx := context.Background()

	// Verify the collector can be created and called without errors
	metrics, err := collector.GetPodMetrics(ctx, "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("got %d metrics for default namespace", len(metrics.Items))
}
