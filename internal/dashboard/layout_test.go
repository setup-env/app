package dashboard

import (
	"strings"
	"testing"
	"time"

	"github.com/setup-env/app/internal/diagnostics"
	"github.com/setup-env/app/internal/system"
)

func TestRenderWideCompactSmallAndHelpLayouts(t *testing.T) {
	snapshot := dashboardTestSnapshot()
	tests := []struct {
		name     string
		state    ViewState
		contains string
	}{
		{name: "wide", state: ViewState{Snapshot: snapshot, Clock: snapshot.Timestamp, Width: 120, Height: 40}, contains: "Development and health"},
		{name: "compact", state: ViewState{Snapshot: snapshot, Clock: snapshot.Timestamp, Width: 70, Height: 20}, contains: "HEALTH healthy"},
		{name: "small", state: ViewState{Snapshot: snapshot, Clock: snapshot.Timestamp, Width: 30, Height: 8}, contains: "Terminal too small"},
		{name: "help", state: ViewState{Snapshot: snapshot, Clock: snapshot.Timestamp, Width: 70, Height: 20, Help: true}, contains: "Ctrl+C"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Render(test.state)
			if !strings.Contains(got, test.contains) {
				t.Fatalf("render = %q", got)
			}
			for _, line := range strings.Split(got, "\n") {
				if len([]rune(line)) > max(test.state.Width, len([]rune("Resize to at least 40x12."))) {
					t.Fatalf("line exceeds width %d: %q", test.state.Width, line)
				}
			}
		})
	}
}

func TestFilesystemLayoutIsDeterministicAndReportsHiddenRows(t *testing.T) {
	percent := 25.0
	total, used := uint64(100), uint64(25)
	filesystems := []system.Filesystem{
		{Mountpoint: "/", TotalBytes: &total, UsedBytes: &used, UtilizationPercent: &percent},
		{Mountpoint: "/home", TotalBytes: &total, UsedBytes: &used, UtilizationPercent: &percent},
	}
	lines := filesystemLines(filesystems, 80, 1)
	if len(lines) != 2 || !strings.HasPrefix(lines[0], "/") || !strings.Contains(lines[1], "1 more") {
		t.Fatalf("lines = %#v", lines)
	}
}

func dashboardTestSnapshot() system.Snapshot {
	percent := 50.0
	total, used, available, uptime := uint64(8<<30), uint64(4<<30), uint64(4<<30), uint64(3600)
	cores := 4
	authenticated := true
	return system.Snapshot{
		SchemaVersion: system.SnapshotSchemaVersion,
		Timestamp:     time.Date(2026, 7, 24, 3, 0, 0, 0, time.FixedZone("SAST", 7200)),
		Host:          system.Host{Hostname: "example", UptimeSeconds: &uptime},
		OperatingSystem: system.OperatingSystem{
			DisplayName:  "Ubuntu",
			Version:      "26.04",
			Architecture: "amd64",
		},
		CPU: system.CPU{
			Model:              "Example CPU",
			PhysicalCores:      &cores,
			LogicalCPUs:        &cores,
			UtilizationPercent: &percent,
		},
		Memory: system.Memory{
			TotalBytes:         &total,
			UsedBytes:          &used,
			AvailableBytes:     &available,
			UtilizationPercent: &percent,
		},
		Filesystems: []system.Filesystem{{
			Mountpoint:         "/",
			TotalBytes:         &total,
			UsedBytes:          &used,
			UtilizationPercent: &percent,
		}},
		Networks: []system.NetworkInterface{{Name: "eth0", Status: "up"}},
		Development: system.Development{
			Root:                diagnostics.DevelopmentRootStatus{Path: "/home/test/dev", Writable: true},
			Git:                 diagnostics.ToolStatus{Available: true},
			GitHubCLI:           diagnostics.ToolStatus{Available: true},
			GitHubAuthenticated: &authenticated,
		},
		Diagnostics: system.DiagnosticSummary{Health: system.HealthHealthy},
	}
}
