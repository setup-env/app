package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/setup-env/app/internal/app"
	"github.com/setup-env/app/internal/catalog"
	"github.com/setup-env/app/internal/config"
	"github.com/setup-env/app/internal/paths"
)

func TestNoArgumentsDisplaysHelp(t *testing.T) {
	var output bytes.Buffer
	if err := run(context.Background(), nil, &output, &output, app.Service{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "setup-env <command>") {
		t.Fatalf("help output = %q", output.String())
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

func TestUnsupportedInfoOption(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), []string{"info", "--yaml"}, &output, &output, app.Service{})
	if err == nil || !strings.Contains(err.Error(), "--json") {
		t.Fatalf("run() error = %v", err)
	}
}
