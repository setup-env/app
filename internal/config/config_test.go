package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	settings := Default()
	if settings.SchemaVersion != SchemaVersion {
		t.Fatalf("SchemaVersion = %d, want %d", settings.SchemaVersion, SchemaVersion)
	}
	if settings.CloneProtocol != CloneProtocolHTTPS {
		t.Fatalf("CloneProtocol = %q, want %q", settings.CloneProtocol, CloneProtocolHTTPS)
	}
	if !settings.Settings.CheckUpdates {
		t.Fatal("CheckUpdates = false, want true")
	}
}

func TestConfigurationPathUsesOSConfigDirectory(t *testing.T) {
	base := t.TempDir()
	resolver := LocationResolver{UserConfigDir: func() (string, error) { return base, nil }}
	got, err := resolver.Path()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(base, "setup-env", "config.json")
	if got != want {
		t.Fatalf("Path() = %q, want %q", got, want)
	}
}

func TestMissingConfigurationIsOptional(t *testing.T) {
	got, loaded, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatal(err)
	}
	if loaded {
		t.Fatal("loaded = true, want false")
	}
	if got.SchemaVersion != SchemaVersion {
		t.Fatalf("SchemaVersion = %d", got.SchemaVersion)
	}
}

func TestLoadRejectsUnknownSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"schema_version":99,"clone_protocol":"https"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Load(path); err == nil {
		t.Fatal("Load() error = nil, want schema error")
	}
}
