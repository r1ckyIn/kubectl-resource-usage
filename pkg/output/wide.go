package output

import (
	"fmt"
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
	"k8s.io/apimachinery/pkg/api/resource"
)

// WideFormatter formats output as a wide table with requests/limits raw values
type WideFormatter struct{}

// Format writes pod usages as a wide table
func (f *WideFormatter) Format(w io.Writer, podUsages []calculator.PodUsage) error {
	// Print header
	fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s\n",
		wideColNamespace, "NAMESPACE",
		wideColPod, "POD",
		wideColUsage, "CPU_USAGE",
		wideColReqLim, "CPU_REQ",
		wideColReqLim, "CPU_LIM",
		wideColPercent, "CPU_R%",
		wideColPercent, "CPU_L%",
		wideColUsage, "MEM_USAGE",
		wideColReqLim, "MEM_REQ",
		wideColReqLim, "MEM_LIM",
		wideColPercent, "MEM_R%",
		wideColPercent, "MEM_L%",
		wideColNode, "NODE")

	// Print rows
	for _, pu := range podUsages {
		fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s\n",
			wideColNamespace, truncate(pu.Namespace, wideColNamespace),
			wideColPod, truncate(pu.Name, wideColPod),
			wideColUsage, formatCPU(pu.CPU.Usage.MilliValue()),
			wideColReqLim, formatQuantityOrNA(pu.CPU.Requests),
			wideColReqLim, formatQuantityOrNA(pu.CPU.Limits),
			wideColPercent, formatPercent(pu.CPU.RequestPercent),
			wideColPercent, formatPercent(pu.CPU.LimitPercent),
			wideColUsage, formatMemory(pu.Memory.Usage.Value()),
			wideColReqLim, formatMemoryQuantityOrNA(pu.Memory.Requests),
			wideColReqLim, formatMemoryQuantityOrNA(pu.Memory.Limits),
			wideColPercent, formatPercent(pu.Memory.RequestPercent),
			wideColPercent, formatPercent(pu.Memory.LimitPercent),
			wideColNode, truncate(pu.Node, wideColNode),
		)
	}

	return nil
}

// formatQuantityOrNA formats a CPU quantity or returns "N/A"
func formatQuantityOrNA(q *resource.Quantity) string {
	if q == nil {
		return "N/A"
	}
	return formatCPU(q.MilliValue())
}

// formatMemoryQuantityOrNA formats a memory quantity or returns "N/A"
func formatMemoryQuantityOrNA(q *resource.Quantity) string {
	if q == nil {
		return "N/A"
	}
	return formatMemory(q.Value())
}
