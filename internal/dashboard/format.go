package dashboard

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

func ByteRate(value *float64) string {
	if value == nil || math.IsNaN(*value) || math.IsInf(*value, 0) || *value < 0 {
		return "unavailable"
	}
	const (
		kib = 1024
		mib = kib * 1024
		gib = mib * 1024
	)
	switch {
	case *value >= gib:
		return fmt.Sprintf("%.1f GiB/s", *value/gib)
	case *value >= mib:
		return fmt.Sprintf("%.1f MiB/s", *value/mib)
	case *value >= kib:
		return fmt.Sprintf("%.1f KiB/s", *value/kib)
	default:
		return fmt.Sprintf("%.0f B/s", *value)
	}
}

func UsageBar(percent *float64, width int) string {
	if width < 3 {
		return ""
	}
	inside := width - 2
	if percent == nil {
		return "[" + strings.Repeat("?", inside) + "]"
	}
	value := math.Max(0, math.Min(100, *percent))
	filled := int(math.Round(value / 100 * float64(inside)))
	return "[" + strings.Repeat("#", filled) + strings.Repeat("-", inside-filled) + "]"
}

func Sparkline(values []float64, width int) string {
	if width <= 0 || len(values) == 0 {
		return ""
	}
	const levels = " .:-=+*#%@"
	start := 0
	if len(values) > width {
		start = len(values) - width
	}
	var result strings.Builder
	for _, value := range values[start:] {
		value = math.Max(0, math.Min(100, value))
		index := int(math.Round(value / 100 * float64(len(levels)-1)))
		result.WriteByte(levels[index])
	}
	if padding := width - result.Len(); padding > 0 {
		return strings.Repeat(" ", padding) + result.String()
	}
	return result.String()
}

func Truncate(value string, width int) string {
	if width <= 0 {
		return ""
	}
	if utf8.RuneCountInString(value) <= width {
		return value
	}
	if width <= 3 {
		return string([]rune(value)[:width])
	}
	runes := []rune(value)
	return string(runes[:width-3]) + "..."
}
