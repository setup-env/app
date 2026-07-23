package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validManifest(t *testing.T) Manifest {
	t.Helper()
	path := filepath.Join("..", "..", "examples", "setup-env.yaml")
	value, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestReferenceManifestParsesAndValidates(t *testing.T) {
	value := validManifest(t)
	if err := ValidateSchema(value); err != nil {
		t.Fatal(err)
	}
	if err := ValidateSemantics(value); err != nil {
		t.Fatal(err)
	}
}

func TestParseRejectsUnknownTrustField(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "examples", "setup-env.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, []byte("\ntrust: official\n")...)
	if _, err := Parse(data); err == nil {
		t.Fatal("Parse() error = nil, want unknown trust field error")
	}
}

func TestMissingRequiredFieldsAreReportedTogether(t *testing.T) {
	err := Validate(Manifest{SchemaVersion: SchemaVersion})
	validationError, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("Validate() error = %T, want *ValidationError", err)
	}
	if len(validationError.Problems) < 8 {
		t.Fatalf("Validate() problems = %v, want multiple required-field errors", validationError.Problems)
	}
}

func TestInvalidModuleID(t *testing.T) {
	value := validManifest(t)
	value.ID = "Workstation Setup"
	assertProblemContains(t, ValidateSemantics(value), "id must be lowercase")
}

func TestDuplicateWorkflowID(t *testing.T) {
	value := validManifest(t)
	value.Workflows = append(value.Workflows, value.Workflows[0])
	assertProblemContains(t, ValidateSemantics(value), "workflow id install is duplicated")
}

func TestInvalidPlatformAndArchitecture(t *testing.T) {
	value := validManifest(t)
	value.Platforms.OperatingSystems = []string{"plan9"}
	value.Platforms.Architectures = []string{"quantum"}
	err := ValidateSemantics(value)
	assertProblemContains(t, err, "unsupported value plan9")
	assertProblemContains(t, err, "unsupported value quantum")
}

func TestMalformedMinimumAppVersion(t *testing.T) {
	value := validManifest(t)
	value.MinimumAppVersion = "latest"
	assertProblemContains(t, ValidateSemantics(value), "minimum_app_version")
}

func TestDuplicateCategory(t *testing.T) {
	value := validManifest(t)
	value.Categories = []string{"workstation", "workstation"}
	assertProblemContains(t, ValidateSemantics(value), "duplicate value workstation")
}

func assertProblemContains(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), expected) {
		t.Fatalf("error = %v, want substring %q", err, expected)
	}
}
