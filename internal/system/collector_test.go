package system

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/setup-env/app/internal/diagnostics"
)

type fakeSection struct {
	name    string
	collect func(*Snapshot)
	err     error
}

func (f fakeSection) Name() string { return f.name }

func (f fakeSection) Collect(_ context.Context, snapshot *Snapshot) error {
	if f.collect != nil {
		f.collect(snapshot)
	}
	return f.err
}

func TestCollectorPreservesPartialResultsAndWarnings(t *testing.T) {
	now := time.Date(2026, 7, 24, 1, 30, 0, 0, time.UTC)
	collector := Collector{
		Now: func() time.Time { return now },
		Sections: []SectionCollector{
			fakeSection{name: "host", collect: func(snapshot *Snapshot) { snapshot.Host.Hostname = "test-host" }},
			fakeSection{name: "cpu", err: errors.New("counter unavailable")},
		},
	}
	snapshot, err := collector.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Host.Hostname != "test-host" || len(snapshot.Warnings) != 1 {
		t.Fatalf("Collect() = %#v", snapshot)
	}
	if snapshot.Diagnostics.Health != HealthWarning || snapshot.Diagnostics.WarningCount != 1 {
		t.Fatalf("health = %#v", snapshot.Diagnostics)
	}
}

func TestCollectorReturnsErrorWhenEverySectionFails(t *testing.T) {
	collector := Collector{Sections: []SectionCollector{
		fakeSection{name: "host", err: errors.New("unavailable")},
	}}
	snapshot, err := collector.Collect(context.Background())
	if err == nil || !strings.Contains(err.Error(), "every section failed") {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(snapshot.Warnings) != 1 {
		t.Fatalf("warnings = %#v", snapshot.Warnings)
	}
}

func TestHealthUsesDiagnosticFailures(t *testing.T) {
	snapshot := Snapshot{
		Diagnostics: DiagnosticSummary{Checks: []diagnostics.Check{
			{Name: "git", Status: diagnostics.StatusPass},
			{Name: "development-root", Status: diagnostics.StatusFail},
		}},
	}
	finalizeHealth(&snapshot)
	if snapshot.Diagnostics.Health != HealthUnhealthy || snapshot.Diagnostics.FailureCount != 1 {
		t.Fatalf("health = %#v", snapshot.Diagnostics)
	}
}

func TestCollectorHonorsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (Collector{}).Collect(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Collect() error = %v, want context.Canceled", err)
	}
}
