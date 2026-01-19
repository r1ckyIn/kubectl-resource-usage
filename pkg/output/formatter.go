package output

import (
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
)

// Formatter is the interface for output formatters
type Formatter interface {
	Format(w io.Writer, podUsages []calculator.PodUsage) error
}

// NewFormatter creates a formatter based on the format type
func NewFormatter(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	default:
		return &TableFormatter{}
	}
}
