package system

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestSnapshotJSONContract(t *testing.T) {
	timestamp := time.Date(2026, 7, 24, 1, 30, 0, 0, time.FixedZone("SAST", 7200))
	snapshot := Snapshot{
		SchemaVersion: SnapshotSchemaVersion,
		Timestamp:     timestamp,
		TimeZone:      TimeZone{Name: "SAST", UTCOffsetSeconds: 7200},
		Filesystems:   []Filesystem{},
		Networks:      []NetworkInterface{},
		Warnings:      []Warning{},
	}
	data, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, `"schema_version":1`) {
		t.Fatalf("JSON = %s", text)
	}
	if !strings.Contains(text, `"timestamp":"2026-07-24T01:30:00+02:00"`) {
		t.Fatalf("timestamp is not RFC 3339: %s", text)
	}
	if !strings.Contains(text, `"total_bytes":null`) {
		t.Fatalf("unavailable numeric value is not null: %s", text)
	}
	if strings.Contains(text, "\x1b") {
		t.Fatalf("JSON contains ANSI escape: %q", text)
	}
}

func TestUtilization(t *testing.T) {
	tests := []struct {
		used  uint64
		total uint64
		want  float64
	}{
		{used: 50, total: 100, want: 50},
		{used: 0, total: 0, want: 0},
		{used: 200, total: 100, want: 100},
	}
	for _, test := range tests {
		if got := Utilization(test.used, test.total); got != test.want {
			t.Fatalf("Utilization(%d, %d) = %v, want %v", test.used, test.total, got, test.want)
		}
	}
}
