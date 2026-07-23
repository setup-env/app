package paths

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestDevelopmentRootDefaultsToHomeDev(t *testing.T) {
	resolver := Resolver{UserHomeDir: func() (string, error) {
		return filepath.Join(string(filepath.Separator), "users", "person"), nil
	}}
	root, err := resolver.DevelopmentRoot("")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(string(filepath.Separator), "users", "person", "dev")
	if root != want {
		t.Fatalf("DevelopmentRoot() = %q, want %q", root, want)
	}
}

func TestDevelopmentRootUsesAbsoluteOverride(t *testing.T) {
	resolver := Resolver{UserHomeDir: func() (string, error) {
		return "", errors.New("must not be called")
	}}
	root, err := resolver.DevelopmentRoot(filepath.Join(".", "custom-dev"))
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(root) {
		t.Fatalf("DevelopmentRoot() = %q, want absolute path", root)
	}
}

func TestHomeErrorIsActionable(t *testing.T) {
	resolver := Resolver{UserHomeDir: func() (string, error) {
		return "", errors.New("lookup failed")
	}}
	if _, err := resolver.Home(); err == nil {
		t.Fatal("Home() error = nil, want error")
	}
}
