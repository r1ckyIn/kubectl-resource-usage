package calculator

// FilterOptions contains options for filtering pod usages
type FilterOptions struct {
	Above    int    // Filter pods with usage >= Above%, -1 means not set
	Below    int    // Filter pods with usage <= Below%, -1 means not set
	NoLimits bool   // Filter pods without limits set
	Field    string // Field to filter by: "cpu" or "memory"
}

// NewFilterOptions creates a FilterOptions with default values
func NewFilterOptions() FilterOptions {
	return FilterOptions{
		Above: -1,
		Below: -1,
		Field: "memory",
	}
}

// FilterPodUsages filters pod usages based on the provided options
func FilterPodUsages(pods []PodUsage, opts FilterOptions) []PodUsage {
	if opts.Above == -1 && opts.Below == -1 && !opts.NoLimits {
		return pods
	}

	var result []PodUsage
	for _, pod := range pods {
		if matchesFilter(pod, opts) {
			result = append(result, pod)
		}
	}
	return result
}

// matchesFilter checks if a pod matches the filter criteria
func matchesFilter(pod PodUsage, opts FilterOptions) bool {
	// Handle --no-limits filter
	if opts.NoLimits {
		if pod.CPU.Limits == nil || pod.Memory.Limits == nil {
			return true
		}
		return false
	}

	// Get the percentage based on the field
	var percent *int
	if opts.Field == "cpu" {
		percent = pod.CPU.LimitPercent
	} else {
		percent = pod.Memory.LimitPercent
	}

	// If no limit percentage is available, exclude from filter results
	if percent == nil {
		return false
	}

	// Check above threshold
	if opts.Above != -1 && *percent < opts.Above {
		return false
	}

	// Check below threshold
	if opts.Below != -1 && *percent > opts.Below {
		return false
	}

	return true
}
