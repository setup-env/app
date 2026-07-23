package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/setup-env/app/internal/catalog"
	"github.com/setup-env/app/internal/config"
	"github.com/setup-env/app/internal/manifest"
	"github.com/setup-env/app/internal/paths"
)

func TestModuleInfoLoadsMatchingLocalManifest(t *testing.T) {
	home := t.TempDir()
	repository := filepath.Join(home, "dev", "setup-env", "workstation")
	if err := os.MkdirAll(repository, 0o700); err != nil {
		t.Fatal(err)
	}
	example, err := os.ReadFile(filepath.Join("..", "..", "examples", "setup-env.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "setup-env.yaml"), example, 0o600); err != nil {
		t.Fatal(err)
	}

	service := DefaultService()
	service.PathResolver = paths.Resolver{UserHomeDir: func() (string, error) { return home, nil }}
	service.ConfigLocation = config.LocationResolver{UserConfigDir: func() (string, error) { return t.TempDir(), nil }}
	info, err := service.ModuleInfo(context.Background(), "workstation")
	if err != nil {
		t.Fatal(err)
	}
	if !info.ManifestAvailable || info.Manifest == nil || info.Manifest.ID != "workstation" {
		t.Fatalf("ModuleInfo() = %#v", info)
	}
	if !strings.Contains(info.Message, "workflow execution") {
		t.Fatalf("ModuleInfo() message = %q", info.Message)
	}
}

func TestCatalogManifestMismatchIsRejected(t *testing.T) {
	entry := catalog.Entry{ID: "terraform", Repository: "setup-env/terraform"}
	moduleManifest := validAppTestManifest(t)
	if err := matchCatalogEntry(entry, moduleManifest); err == nil {
		t.Fatal("matchCatalogEntry() error = nil, want mismatch")
	}
}

func TestValidateManifestAcceptsRepositoryDirectory(t *testing.T) {
	directory := t.TempDir()
	example, err := os.ReadFile(filepath.Join("..", "..", "examples", "setup-env.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "setup-env.yaml"), example, 0o600); err != nil {
		t.Fatal(err)
	}
	report := DefaultService().ValidateManifest(directory)
	if !report.Valid {
		t.Fatalf("ValidateManifest() = %#v", report)
	}
}

func validAppTestManifest(t *testing.T) manifest.Manifest {
	t.Helper()
	value, err := manifest.ParseFile(filepath.Join("..", "..", "examples", "setup-env.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return value
}
