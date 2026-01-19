package output

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// ColorMode represents the color output mode
type ColorMode string

const (
	ColorModeAuto   ColorMode = "auto"
	ColorModeAlways ColorMode = "always"
	ColorModeNever  ColorMode = "never"
)

// Colorizer handles colorizing output based on usage percentage
type Colorizer struct {
	enabled bool
}

// NewColorizer creates a new Colorizer based on the color mode
func NewColorizer(mode ColorMode) *Colorizer {
	var enabled bool
	switch mode {
	case ColorModeAlways:
		enabled = true
	case ColorModeNever:
		enabled = false
	default:
		enabled = isTerminal()
	}
	return &Colorizer{enabled: enabled}
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// FormatPercent formats percentage with color based on value
// High usage (>=80%) -> Red
// Medium usage (50-79%) -> Yellow
// Low usage (<50%) -> Green
// N/A -> No color
func (c *Colorizer) FormatPercent(p *int, width int) string {
	if p == nil {
		return fmt.Sprintf("%-*s", width, "N/A")
	}

	percentStr := fmt.Sprintf("%d%%", *p)

	if !c.enabled {
		return fmt.Sprintf("%-*s", width, percentStr)
	}

	var color string
	switch {
	case *p >= 80:
		color = colorRed
	case *p >= 50:
		color = colorYellow
	default:
		color = colorGreen
	}

	return fmt.Sprintf("%s%-*s%s", color, width, percentStr, colorReset)
}

// Enabled returns whether colorization is enabled
func (c *Colorizer) Enabled() bool {
	return c.enabled
}
