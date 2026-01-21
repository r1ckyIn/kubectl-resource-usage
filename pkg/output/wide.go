package output

import (
	"fmt"
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
	"k8s.io/apimachinery/pkg/api/resource"
)

// WideFormatter formats output as a wide table with requests/limits raw values
type WideFormatter struct {
	colorizer     *Colorizer
	unitFormatter *UnitFormatter
}

// Format writes pod usages as a wide table
func (f *WideFormatter) Format(w io.Writer, podUsages []calculator.PodUsage) error {
	// Print header
	if _, err := fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s\n",
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
		wideColNode, "NODE"); err != nil {
		return err
	}

	// Print rows
	for _, pu := range podUsages {
		if _, err := fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %-*s %s %s %-*s %-*s %-*s %s %s %-*s\n",
			wideColNamespace, truncate(pu.Namespace, wideColNamespace),
			wideColPod, truncate(pu.Name, wideColPod),
			wideColUsage, f.unitFormatter.FormatCPU(pu.CPU.Usage.MilliValue()),
			wideColReqLim, f.formatCPUQuantityOrNA(pu.CPU.Requests),
			wideColReqLim, f.formatCPUQuantityOrNA(pu.CPU.Limits),
			f.colorizer.FormatPercent(pu.CPU.RequestPercent, wideColPercent),
			f.colorizer.FormatPercent(pu.CPU.LimitPercent, wideColPercent),
			wideColUsage, f.unitFormatter.FormatMemory(pu.Memory.Usage.Value()),
			wideColReqLim, f.formatMemoryQuantityOrNA(pu.Memory.Requests),
			wideColReqLim, f.formatMemoryQuantityOrNA(pu.Memory.Limits),
			f.colorizer.FormatPercent(pu.Memory.RequestPercent, wideColPercent),
			f.colorizer.FormatPercent(pu.Memory.LimitPercent, wideColPercent),
			wideColNode, truncate(pu.Node, wideColNode),
		); err != nil {
			return err
		}
	}

	return nil
}

// formatCPUQuantityOrNA formats a CPU quantity or returns "N/A"
func (f *WideFormatter) formatCPUQuantityOrNA(q *resource.Quantity) string {
	if q == nil {
		return "N/A"
	}
	return f.unitFormatter.FormatCPU(q.MilliValue())
}

// formatMemoryQuantityOrNA formats a memory quantity or returns "N/A"
func (f *WideFormatter) formatMemoryQuantityOrNA(q *resource.Quantity) string {
	if q == nil {
		return "N/A"
	}
	return f.unitFormatter.FormatMemory(q.Value())
}
