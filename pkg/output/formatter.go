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
	case "yaml":
		return &YAMLFormatter{}
	case "wide":
		return &WideFormatter{}
	default:
		return &TableFormatter{}
	}
}

// StructuredOutput is the structured output format used by JSON and YAML formatters
type StructuredOutput struct {
	Items []StructuredPodUsage `json:"items" yaml:"items"`
}

// StructuredPodUsage represents a pod's resource usage in structured format
type StructuredPodUsage struct {
	Namespace string                  `json:"namespace" yaml:"namespace"`
	Pod       string                  `json:"pod" yaml:"pod"`
	Node      string                  `json:"node" yaml:"node"`
	CPU       StructuredResourceUsage `json:"cpu" yaml:"cpu"`
	Memory    StructuredResourceUsage `json:"memory" yaml:"memory"`
}

// StructuredResourceUsage represents CPU or Memory usage in structured format
type StructuredResourceUsage struct {
	Usage          string  `json:"usage" yaml:"usage"`
	Requests       *string `json:"requests" yaml:"requests"`
	Limits         *string `json:"limits" yaml:"limits"`
	RequestPercent *int    `json:"requestPercent" yaml:"requestPercent"`
	LimitPercent   *int    `json:"limitPercent" yaml:"limitPercent"`
}

// toStructuredOutput converts pod usages to structured output format
func toStructuredOutput(podUsages []calculator.PodUsage) StructuredOutput {
	output := StructuredOutput{
		Items: make([]StructuredPodUsage, 0, len(podUsages)),
	}

	for _, pu := range podUsages {
		structuredPod := StructuredPodUsage{
			Namespace: pu.Namespace,
			Pod:       pu.Name,
			Node:      pu.Node,
			CPU:       toStructuredResourceUsage(pu.CPU),
			Memory:    toStructuredResourceUsage(pu.Memory),
		}
		output.Items = append(output.Items, structuredPod)
	}

	return output
}

// toStructuredResourceUsage converts ResourceUsage to StructuredResourceUsage
func toStructuredResourceUsage(ru calculator.ResourceUsage) StructuredResourceUsage {
	result := StructuredResourceUsage{
		Usage:          ru.Usage.String(),
		RequestPercent: ru.RequestPercent,
		LimitPercent:   ru.LimitPercent,
	}

	if ru.Requests != nil {
		s := ru.Requests.String()
		result.Requests = &s
	}

	if ru.Limits != nil {
		s := ru.Limits.String()
		result.Limits = &s
	}

	return result
}
