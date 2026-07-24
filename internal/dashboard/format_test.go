package dashboard

import (
	"strings"
	"testing"
)

func TestByteRateAndUsageBar(t *testing.T) {
	rate := float64(1536)
	if got := ByteRate(&rate); got != "1.5 KiB/s" {
		t.Fatalf("ByteRate() = %q", got)
	}
	percent := 50.0
	if got := UsageBar(&percent, 10); got != "[####----]" {
		t.Fatalf("UsageBar() = %q", got)
	}
	if got := UsageBar(nil, 6); got != "[????]" {
		t.Fatalf("UsageBar(nil) = %q", got)
	}
}

func TestSparklineIsBoundedAndASCIISafe(t *testing.T) {
	got := Sparkline([]float64{0, 25, 50, 75, 100}, 4)
	if len(got) != 4 || strings.ContainsAny(got, "▁▂▃▄▅▆▇█") {
		t.Fatalf("Sparkline() = %q", got)
	}
}

func TestTruncate(t *testing.T) {
	if got := Truncate("abcdefgh", 6); got != "abc..." {
		t.Fatalf("Truncate() = %q", got)
	}
}
