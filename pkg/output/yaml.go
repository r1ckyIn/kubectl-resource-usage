package output

import (
	"io"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
	"gopkg.in/yaml.v3"
)

// YAMLFormatter formats output as YAML
type YAMLFormatter struct{}

// Format writes pod usages as YAML
func (f *YAMLFormatter) Format(w io.Writer, podUsages []calculator.PodUsage) error {
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

	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(output)
}
