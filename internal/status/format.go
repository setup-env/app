package status

import (
	"fmt"
	"math"
	"time"
)

const unavailable = "unavailable"

func Bytes(value *uint64) string {
	if value == nil {
		return unavailable
	}
	const (
		kib = uint64(1024)
		mib = kib * 1024
		gib = mib * 1024
		tib = gib * 1024
	)
	switch {
	case *value >= tib:
		return fmt.Sprintf("%.1f TiB", float64(*value)/float64(tib))
	case *value >= gib:
		return fmt.Sprintf("%.1f GiB", float64(*value)/float64(gib))
	case *value >= mib:
		return fmt.Sprintf("%.1f MiB", float64(*value)/float64(mib))
	case *value >= kib:
		return fmt.Sprintf("%.1f KiB", float64(*value)/float64(kib))
	default:
		return fmt.Sprintf("%d B", *value)
	}
}

func Percentage(value *float64) string {
	if value == nil || math.IsNaN(*value) || math.IsInf(*value, 0) {
		return unavailable
	}
	return fmt.Sprintf("%.1f%%", *value)
}

func Duration(seconds *uint64) string {
	if seconds == nil {
		return unavailable
	}
	duration := time.Duration(*seconds) * time.Second
	days := duration / (24 * time.Hour)
	duration %= 24 * time.Hour
	hours := duration / time.Hour
	duration %= time.Hour
	minutes := duration / time.Minute
	remainingSeconds := duration % time.Minute / time.Second
	if days > 0 {
		return fmt.Sprintf("%d days %02d:%02d:%02d", days, hours, minutes, remainingSeconds)
	}
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, remainingSeconds)
}

func Timestamp(value time.Time) string {
	if value.IsZero() {
		return unavailable
	}
	return value.Format("2006-01-02 15:04:05 MST")
}

func UTCOffset(seconds int) string {
	sign := "+"
	if seconds < 0 {
		sign = "-"
		seconds = -seconds
	}
	hours := seconds / 3600
	minutes := seconds % 3600 / 60
	return fmt.Sprintf("UTC%s%02d:%02d", sign, hours, minutes)
}

func Text(value string) string {
	if value == "" {
		return unavailable
	}
	return value
}

func YesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func OptionalBool(value *bool) string {
	if value == nil {
		return unavailable
	}
	return YesNo(*value)
}
