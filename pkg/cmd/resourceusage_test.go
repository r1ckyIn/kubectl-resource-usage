package cmd

import (
	"strings"
	"testing"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestResourceUsageOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ResourceUsageOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid default options",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid output format",
			opts: &ResourceUsageOptions{
				output:   "invalid",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "invalid output format",
		},
		{
			name: "invalid color mode",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "invalid",
				unit:     "auto",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "invalid color mode",
		},
		{
			name: "invalid sort field",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				sortBy:   "invalid",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "invalid sort field",
		},
		{
			name: "above greater than below",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    80,
				below:    50,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "--above (80) cannot be greater than --below (50)",
		},
		{
			name: "above out of range",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    150,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "invalid --above value",
		},
		{
			name: "below out of range",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    150,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "invalid --below value",
		},
		{
			name: "no-limits with above",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    50,
				below:    -1,
				noLimits: true,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "--no-limits cannot be used with --above or --below",
		},
		{
			name: "no-limits with below",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    50,
				noLimits: true,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "--no-limits cannot be used with --above or --below",
		},
		{
			name: "watch mode with json output",
			opts: &ResourceUsageOptions{
				output:   "json",
				color:    "auto",
				unit:     "auto",
				watch:    true,
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "watch mode is not supported with json output format",
		},
		{
			name: "watch mode with yaml output",
			opts: &ResourceUsageOptions{
				output:   "yaml",
				color:    "auto",
				unit:     "auto",
				watch:    true,
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: true,
			errMsg:  "watch mode is not supported with yaml output format",
		},
		{
			name: "interval too short",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    -1,
				interval: 500 * time.Millisecond,
			},
			wantErr: true,
			errMsg:  "interval must be at least 1 second",
		},
		{
			name: "valid label selector",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				selector: "app=test",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid json output",
			opts: &ResourceUsageOptions{
				output:   "json",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid yaml output",
			opts: &ResourceUsageOptions{
				output:   "yaml",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid wide output",
			opts: &ResourceUsageOptions{
				output:   "wide",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    -1,
				interval: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid above and below",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    20,
				below:    80,
				interval: 2 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid no-limits alone",
			opts: &ResourceUsageOptions{
				output:   "table",
				color:    "auto",
				unit:     "auto",
				above:    -1,
				below:    -1,
				noLimits: true,
				interval: 2 * time.Second,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestNewResourceUsageOptions(t *testing.T) {
	streams := genericclioptions.IOStreams{}
	opts := NewResourceUsageOptions(streams)

	if opts.output != "table" {
		t.Errorf("expected output 'table', got %q", opts.output)
	}
	if opts.color != "auto" {
		t.Errorf("expected color 'auto', got %q", opts.color)
	}
	if opts.unit != "auto" {
		t.Errorf("expected unit 'auto', got %q", opts.unit)
	}
	if opts.above != -1 {
		t.Errorf("expected above -1, got %d", opts.above)
	}
	if opts.below != -1 {
		t.Errorf("expected below -1, got %d", opts.below)
	}
}
