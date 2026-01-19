package calculator

import (
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// ResourceUsage represents CPU or Memory usage with requests/limits
type ResourceUsage struct {
	Usage          resource.Quantity
	Requests       *resource.Quantity
	Limits         *resource.Quantity
	RequestPercent *int
	LimitPercent   *int
}

// PodUsage represents resource usage for a single pod
type PodUsage struct {
	Namespace string
	Name      string
	Node      string
	CPU       ResourceUsage
	Memory    ResourceUsage
}

// CalculatePercent calculates usage percentage relative to base
// Returns nil if base is nil or zero
func CalculatePercent(usage, base *resource.Quantity) *int {
	if base == nil || base.IsZero() {
		return nil
	}
	percent := int(usage.MilliValue() * 100 / base.MilliValue())
	return &percent
}

// CalculatePodUsage calculates resource usage for a pod
func CalculatePodUsage(podMetric metricsv1beta1.PodMetrics, pod corev1.Pod) PodUsage {
	// Sum up container metrics
	var totalCPU, totalMem resource.Quantity
	for _, container := range podMetric.Containers {
		totalCPU.Add(*container.Usage.Cpu())
		totalMem.Add(*container.Usage.Memory())
	}

	// Sum up container requests/limits
	var cpuReq, cpuLim, memReq, memLim resource.Quantity
	var hasCPUReq, hasCPULim, hasMemReq, hasMemLim bool

	for _, container := range pod.Spec.Containers {
		if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
			cpuReq.Add(req)
			hasCPUReq = true
		}
		if lim, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
			cpuLim.Add(lim)
			hasCPULim = true
		}
		if req, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
			memReq.Add(req)
			hasMemReq = true
		}
		if lim, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
			memLim.Add(lim)
			hasMemLim = true
		}
	}

	// Build ResourceUsage for CPU
	cpuUsage := ResourceUsage{
		Usage: totalCPU,
	}
	if hasCPUReq {
		cpuUsage.Requests = &cpuReq
		cpuUsage.RequestPercent = CalculatePercent(&totalCPU, &cpuReq)
	}
	if hasCPULim {
		cpuUsage.Limits = &cpuLim
		cpuUsage.LimitPercent = CalculatePercent(&totalCPU, &cpuLim)
	}

	// Build ResourceUsage for Memory
	memUsage := ResourceUsage{
		Usage: totalMem,
	}
	if hasMemReq {
		memUsage.Requests = &memReq
		memUsage.RequestPercent = CalculatePercent(&totalMem, &memReq)
	}
	if hasMemLim {
		memUsage.Limits = &memLim
		memUsage.LimitPercent = CalculatePercent(&totalMem, &memLim)
	}

	return PodUsage{
		Namespace: podMetric.Namespace,
		Name:      podMetric.Name,
		Node:      pod.Spec.NodeName,
		CPU:       cpuUsage,
		Memory:    memUsage,
	}
}

// SortPodUsages sorts pod usages by the specified field
// field can be "cpu" or "memory"
// N/A values are sorted to the end
func SortPodUsages(pods []PodUsage, field string, ascending bool) {
	sort.Slice(pods, func(i, j int) bool {
		var vi, vj *int

		if field == "cpu" {
			vi, vj = pods[i].CPU.LimitPercent, pods[j].CPU.LimitPercent
		} else {
			vi, vj = pods[i].Memory.LimitPercent, pods[j].Memory.LimitPercent
		}

		// N/A values go to the end
		if vi == nil && vj == nil {
			return false
		}
		if vi == nil {
			return false // i goes after j
		}
		if vj == nil {
			return true // j goes after i
		}

		if ascending {
			return *vi < *vj
		}
		return *vi > *vj
	})
}
