package output

import (
	"fmt"
)

// Unit represents the unit for resource display
type Unit string

const (
	UnitAuto Unit = "auto"
	// Memory units
	UnitKi Unit = "Ki"
	UnitMi Unit = "Mi"
	UnitGi Unit = "Gi"
	// CPU units
	UnitMillicores Unit = "m"
	UnitCores      Unit = "cores"
)

// ValidUnits returns the list of valid unit options
func ValidUnits() []string {
	return []string{"auto", "Ki", "Mi", "Gi", "m", "cores"}
}

// IsValidUnit checks if the given unit is valid
func IsValidUnit(u string) bool {
	for _, valid := range ValidUnits() {
		if u == valid {
			return true
		}
	}
	return false
}

// UnitFormatter handles unit conversion for display
type UnitFormatter struct {
	unit Unit
}

// NewUnitFormatter creates a new UnitFormatter
func NewUnitFormatter(unit string) *UnitFormatter {
	return &UnitFormatter{unit: Unit(unit)}
}

// FormatCPU formats CPU value based on the unit setting
func (f *UnitFormatter) FormatCPU(milliCores int64) string {
	switch f.unit {
	case UnitCores:
		cores := float64(milliCores) / 1000
		if cores >= 1 {
			return fmt.Sprintf("%.1f", cores)
		}
		return fmt.Sprintf("%.2f", cores)
	case UnitMillicores:
		return fmt.Sprintf("%dm", milliCores)
	default:
		return fmt.Sprintf("%dm", milliCores)
	}
}

// FormatMemory formats memory value based on the unit setting
func (f *UnitFormatter) FormatMemory(bytes int64) string {
	const (
		Ki = 1024
		Mi = Ki * 1024
		Gi = Mi * 1024
	)

	switch f.unit {
	case UnitKi:
		return fmt.Sprintf("%dKi", bytes/Ki)
	case UnitMi:
		return fmt.Sprintf("%dMi", bytes/Mi)
	case UnitGi:
		if bytes >= Gi {
			return fmt.Sprintf("%dGi", bytes/Gi)
		}
		gb := float64(bytes) / float64(Gi)
		return fmt.Sprintf("%.2fGi", gb)
	default:
		switch {
		case bytes >= Gi:
			return fmt.Sprintf("%dGi", bytes/Gi)
		case bytes >= Mi:
			return fmt.Sprintf("%dMi", bytes/Mi)
		case bytes >= Ki:
			return fmt.Sprintf("%dKi", bytes/Ki)
		default:
			return fmt.Sprintf("%d", bytes)
		}
	}
}
