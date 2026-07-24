package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestWorkflowYAMLAndPermissionBoundaries(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, name := range []string{"ci.yml", "release.yml"} {
		path := filepath.Join(root, ".github", "workflows", name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var document map[string]any
		if err := yaml.Unmarshal(data, &document); err != nil {
			t.Fatalf("%s is invalid YAML: %v", name, err)
		}
		text := string(data)
		if !strings.Contains(text, "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1") ||
			!strings.Contains(text, "actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e") {
			t.Fatalf("%s does not pin release-critical official actions", name)
		}
	}

	ci, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ci), "permissions:\n  contents: read") {
		t.Fatal("pull-request CI does not have explicit read-only contents permission")
	}

	releaseWorkflow, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(releaseWorkflow)
	if !strings.Contains(text, "permissions:\n  contents: write") ||
		!strings.Contains(text, "tags:\n      - \"v*\"") ||
		strings.Contains(text, "pull_request:") {
		t.Fatal("release workflow trigger or permission boundary is unsafe")
	}
}
