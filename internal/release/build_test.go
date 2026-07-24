package release

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadChecksums(t *testing.T) {
	path := filepath.Join(t.TempDir(), ChecksumsName)
	line := strings.Repeat("a", 64) + "  setup-env_0.1.0_linux_amd64.tar.gz\n"
	if err := os.WriteFile(path, []byte(line), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := readChecksums(path)
	if err != nil {
		t.Fatal(err)
	}
	if got["setup-env_0.1.0_linux_amd64.tar.gz"] != strings.Repeat("a", 64) {
		t.Fatalf("checksums = %v", got)
	}
}

func TestPrepareOutputRemovesOnlyOwnedFiles(t *testing.T) {
	output := t.TempDir()
	owned := filepath.Join(output, "setup-env_0.1.0_linux_amd64.tar.gz")
	if err := os.WriteFile(owned, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := prepareOutput(output, true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(owned); !os.IsNotExist(err) {
		t.Fatalf("owned artifact still exists: %v", err)
	}

	unrelated := filepath.Join(output, "keep.txt")
	if err := os.WriteFile(unrelated, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := prepareOutput(output, true); err == nil {
		t.Fatal("prepareOutput() removed or accepted unrelated file")
	}
	if _, err := os.Stat(unrelated); err != nil {
		t.Fatalf("unrelated file was changed: %v", err)
	}
}
