package status

import (
	"testing"
	"time"
)

func TestBytes(t *testing.T) {
	tests := []struct {
		value *uint64
		want  string
	}{
		{value: nil, want: "unavailable"},
		{value: uint64Value(0), want: "0 B"},
		{value: uint64Value(1024), want: "1.0 KiB"},
		{value: uint64Value(1024 * 1024), want: "1.0 MiB"},
		{value: uint64Value(3 * 1024 * 1024 * 1024), want: "3.0 GiB"},
		{value: uint64Value(2 * 1024 * 1024 * 1024 * 1024), want: "2.0 TiB"},
	}
	for _, test := range tests {
		if got := Bytes(test.value); got != test.want {
			t.Fatalf("Bytes() = %q, want %q", got, test.want)
		}
	}
}

func TestPercentage(t *testing.T) {
	if got := Percentage(nil); got != "unavailable" {
		t.Fatalf("Percentage(nil) = %q", got)
	}
	value := 43.75
	if got := Percentage(&value); got != "43.8%" {
		t.Fatalf("Percentage() = %q", got)
	}
}

func TestDuration(t *testing.T) {
	if got := Duration(nil); got != "unavailable" {
		t.Fatalf("Duration(nil) = %q", got)
	}
	value := uint64(4*24*60*60 + 3*60*60 + 12*60 + 18)
	if got := Duration(&value); got != "4 days 03:12:18" {
		t.Fatalf("Duration() = %q", got)
	}
}

func TestTimestampAndOffset(t *testing.T) {
	value := time.Date(2026, 7, 24, 1, 30, 0, 0, time.FixedZone("SAST", 7200))
	if got := Timestamp(value); got != "2026-07-24 01:30:00 SAST" {
		t.Fatalf("Timestamp() = %q", got)
	}
	if got := UTCOffset(7200); got != "UTC+02:00" {
		t.Fatalf("UTCOffset() = %q", got)
	}
}

func uint64Value(value uint64) *uint64 { return &value }
