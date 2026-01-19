package output

import (
	"encoding/json"
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
)

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// Format writes pod usages as JSON
func (f *JSONFormatter) Format(w io.Writer, podUsages []calculator.PodUsage) error {
	output := toStructuredOutput(podUsages)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
