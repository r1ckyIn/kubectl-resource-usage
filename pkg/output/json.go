package output

import (
	"encoding/json"
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
)

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// JSONOutput is the JSON output structure
type JSONOutput struct {
	Items []JSONPodUsage `json:"items"`
}

// JSONPodUsage represents a pod's resource usage in JSON format
type JSONPodUsage struct {
	Namespace string            `json:"namespace"`
	Pod       string            `json:"pod"`
	Node      string            `json:"node"`
	CPU       JSONResourceUsage `json:"cpu"`
	Memory    JSONResourceUsage `json:"memory"`
}

// JSONResourceUsage represents CPU or Memory usage in JSON format
type JSONResourceUsage struct {
	Usage          string `json:"usage"`
	Requests       *string `json:"requests"`
	Limits         *string `json:"limits"`
	RequestPercent *int    `json:"requestPercent"`
	LimitPercent   *int    `json:"limitPercent"`
}

// Format writes pod usages as JSON
func (f *JSONFormatter) Format(w io.Writer, podUsages []calculator.PodUsage) error {
	output := JSONOutput{
		Items: make([]JSONPodUsage, 0, len(podUsages)),
	}

	for _, pu := range podUsages {
		jsonPod := JSONPodUsage{
			Namespace: pu.Namespace,
			Pod:       pu.Name,
			Node:      pu.Node,
			CPU:       toJSONResourceUsage(pu.CPU),
			Memory:    toJSONResourceUsage(pu.Memory),
		}
		output.Items = append(output.Items, jsonPod)
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// toJSONResourceUsage converts ResourceUsage to JSONResourceUsage
func toJSONResourceUsage(ru calculator.ResourceUsage) JSONResourceUsage {
	result := JSONResourceUsage{
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
