package app

import (
	"context"
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/setup-env/app/internal/config"
	"github.com/setup-env/app/internal/paths"
	"github.com/setup-env/app/internal/platform"
	"github.com/setup-env/app/internal/system"
)

type statusTestRunner struct{}

func (statusTestRunner) LookPath(name string) (string, error) {
	if name == "git" || name == "gh" {
		return name, nil
	}
	return "", errors.New("not found")
}

func (statusTestRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	if name == "git" {
		return []byte("git version 2.50.0\n"), nil
	}
	if len(args) > 0 && args[0] == "auth" {
		return []byte("authenticated"), nil
	}
	return []byte("gh version 2.80.0\n"), nil
}

func TestDevelopmentCollectorReusesInfoAndDoctorModels(t *testing.T) {
	home := t.TempDir()
	root := filepath.Join(home, "dev")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	service := Service{
		PlatformDetector: platform.Detector{
			GOOS:   "windows",
			GOARCH: "amd64",
			CurrentUser: func() (*user.User, error) {
				return &user.User{Username: "tester"}, nil
			},
			UserHomeDir: func() (string, error) { return home, nil },
			LookupEnv:   func(string) (string, bool) { return "", false },
		},
		PathResolver:   paths.Resolver{UserHomeDir: func() (string, error) { return home, nil }},
		ConfigLocation: config.LocationResolver{UserConfigDir: func() (string, error) { return t.TempDir(), nil }},
		Getwd:          func() (string, error) { return root, nil },
		Commands:       statusTestRunner{},
	}

	snapshot := system.Snapshot{}
	if err := (DevelopmentCollector{Service: service}).Collect(context.Background(), &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.Development.Root.Path != root ||
		!snapshot.Development.Root.Exists ||
		!snapshot.Development.Root.Writable {
		t.Fatalf("development root = %#v", snapshot.Development.Root)
	}
	if snapshot.Development.Directory.DevelopmentRoot != root {
		t.Fatalf("directory = %#v", snapshot.Development.Directory)
	}
	if !snapshot.Development.Git.Available ||
		!snapshot.Development.GitHubCLI.Available ||
		snapshot.Development.GitHubAuthenticated == nil ||
		!*snapshot.Development.GitHubAuthenticated {
		t.Fatalf("tool status = %#v", snapshot.Development)
	}
	if len(snapshot.Diagnostics.Checks) != 4 {
		t.Fatalf("diagnostic checks = %#v", snapshot.Diagnostics.Checks)
	}
}
