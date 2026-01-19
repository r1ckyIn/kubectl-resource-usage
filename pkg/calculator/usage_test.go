package calculator

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func TestCalculatePercent(t *testing.T) {
	tests := []struct {
		name     string
		usage    *resource.Quantity
		base     *resource.Quantity
		expected *int
	}{
		{
			name:     "normal calculation",
			usage:    resourcePtr(resource.MustParse("500m")),
			base:     resourcePtr(resource.MustParse("1000m")),
			expected: intPtr(50),
		},
		{
			name:     "over 100%",
			usage:    resourcePtr(resource.MustParse("2000m")),
			base:     resourcePtr(resource.MustParse("1000m")),
			expected: intPtr(200),
		},
		{
			name:     "nil base returns nil",
			usage:    resourcePtr(resource.MustParse("500m")),
			base:     nil,
			expected: nil,
		},
		{
			name:     "zero base returns nil",
			usage:    resourcePtr(resource.MustParse("500m")),
			base:     resourcePtr(resource.MustParse("0")),
			expected: nil,
		},
		{
			name:     "memory calculation",
			usage:    resourcePtr(resource.MustParse("256Mi")),
			base:     resourcePtr(resource.MustParse("512Mi")),
			expected: intPtr(50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePercent(tt.usage, tt.base)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %d", *result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %d, got nil", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("expected %d, got %d", *tt.expected, *result)
				}
			}
		})
	}
}

func TestCalculatePodUsage(t *testing.T) {
	podMetric := metricsv1beta1.PodMetrics{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Containers: []metricsv1beta1.ContainerMetrics{
			{
				Name: "container1",
				Usage: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("128Mi"),
				},
			},
		},
	}

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			NodeName: "node-1",
			Containers: []corev1.Container{
				{
					Name: "container1",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("512Mi"),
						},
					},
				},
			},
		},
	}

	result := CalculatePodUsage(podMetric, pod)

	if result.Namespace != "default" {
		t.Errorf("expected namespace 'default', got '%s'", result.Namespace)
	}
	if result.Name != "test-pod" {
		t.Errorf("expected name 'test-pod', got '%s'", result.Name)
	}
	if result.Node != "node-1" {
		t.Errorf("expected node 'node-1', got '%s'", result.Node)
	}

	// CPU: 100m / 200m = 50% request, 100m / 500m = 20% limit
	if result.CPU.RequestPercent == nil || *result.CPU.RequestPercent != 50 {
		t.Errorf("expected CPU request percent 50, got %v", result.CPU.RequestPercent)
	}
	if result.CPU.LimitPercent == nil || *result.CPU.LimitPercent != 20 {
		t.Errorf("expected CPU limit percent 20, got %v", result.CPU.LimitPercent)
	}

	// Memory: 128Mi / 256Mi = 50% request, 128Mi / 512Mi = 25% limit
	if result.Memory.RequestPercent == nil || *result.Memory.RequestPercent != 50 {
		t.Errorf("expected Memory request percent 50, got %v", result.Memory.RequestPercent)
	}
	if result.Memory.LimitPercent == nil || *result.Memory.LimitPercent != 25 {
		t.Errorf("expected Memory limit percent 25, got %v", result.Memory.LimitPercent)
	}
}

func TestCalculatePodUsageNoLimits(t *testing.T) {
	podMetric := metricsv1beta1.PodMetrics{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Containers: []metricsv1beta1.ContainerMetrics{
			{
				Name: "container1",
				Usage: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("128Mi"),
				},
			},
		},
	}

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			NodeName: "node-1",
			Containers: []corev1.Container{
				{
					Name:      "container1",
					Resources: corev1.ResourceRequirements{},
				},
			},
		},
	}

	result := CalculatePodUsage(podMetric, pod)

	if result.CPU.RequestPercent != nil {
		t.Errorf("expected nil CPU request percent, got %d", *result.CPU.RequestPercent)
	}
	if result.CPU.LimitPercent != nil {
		t.Errorf("expected nil CPU limit percent, got %d", *result.CPU.LimitPercent)
	}
	if result.Memory.RequestPercent != nil {
		t.Errorf("expected nil Memory request percent, got %d", *result.Memory.RequestPercent)
	}
	if result.Memory.LimitPercent != nil {
		t.Errorf("expected nil Memory limit percent, got %d", *result.Memory.LimitPercent)
	}
}

func TestSortPodUsages(t *testing.T) {
	pods := []PodUsage{
		{Name: "pod1", CPU: ResourceUsage{LimitPercent: intPtr(50)}, Memory: ResourceUsage{LimitPercent: intPtr(30)}},
		{Name: "pod2", CPU: ResourceUsage{LimitPercent: intPtr(80)}, Memory: ResourceUsage{LimitPercent: intPtr(60)}},
		{Name: "pod3", CPU: ResourceUsage{LimitPercent: nil}, Memory: ResourceUsage{LimitPercent: nil}},
		{Name: "pod4", CPU: ResourceUsage{LimitPercent: intPtr(20)}, Memory: ResourceUsage{LimitPercent: intPtr(90)}},
	}

	// Sort by CPU descending
	SortPodUsages(pods, "cpu", false)
	if pods[0].Name != "pod2" || pods[1].Name != "pod1" || pods[2].Name != "pod4" || pods[3].Name != "pod3" {
		t.Errorf("CPU descending sort failed: %v", []string{pods[0].Name, pods[1].Name, pods[2].Name, pods[3].Name})
	}

	// Sort by Memory ascending
	SortPodUsages(pods, "memory", true)
	if pods[0].Name != "pod1" || pods[1].Name != "pod2" || pods[2].Name != "pod4" || pods[3].Name != "pod3" {
		t.Errorf("Memory ascending sort failed: %v", []string{pods[0].Name, pods[1].Name, pods[2].Name, pods[3].Name})
	}
}

func intPtr(i int) *int {
	return &i
}

func resourcePtr(r resource.Quantity) *resource.Quantity {
	return &r
}
