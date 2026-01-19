package cmd

import (
	"context"
	"fmt"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
	"github.com/r1ckyIn/kubectl-resource-usage/pkg/collector"
	"github.com/r1ckyIn/kubectl-resource-usage/pkg/output"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// ResourceUsageOptions contains the options for the resource-usage command
type ResourceUsageOptions struct {
	configFlags *genericclioptions.ConfigFlags
	genericclioptions.IOStreams

	selector  string
	sortBy    string
	ascending bool
	output    string

	// Filter options
	above    int
	below    int
	noLimits bool
}

// NewResourceUsageOptions creates a new ResourceUsageOptions with default values
func NewResourceUsageOptions(streams genericclioptions.IOStreams) *ResourceUsageOptions {
	return &ResourceUsageOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
		output:      "table",
		above:       -1,
		below:       -1,
	}
}

// NewCmdResourceUsage creates the resource-usage command
func NewCmdResourceUsage(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewResourceUsageOptions(streams)

	cmd := &cobra.Command{
		Use:   "resource-usage",
		Short: "Display Pod CPU/Memory usage percentage relative to requests and limits",
		Long: `Display Pod-level CPU/Memory resource usage rates.
Unlike 'kubectl top pods', this plugin calculates usage percentages
relative to requests and limits, helping SREs quickly identify resource issues.`,
		Example: `  # View resource usage for all namespaces
  kubectl resource-usage

  # Filter by namespace
  kubectl resource-usage -n payment

  # Filter by label
  kubectl resource-usage -l app=api

  # Sort by memory usage (descending)
  kubectl resource-usage --sort memory

  # Output as JSON
  kubectl resource-usage -o json

  # Show pods with memory usage >= 80%
  kubectl resource-usage --above 80

  # Show pods with CPU usage <= 50%
  kubectl resource-usage --below 50 --sort cpu

  # Show pods without limits configured
  kubectl resource-usage --no-limits

  # Output as YAML or wide format
  kubectl resource-usage -o yaml
  kubectl resource-usage -o wide`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run(cmd.Context())
		},
	}

	// Add kubectl config flags (--kubeconfig, --context, --namespace, etc.)
	o.configFlags.AddFlags(cmd.Flags())

	// Add custom flags
	cmd.Flags().StringVarP(&o.selector, "selector", "l", "", "Filter by label selector (e.g., app=api)")
	cmd.Flags().StringVar(&o.sortBy, "sort", "", "Sort by field: cpu or memory")
	cmd.Flags().BoolVar(&o.ascending, "asc", false, "Sort in ascending order (default: descending)")
	cmd.Flags().StringVarP(&o.output, "output", "o", "table", "Output format: table, json, yaml, or wide")

	// Filter flags
	cmd.Flags().IntVar(&o.above, "above", -1, "Show pods with usage >= N% (uses --sort field, default: memory)")
	cmd.Flags().IntVar(&o.below, "below", -1, "Show pods with usage <= N% (uses --sort field, default: memory)")
	cmd.Flags().BoolVar(&o.noLimits, "no-limits", false, "Show pods without limits configured")

	return cmd
}

// Complete fills in any fields not set by flags
func (o *ResourceUsageOptions) Complete(cmd *cobra.Command) error {
	return nil
}

// Validate validates the options
func (o *ResourceUsageOptions) Validate() error {
	if o.selector != "" {
		if _, err := labels.Parse(o.selector); err != nil {
			return fmt.Errorf("invalid label selector: %w", err)
		}
	}
	if o.sortBy != "" && o.sortBy != "cpu" && o.sortBy != "memory" {
		return fmt.Errorf("invalid sort field: %s (must be 'cpu' or 'memory')", o.sortBy)
	}
	validOutputs := map[string]bool{"table": true, "json": true, "yaml": true, "wide": true}
	if !validOutputs[o.output] {
		return fmt.Errorf("invalid output format: %s (must be 'table', 'json', 'yaml', or 'wide')", o.output)
	}
	if o.above != -1 && (o.above < 0 || o.above > 100) {
		return fmt.Errorf("invalid --above value: %d (must be between 0 and 100)", o.above)
	}
	if o.below != -1 && (o.below < 0 || o.below > 100) {
		return fmt.Errorf("invalid --below value: %d (must be between 0 and 100)", o.below)
	}
	if o.above != -1 && o.below != -1 && o.above > o.below {
		return fmt.Errorf("--above (%d) cannot be greater than --below (%d)", o.above, o.below)
	}
	return nil
}

// Run executes the resource-usage command
func (o *ResourceUsageOptions) Run(ctx context.Context) error {
	// Create REST config from flags
	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("failed to create REST config: %w", err)
	}

	// Get namespace from config flags (empty string means all namespaces)
	namespace := ""
	if o.configFlags.Namespace != nil && *o.configFlags.Namespace != "" {
		namespace = *o.configFlags.Namespace
	}

	// Create collectors
	metricsCollector, err := collector.NewMetricsCollector(restConfig)
	if err != nil {
		return fmt.Errorf("failed to create metrics collector: %w", err)
	}

	podCollector, err := collector.NewPodCollector(restConfig)
	if err != nil {
		return fmt.Errorf("failed to create pod collector: %w", err)
	}

	// Fetch pod metrics
	podMetrics, err := metricsCollector.GetPodMetrics(ctx, namespace)
	if err != nil {
		return fmt.Errorf("failed to get pod metrics: %w", err)
	}

	// Fetch pods with label selector
	pods, err := podCollector.GetPods(ctx, namespace, o.selector)
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	// Build pod map for quick lookup
	podMap := make(map[string]int)
	for i, pod := range pods.Items {
		key := pod.Namespace + "/" + pod.Name
		podMap[key] = i
	}

	// Calculate usage for each pod
	var podUsages []calculator.PodUsage
	for _, pm := range podMetrics.Items {
		key := pm.Namespace + "/" + pm.Name
		podIndex, exists := podMap[key]
		if !exists {
			continue
		}
		podUsage := calculator.CalculatePodUsage(pm, pods.Items[podIndex])
		podUsages = append(podUsages, podUsage)
	}

	// Apply filters
	filterField := o.sortBy
	if filterField == "" {
		filterField = "memory"
	}
	filterOpts := calculator.FilterOptions{
		Above:    o.above,
		Below:    o.below,
		NoLimits: o.noLimits,
		Field:    filterField,
	}
	podUsages = calculator.FilterPodUsages(podUsages, filterOpts)

	// Handle empty results
	if len(podUsages) == 0 {
		fmt.Fprintln(o.Out, "No pods found matching the criteria")
		return nil
	}

	// Sort if requested
	if o.sortBy != "" {
		calculator.SortPodUsages(podUsages, o.sortBy, o.ascending)
	}

	// Output using formatter
	formatter := output.NewFormatter(o.output)
	return formatter.Format(o.Out, podUsages)
}
