package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/setup-env/app/internal/app"
	"github.com/setup-env/app/internal/catalog"
	"github.com/setup-env/app/internal/config"
	"github.com/setup-env/app/internal/paths"
	"github.com/setup-env/app/internal/system"
)

type fakeSnapshotCollector struct {
	snapshot system.Snapshot
	err      error
}

func (f fakeSnapshotCollector) Collect(context.Context) (system.Snapshot, error) {
	return f.snapshot, f.err
}

func TestNoArgumentsDisplaysStaticStatus(t *testing.T) {
	var output bytes.Buffer
	service := app.Service{SystemCollector: fakeSnapshotCollector{snapshot: testSnapshot()}}
	if err := run(context.Background(), nil, &output, &output, service); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "SETUP ENV SYSTEM STATUS") {
		t.Fatalf("status output = %q", output.String())
	}
}

func TestExplicitHelpStillDisplaysHelp(t *testing.T) {
	var output bytes.Buffer
	if err := run(context.Background(), []string{"--help"}, &output, &output, app.Service{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "setup-env <command>") {
		t.Fatalf("help output = %q", output.String())
	}
}

func TestStatusJSON(t *testing.T) {
	var output bytes.Buffer
	service := app.Service{SystemCollector: fakeSnapshotCollector{snapshot: testSnapshot()}}
	if err := run(context.Background(), []string{"status", "--json"}, &output, &output, service); err != nil {
		t.Fatal(err)
	}
	var snapshot system.Snapshot
	if err := json.Unmarshal(output.Bytes(), &snapshot); err != nil {
		t.Fatalf("status JSON is invalid: %v\n%s", err, output.String())
	}
	if snapshot.SchemaVersion != system.SnapshotSchemaVersion || strings.Contains(output.String(), "\x1b") {
		t.Fatalf("status JSON = %q", output.String())
	}
}

func TestStatusCollectorFailureReturnsError(t *testing.T) {
	var output bytes.Buffer
	service := app.Service{SystemCollector: fakeSnapshotCollector{err: errors.New("all collectors failed")}}
	err := run(context.Background(), []string{"status"}, &output, &output, service)
	if err == nil || !strings.Contains(err.Error(), "collect system status") {
		t.Fatalf("run() error = %v", err)
	}
}

func TestVersionCommand(t *testing.T) {
	var output bytes.Buffer
	if err := run(context.Background(), []string{"version"}, &output, &output, app.Service{}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(output.String(), "setup-env ") {
		t.Fatalf("version output = %q", output.String())
	}
}

func TestUnknownCommandIsActionable(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), []string{"unknown"}, &output, &output, app.Service{})
	if err == nil || !strings.Contains(err.Error(), "setup-env help") {
		t.Fatalf("run() error = %v", err)
	}
}

func TestModuleListHumanAndJSON(t *testing.T) {
	service := app.Service{CatalogSource: catalog.EmbeddedSource{}}
	var human bytes.Buffer
	if err := run(context.Background(), []string{"module", "list", "--status", "planned"}, &human, &human, service); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(human.String(), "workstation") || !strings.Contains(human.String(), "official") {
		t.Fatalf("human output = %q", human.String())
	}

	var machine bytes.Buffer
	if err := run(context.Background(), []string{"module", "list", "--category", "cloud", "--json"}, &machine, &machine, service); err != nil {
		t.Fatal(err)
	}
	var result struct {
		Modules []catalog.Entry `json:"modules"`
	}
	if err := json.Unmarshal(machine.Bytes(), &result); err != nil {
		t.Fatalf("JSON output is invalid: %v\n%s", err, machine.String())
	}
	if len(result.Modules) != 3 {
		t.Fatalf("cloud modules = %d, want 3", len(result.Modules))
	}
}

func TestModuleInfoWithoutLocalManifestIsHonest(t *testing.T) {
	home := t.TempDir()
	configDirectory := t.TempDir()
	service := app.DefaultService()
	service.PathResolver = paths.Resolver{UserHomeDir: func() (string, error) { return home, nil }}
	service.ConfigLocation = config.LocationResolver{UserConfigDir: func() (string, error) { return configDirectory, nil }}

	var output bytes.Buffer
	if err := run(context.Background(), []string{"module", "info", "workstation", "--json"}, &output, &output, service); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"manifest_available": false`) || !strings.Contains(output.String(), "Downloading is not implemented") {
		t.Fatalf("info output = %q", output.String())
	}
}

func TestModuleValidateValidManifest(t *testing.T) {
	service := app.DefaultService()
	path := filepath.Join("..", "..", "examples", "setup-env.yaml")
	var output bytes.Buffer
	if err := run(context.Background(), []string{"module", "validate", path, "--json"}, &output, &output, service); err != nil {
		t.Fatal(err)
	}
	var result app.ManifestValidationReport
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !result.Valid {
		t.Fatalf("validation result = %#v", result)
	}
}

func TestModuleValidateInvalidManifestReturnsErrorAndJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "setup-env.yaml")
	if err := os.WriteFile(path, []byte("schema_version: 1\nid: INVALID\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	service := app.DefaultService()
	var output bytes.Buffer
	err := run(context.Background(), []string{"module", "validate", path, "--json"}, &output, &output, service)
	if err == nil {
		t.Fatal("run() error = nil, want validation error")
	}
	var result app.ManifestValidationReport
	if jsonErr := json.Unmarshal(output.Bytes(), &result); jsonErr != nil {
		t.Fatalf("JSON output is invalid: %v\n%s", jsonErr, output.String())
	}
	if result.Valid || len(result.Problems) == 0 {
		t.Fatalf("validation result = %#v", result)
	}
}

func TestValidateCatalogCommand(t *testing.T) {
	service := app.Service{CatalogSource: catalog.EmbeddedSource{}}
	var output bytes.Buffer
	if err := run(context.Background(), []string{"module", "validate-catalog", "--json"}, &output, &output, service); err != nil {
		t.Fatal(err)
	}
	var result app.CatalogValidationReport
	if err := json.Unmarshal(output.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !result.Valid || result.ModuleCount != 10 {
		t.Fatalf("validation result = %#v", result)
	}
}

func testSnapshot() system.Snapshot {
	total := uint64(8 * 1024 * 1024 * 1024)
	used := uint64(4 * 1024 * 1024 * 1024)
	percent := 50.0
	uptime := uint64(3661)
	return system.Snapshot{
		SchemaVersion: system.SnapshotSchemaVersion,
		Timestamp:     time.Date(2026, 7, 24, 1, 30, 0, 0, time.FixedZone("SAST", 7200)),
		TimeZone:      system.TimeZone{Name: "SAST", UTCOffsetSeconds: 7200},
		Host:          system.Host{Hostname: "test-host", UptimeSeconds: &uptime},
		OperatingSystem: system.OperatingSystem{
			OS:           "linux",
			DisplayName:  "Ubuntu",
			Version:      "26.04",
			Architecture: "amd64",
		},
		User:        system.User{Username: "tester", Home: "/home/tester"},
		Memory:      system.Memory{TotalBytes: &total, UsedBytes: &used, AvailableBytes: &used, UtilizationPercent: &percent},
		Filesystems: []system.Filesystem{},
		Networks:    []system.NetworkInterface{},
		Warnings:    []system.Warning{},
		Diagnostics: system.DiagnosticSummary{Health: system.HealthHealthy},
	}
}

func TestUnsupportedInfoOption(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), []string{"info", "--yaml"}, &output, &output, app.Service{})
	if err == nil || !strings.Contains(err.Error(), "--json") {
		t.Fatalf("run() error = %v", err)
	}
}
