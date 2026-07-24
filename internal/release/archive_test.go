package release

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestArchivesAreDeterministicAndMinimal(t *testing.T) {
	timestamp := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)
	files := []archiveFile{
		{Name: "setup-env", Mode: 0o755, Data: []byte("binary")},
		{Name: "LICENSE", Mode: 0o644, Data: []byte("license")},
		{Name: "README.md", Mode: 0o644, Data: []byte("readme")},
		{Name: "THIRD_PARTY_NOTICES.md", Mode: 0o644, Data: []byte("notices")},
	}
	for _, target := range []Target{{OS: "windows", Arch: "amd64"}, {OS: "linux", Arch: "amd64"}} {
		targetFiles := append([]archiveFile(nil), files...)
		targetFiles[0].Name = target.BinaryName()
		first := filepath.Join(t.TempDir(), target.ArchiveName("v0.1.0"))
		second := filepath.Join(t.TempDir(), target.ArchiveName("v0.1.0"))
		if err := writeArchive(first, target, targetFiles, timestamp); err != nil {
			t.Fatal(err)
		}
		if err := writeArchive(second, target, targetFiles, timestamp); err != nil {
			t.Fatal(err)
		}
		firstData, _ := os.ReadFile(first)
		secondData, _ := os.ReadFile(second)
		if !bytes.Equal(firstData, secondData) {
			t.Fatalf("%s archives are not reproducible", target)
		}
		if err := verifyArchive(first, target); err != nil {
			t.Fatalf("verifyArchive(%s): %v", target, err)
		}
	}
}

func TestVerifyRejectsChecksumMismatch(t *testing.T) {
	target := Target{OS: "windows", Arch: "amd64"}
	output := t.TempDir()
	path := filepath.Join(output, target.ArchiveName("v0.1.0"))
	files := []archiveFile{
		{Name: target.BinaryName(), Mode: 0o755, Data: []byte("binary")},
		{Name: "LICENSE", Mode: 0o644, Data: []byte("license")},
		{Name: "README.md", Mode: 0o644, Data: []byte("readme")},
		{Name: "THIRD_PARTY_NOTICES.md", Mode: 0o644, Data: []byte("notices")},
	}
	if err := writeArchive(path, target, files, time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, ChecksumsName), []byte(
		"0000000000000000000000000000000000000000000000000000000000000000  "+filepath.Base(path)+"\n",
	), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Verify(output, "v0.1.0", []Target{target}); err == nil {
		t.Fatal("Verify() error = nil, want checksum mismatch")
	}
}

func TestZipEntryNames(t *testing.T) {
	target := Target{OS: "windows", Arch: "arm64"}
	path := filepath.Join(t.TempDir(), target.ArchiveName("v0.1.0"))
	files := []archiveFile{
		{Name: "README.md", Mode: 0o644, Data: []byte("readme")},
		{Name: target.BinaryName(), Mode: 0o755, Data: []byte("binary")},
	}
	if err := writeArchive(path, target, files, time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	reader, err := zip.OpenReader(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	var names []string
	for _, file := range reader.File {
		names = append(names, file.Name)
	}
	sort.Strings(names)
	if !reflect.DeepEqual(names, []string{"README.md", "setup-env.exe"}) {
		t.Fatalf("zip names = %v", names)
	}
}
