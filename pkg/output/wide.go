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
	fmt.Fprintf(w, "%-12s %-30s %-9s %-9s %-9s %-8s %-8s %-9s %-9s %-9s %-8s %-8s %-12s\n",
		"NAMESPACE", "POD", "CPU_USAGE", "CPU_REQ", "CPU_LIM", "CPU_R%", "CPU_L%",
		"MEM_USAGE", "MEM_REQ", "MEM_LIM", "MEM_R%", "MEM_L%", "NODE")

	// Print rows
	for _, pu := range podUsages {
		fmt.Fprintf(w, "%-12s %-30s %-9s %-9s %-9s %-8s %-8s %-9s %-9s %-9s %-8s %-8s %-12s\n",
			truncate(pu.Namespace, 12),
			truncate(pu.Name, 30),
			formatCPU(pu.CPU.Usage.MilliValue()),
			formatQuantityOrNA(pu.CPU.Requests),
			formatQuantityOrNA(pu.CPU.Limits),
			formatPercent(pu.CPU.RequestPercent),
			formatPercent(pu.CPU.LimitPercent),
			formatMemory(pu.Memory.Usage.Value()),
			formatMemoryQuantityOrNA(pu.Memory.Requests),
			formatMemoryQuantityOrNA(pu.Memory.Limits),
			formatPercent(pu.Memory.RequestPercent),
			formatPercent(pu.Memory.LimitPercent),
			truncate(pu.Node, 12),
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
