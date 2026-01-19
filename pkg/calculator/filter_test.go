package calculator

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
)

func quantityPtr(s string) *resource.Quantity {
	q := resource.MustParse(s)
	return &q
}

func TestFilterPodUsages_Above(t *testing.T) {
	pods := []PodUsage{
		{Name: "pod1", Memory: ResourceUsage{LimitPercent: intPtr(90)}},
		{Name: "pod2", Memory: ResourceUsage{LimitPercent: intPtr(50)}},
		{Name: "pod3", Memory: ResourceUsage{LimitPercent: intPtr(80)}},
		{Name: "pod4", Memory: ResourceUsage{LimitPercent: nil}},
	}

	opts := FilterOptions{Above: 80, Below: -1, Field: "memory"}
	result := FilterPodUsages(pods, opts)

	if len(result) != 2 {
		t.Errorf("expected 2 pods, got %d", len(result))
	}

	expectedNames := map[string]bool{"pod1": true, "pod3": true}
	for _, pod := range result {
		if !expectedNames[pod.Name] {
			t.Errorf("unexpected pod: %s", pod.Name)
		}
	}
}

func TestFilterPodUsages_Below(t *testing.T) {
	pods := []PodUsage{
		{Name: "pod1", Memory: ResourceUsage{LimitPercent: intPtr(90)}},
		{Name: "pod2", Memory: ResourceUsage{LimitPercent: intPtr(30)}},
		{Name: "pod3", Memory: ResourceUsage{LimitPercent: intPtr(50)}},
	}

	opts := FilterOptions{Above: -1, Below: 50, Field: "memory"}
	result := FilterPodUsages(pods, opts)

	if len(result) != 2 {
		t.Errorf("expected 2 pods, got %d", len(result))
	}

	expectedNames := map[string]bool{"pod2": true, "pod3": true}
	for _, pod := range result {
		if !expectedNames[pod.Name] {
			t.Errorf("unexpected pod: %s", pod.Name)
		}
	}
}

func TestFilterPodUsages_NoLimits(t *testing.T) {
	pods := []PodUsage{
		{
			Name: "pod1",
			CPU:  ResourceUsage{Limits: quantityPtr("100m")},
			Memory: ResourceUsage{Limits: quantityPtr("128Mi")},
		},
		{
			Name: "pod2",
			CPU:  ResourceUsage{Limits: nil},
			Memory: ResourceUsage{Limits: quantityPtr("128Mi")},
		},
		{
			Name: "pod3",
			CPU:  ResourceUsage{Limits: quantityPtr("100m")},
			Memory: ResourceUsage{Limits: nil},
		},
		{
			Name: "pod4",
			CPU:  ResourceUsage{Limits: nil},
			Memory: ResourceUsage{Limits: nil},
		},
	}

	opts := FilterOptions{Above: -1, Below: -1, NoLimits: true}
	result := FilterPodUsages(pods, opts)

	if len(result) != 3 {
		t.Errorf("expected 3 pods, got %d", len(result))
	}

	expectedNames := map[string]bool{"pod2": true, "pod3": true, "pod4": true}
	for _, pod := range result {
		if !expectedNames[pod.Name] {
			t.Errorf("unexpected pod: %s", pod.Name)
		}
	}
}

func TestFilterPodUsages_CPUField(t *testing.T) {
	pods := []PodUsage{
		{Name: "pod1", CPU: ResourceUsage{LimitPercent: intPtr(90)}, Memory: ResourceUsage{LimitPercent: intPtr(30)}},
		{Name: "pod2", CPU: ResourceUsage{LimitPercent: intPtr(50)}, Memory: ResourceUsage{LimitPercent: intPtr(90)}},
		{Name: "pod3", CPU: ResourceUsage{LimitPercent: intPtr(80)}, Memory: ResourceUsage{LimitPercent: intPtr(10)}},
	}

	opts := FilterOptions{Above: 80, Below: -1, Field: "cpu"}
	result := FilterPodUsages(pods, opts)

	if len(result) != 2 {
		t.Errorf("expected 2 pods, got %d", len(result))
	}

	expectedNames := map[string]bool{"pod1": true, "pod3": true}
	for _, pod := range result {
		if !expectedNames[pod.Name] {
			t.Errorf("unexpected pod: %s", pod.Name)
		}
	}
}

func TestFilterPodUsages_AboveAndBelow(t *testing.T) {
	pods := []PodUsage{
		{Name: "pod1", Memory: ResourceUsage{LimitPercent: intPtr(90)}},
		{Name: "pod2", Memory: ResourceUsage{LimitPercent: intPtr(30)}},
		{Name: "pod3", Memory: ResourceUsage{LimitPercent: intPtr(60)}},
		{Name: "pod4", Memory: ResourceUsage{LimitPercent: intPtr(70)}},
	}

	opts := FilterOptions{Above: 50, Below: 80, Field: "memory"}
	result := FilterPodUsages(pods, opts)

	if len(result) != 2 {
		t.Errorf("expected 2 pods, got %d", len(result))
	}

	expectedNames := map[string]bool{"pod3": true, "pod4": true}
	for _, pod := range result {
		if !expectedNames[pod.Name] {
			t.Errorf("unexpected pod: %s", pod.Name)
		}
	}
}

func TestFilterPodUsages_NoFilter(t *testing.T) {
	pods := []PodUsage{
		{Name: "pod1"},
		{Name: "pod2"},
	}

	opts := FilterOptions{Above: -1, Below: -1, NoLimits: false}
	result := FilterPodUsages(pods, opts)

	if len(result) != 2 {
		t.Errorf("expected 2 pods, got %d", len(result))
	}
}

func TestNewFilterOptions(t *testing.T) {
	opts := NewFilterOptions()

	if opts.Above != -1 {
		t.Errorf("expected Above to be -1, got %d", opts.Above)
	}
	if opts.Below != -1 {
		t.Errorf("expected Below to be -1, got %d", opts.Below)
	}
	if opts.Field != "memory" {
		t.Errorf("expected Field to be 'memory', got %s", opts.Field)
	}
	if opts.NoLimits {
		t.Errorf("expected NoLimits to be false")
	}
}
