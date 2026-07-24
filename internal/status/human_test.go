package status

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/setup-env/app/internal/diagnostics"
	"github.com/setup-env/app/internal/directory"
	"github.com/setup-env/app/internal/system"
)

func TestWriteHuman(t *testing.T) {
	total := uint64(8 * 1024 * 1024 * 1024)
	used := uint64(4 * 1024 * 1024 * 1024)
	available := total - used
	percent := 50.0
	cores := 4
	logical := 8
	uptime := uint64(3661)
	authenticated := false
	snapshot := system.Snapshot{
		SchemaVersion: system.SnapshotSchemaVersion,
		Timestamp:     time.Date(2026, 7, 24, 1, 30, 0, 0, time.FixedZone("SAST", 7200)),
		TimeZone:      system.TimeZone{Name: "SAST", UTCOffsetSeconds: 7200},
		Host:          system.Host{Hostname: "test-host", UptimeSeconds: &uptime},
		OperatingSystem: system.OperatingSystem{
			OS:            "linux",
			DisplayName:   "Ubuntu",
			Version:       "26.04",
			KernelVersion: "6.8.0",
			Architecture:  "amd64",
		},
		User: system.User{Username: "tester"},
		CPU: system.CPU{
			Model:              "Test CPU",
			PhysicalCores:      &cores,
			LogicalCPUs:        &logical,
			UtilizationPercent: &percent,
		},
		Memory: system.Memory{
			TotalBytes:         &total,
			UsedBytes:          &used,
			AvailableBytes:     &available,
			UtilizationPercent: &percent,
		},
		Filesystems: []system.Filesystem{
			{Mountpoint: "/", Type: "ext4", TotalBytes: &total, UsedBytes: &used, AvailableBytes: &available, UtilizationPercent: &percent},
		},
		Networks: []system.NetworkInterface{
			{Name: "eth0", Status: "up", Addresses: []system.NetworkAddress{{Address: "10.0.0.9", Family: "ipv4", PrefixLength: 24}}},
		},
		Development: system.Development{
			Root:                diagnostics.DevelopmentRootStatus{Path: "/home/tester/dev", Exists: true, Writable: true},
			Directory:           directory.Context{Type: directory.TypeRepository},
			Git:                 diagnostics.ToolStatus{Available: true, Version: "git version 2.50.0"},
			GitHubCLI:           diagnostics.ToolStatus{Available: true, Version: "gh version 2.80.0"},
			GitHubAuthenticated: &authenticated,
		},
		Diagnostics: system.DiagnosticSummary{
			Health:       system.HealthWarning,
			WarningCount: 1,
			Checks: []diagnostics.Check{
				{Name: "github-auth", Status: diagnostics.StatusWarn, Message: "authentication requires renewal"},
			},
		},
	}
	var output bytes.Buffer
	if err := WriteHuman(&output, snapshot); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{
		"SETUP ENV SYSTEM STATUS",
		"2026-07-24 01:30:00 SAST",
		"Ubuntu",
		"Test CPU",
		"4.0 GiB",
		"eth0",
		"10.0.0.9",
		"/home/tester/dev",
		"warning",
		"github-auth: authentication requires renewal",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("output missing %q:\n%s", expected, text)
		}
	}
	if strings.Contains(text, "\x1b") {
		t.Fatalf("output contains ANSI escape: %q", text)
	}
}
