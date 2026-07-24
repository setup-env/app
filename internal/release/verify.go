package release

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Verify(output, version string, targets []Target) error {
	version, err := NormalizeVersion(version)
	if err != nil {
		return err
	}
	checksums, err := readChecksums(filepath.Join(output, ChecksumsName))
	if err != nil {
		return err
	}
	expectedNames := make([]string, 0, len(targets))
	for _, target := range targets {
		name := target.ArchiveName(version)
		expectedNames = append(expectedNames, name)
		expectedChecksum, ok := checksums[name]
		if !ok {
			return fmt.Errorf("checksum missing for %s", name)
		}
		path := filepath.Join(output, name)
		actualChecksum, err := checksumFile(path)
		if err != nil {
			return err
		}
		if actualChecksum != expectedChecksum {
			return fmt.Errorf("%w for %s: got %s, want %s", errChecksumMismatch, name, actualChecksum, expectedChecksum)
		}
		if err := verifyArchive(path, target); err != nil {
			return fmt.Errorf("verify %s: %w", name, err)
		}
	}
	if len(checksums) != len(expectedNames) {
		return fmt.Errorf("checksums cover %d files, expected %d", len(checksums), len(expectedNames))
	}
	return nil
}

func readChecksums(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open checksums: %w", err)
	}
	defer file.Close()
	result := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 || len(fields[0]) != 64 {
			return nil, fmt.Errorf("invalid checksum line %q", scanner.Text())
		}
		if _, exists := result[fields[1]]; exists {
			return nil, fmt.Errorf("duplicate checksum for %s", fields[1])
		}
		result[fields[1]] = strings.ToLower(fields[0])
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func verifyArchive(path string, target Target) error {
	var names []string
	if target.OS == "windows" {
		reader, err := zip.OpenReader(path)
		if err != nil {
			return err
		}
		defer reader.Close()
		for _, file := range reader.File {
			names = append(names, file.Name)
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		entries, err := readTarGzip(file)
		if err != nil {
			return err
		}
		for name := range entries {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	expected := []string{"LICENSE", "README.md", "THIRD_PARTY_NOTICES.md", target.BinaryName()}
	sort.Strings(expected)
	if strings.Join(names, "\n") != strings.Join(expected, "\n") {
		return fmt.Errorf("archive entries = %v, want %v", names, expected)
	}
	for _, name := range names {
		if filepath.IsAbs(name) || strings.Contains(name, "..") || strings.Contains(name, `\`) {
			return fmt.Errorf("unsafe archive path %q", name)
		}
	}
	return nil
}
