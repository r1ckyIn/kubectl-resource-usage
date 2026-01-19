package output

import (
	"fmt"
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
)

// TableFormatter formats output as a table
type TableFormatter struct {
	colorizer     *Colorizer
	unitFormatter *UnitFormatter
}

// Format writes pod usages as a table
func (f *TableFormatter) Format(w io.Writer, podUsages []calculator.PodUsage) error {
	// Print header
	fmt.Fprintf(w, "%-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s %-*s\n",
		tableColNamespace, "NAMESPACE",
		tableColPod, "POD",
		tableColCPUUsage, "CPU_USAGE",
		tableColPercent, "CPU_REQ%",
		tableColPercent, "CPU_LIM%",
		tableColMemUsage, "MEM_USAGE",
		tableColPercent, "MEM_REQ%",
		tableColPercent, "MEM_LIM%",
		tableColNode, "NODE")

	// Print rows
	for _, pu := range podUsages {
		fmt.Fprintf(w, "%-*s %-*s %-*s %s %s %-*s %s %s %-*s\n",
			tableColNamespace, truncate(pu.Namespace, tableColNamespace),
			tableColPod, truncate(pu.Name, tableColPod),
			tableColCPUUsage, f.unitFormatter.FormatCPU(pu.CPU.Usage.MilliValue()),
			f.colorizer.FormatPercent(pu.CPU.RequestPercent, tableColPercent),
			f.colorizer.FormatPercent(pu.CPU.LimitPercent, tableColPercent),
			tableColMemUsage, f.unitFormatter.FormatMemory(pu.Memory.Usage.Value()),
			f.colorizer.FormatPercent(pu.Memory.RequestPercent, tableColPercent),
			f.colorizer.FormatPercent(pu.Memory.LimitPercent, tableColPercent),
			tableColNode, truncate(pu.Node, tableColNode),
		)
	}

	return nil
}

// truncate truncates string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
