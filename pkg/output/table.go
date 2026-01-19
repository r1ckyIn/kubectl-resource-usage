package output

import (
	"fmt"
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
)

// TableFormatter formats output as a table
type TableFormatter struct{}

// Format writes pod usages as a table
func (f *TableFormatter) Format(w io.Writer, podUsages []calculator.PodUsage) error {
	// Print header
	fmt.Fprintf(w, "%-14s %-40s %-11s %-10s %-10s %-11s %-10s %-10s %-15s\n",
		"NAMESPACE", "POD", "CPU_USAGE", "CPU_REQ%", "CPU_LIM%", "MEM_USAGE", "MEM_REQ%", "MEM_LIM%", "NODE")

	// Print rows
	for _, pu := range podUsages {
		fmt.Fprintf(w, "%-14s %-40s %-11s %-10s %-10s %-11s %-10s %-10s %-15s\n",
			truncate(pu.Namespace, 14),
			truncate(pu.Name, 40),
			formatCPU(pu.CPU.Usage.MilliValue()),
			formatPercent(pu.CPU.RequestPercent),
			formatPercent(pu.CPU.LimitPercent),
			formatMemory(pu.Memory.Usage.Value()),
			formatPercent(pu.Memory.RequestPercent),
			formatPercent(pu.Memory.LimitPercent),
			truncate(pu.Node, 15),
		)
	}

	return nil
}

// formatCPU formats millicores to human readable format
func formatCPU(milliCores int64) string {
	return fmt.Sprintf("%dm", milliCores)
}

// formatMemory formats bytes to human readable format
func formatMemory(bytes int64) string {
	const (
		Ki = 1024
		Mi = Ki * 1024
		Gi = Mi * 1024
	)

	switch {
	case bytes >= Gi:
		return fmt.Sprintf("%dGi", bytes/Gi)
	case bytes >= Mi:
		return fmt.Sprintf("%dMi", bytes/Mi)
	case bytes >= Ki:
		return fmt.Sprintf("%dKi", bytes/Ki)
	default:
		return fmt.Sprintf("%d", bytes)
	}
}

// formatPercent formats percentage, returns "N/A" if nil
func formatPercent(p *int) string {
	if p == nil {
		return "N/A"
	}
	return fmt.Sprintf("%d%%", *p)
}

// truncate truncates string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
