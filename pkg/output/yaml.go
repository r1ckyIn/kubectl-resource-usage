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
	output := toStructuredOutput(podUsages)
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(output)
}
