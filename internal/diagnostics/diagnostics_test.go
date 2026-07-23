package diagnostics

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type fakeRunner struct {
	paths   map[string]string
	outputs map[string][]byte
	errors  map[string]error
}

func (f fakeRunner) LookPath(name string) (string, error) {
	if path, ok := f.paths[name]; ok {
		return path, nil
	}
	return "", errors.New("not found")
}

func (f fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	key := filepath.Base(name)
	if len(args) >= 2 && args[0] == "auth" {
		key = "gh-auth"
	}
	return f.outputs[key], f.errors[key]
}

func TestRunReportsToolReadinessWithoutSecrets(t *testing.T) {
	root := t.TempDir()
	runner := fakeRunner{
		paths: map[string]string{"git": "git", "gh": "gh"},
		outputs: map[string][]byte{
			"git":     []byte("git version 2.50.0\n"),
			"gh":      []byte("gh version 2.80.0\n"),
			"gh-auth": []byte("authenticated"),
		},
		errors: map[string]error{},
	}
	report := Run(context.Background(), root, runner)
	if !report.Ready {
		t.Fatalf("Run() = %#v, want ready", report)
	}
	for _, check := range report.Checks {
		if check.Name == "github-auth" && check.Message != "GitHub CLI authentication appears configured" {
			t.Fatalf("authentication message exposed command output: %q", check.Message)
		}
	}
}

func TestDevelopmentRootFailure(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	check := checkDevelopmentRoot(file)
	if check.Status != StatusFail {
		t.Fatalf("status = %q, want %q", check.Status, StatusFail)
	}
}
