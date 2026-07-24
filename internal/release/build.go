package release

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const ChecksumsName = "checksums.txt"

type Config struct {
	Root       string
	Output     string
	Version    string
	Commit     string
	BuildDate  time.Time
	Dirty      bool
	Targets    []Target
	CleanOwned bool
}

type Artifact struct {
	Target   Target
	Name     string
	Path     string
	Checksum string
}

func Build(ctx context.Context, config Config) ([]Artifact, error) {
	var err error
	config.Version, err = NormalizeVersion(config.Version)
	if err != nil {
		return nil, err
	}
	if config.Root == "" {
		config.Root = "."
	}
	config.Root, err = filepath.Abs(config.Root)
	if err != nil {
		return nil, fmt.Errorf("resolve repository root: %w", err)
	}
	if config.Output == "" {
		config.Output = filepath.Join(config.Root, "dist")
	} else if !filepath.IsAbs(config.Output) {
		config.Output = filepath.Join(config.Root, config.Output)
	}
	if len(config.Targets) == 0 {
		config.Targets = append([]Target(nil), SupportedTargets...)
	}
	if config.Commit == "" {
		return nil, fmt.Errorf("release commit is required")
	}
	if config.BuildDate.IsZero() {
		return nil, fmt.Errorf("release build date is required")
	}
	if err := prepareOutput(config.Output, config.CleanOwned); err != nil {
		return nil, err
	}
	license, err := os.ReadFile(filepath.Join(config.Root, "LICENSE"))
	if err != nil {
		return nil, fmt.Errorf("read LICENSE: %w", err)
	}
	releaseReadme, err := os.ReadFile(filepath.Join(config.Root, "release", "README.md"))
	if err != nil {
		return nil, fmt.Errorf("read release README: %w", err)
	}
	notices, err := os.ReadFile(filepath.Join(config.Root, "THIRD_PARTY_NOTICES.md"))
	if err != nil {
		return nil, fmt.Errorf("read third-party notices: %w", err)
	}

	buildDir, err := os.MkdirTemp("", "setup-env-release-build-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(buildDir)

	artifacts := make([]Artifact, 0, len(config.Targets))
	for _, target := range config.Targets {
		binaryPath := filepath.Join(buildDir, target.OS+"-"+target.Arch, target.BinaryName())
		if err := os.MkdirAll(filepath.Dir(binaryPath), 0o755); err != nil {
			return nil, err
		}
		if err := buildBinary(ctx, config, target, binaryPath); err != nil {
			return nil, err
		}
		binary, err := os.ReadFile(binaryPath)
		if err != nil {
			return nil, err
		}
		name := target.ArchiveName(config.Version)
		path := filepath.Join(config.Output, name)
		files := []archiveFile{
			{Name: target.BinaryName(), Mode: 0o755, Data: binary},
			{Name: "LICENSE", Mode: 0o644, Data: license},
			{Name: "README.md", Mode: 0o644, Data: releaseReadme},
			{Name: "THIRD_PARTY_NOTICES.md", Mode: 0o644, Data: notices},
		}
		if err := writeArchive(path, target, files, config.BuildDate); err != nil {
			return nil, fmt.Errorf("archive %s: %w", target, err)
		}
		checksum, err := checksumFile(path)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, Artifact{
			Target:   target,
			Name:     name,
			Path:     path,
			Checksum: checksum,
		})
	}
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Name < artifacts[j].Name })
	var checksums strings.Builder
	for _, artifact := range artifacts {
		fmt.Fprintf(&checksums, "%s  %s\n", artifact.Checksum, artifact.Name)
	}
	if err := writeFileAtomic(filepath.Join(config.Output, ChecksumsName), []byte(checksums.String()), 0o644); err != nil {
		return nil, fmt.Errorf("write checksums: %w", err)
	}
	if err := Verify(config.Output, config.Version, config.Targets); err != nil {
		return nil, fmt.Errorf("verify release output: %w", err)
	}
	return artifacts, nil
}

func buildBinary(ctx context.Context, config Config, target Target, output string) error {
	date := config.BuildDate.UTC().Format(time.RFC3339)
	ldflags := strings.Join([]string{
		"-s",
		"-w",
		"-buildid=",
		"-X", "github.com/setup-env/app/internal/version.Version=" + config.Version,
		"-X", "github.com/setup-env/app/internal/version.Commit=" + config.Commit,
		"-X", "github.com/setup-env/app/internal/version.Date=" + date,
		"-X", "github.com/setup-env/app/internal/version.Dirty=" + strconv.FormatBool(config.Dirty),
	}, " ")
	command := exec.CommandContext(
		ctx,
		"go",
		"build",
		"-mod=readonly",
		"-trimpath",
		"-buildvcs=false",
		"-ldflags", ldflags,
		"-o", output,
		"./cmd/setup-env",
	)
	command.Dir = config.Root
	command.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS="+target.OS,
		"GOARCH="+target.Arch,
	)
	outputText, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build %s: %w\n%s", target, err, strings.TrimSpace(string(outputText)))
	}
	return nil
}

func prepareOutput(path string, cleanOwned bool) error {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	volume := filepath.VolumeName(absolute)
	root := string(filepath.Separator)
	if volume != "" {
		root = volume + string(filepath.Separator)
	}
	if filepath.Clean(absolute) == filepath.Clean(root) {
		return fmt.Errorf("refusing to use filesystem root as release output")
	}
	if err := os.MkdirAll(absolute, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(absolute)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		owned := !entry.IsDir() &&
			(strings.HasPrefix(entry.Name(), "setup-env_") || entry.Name() == ChecksumsName)
		if cleanOwned && owned {
			if err := os.Remove(filepath.Join(absolute, entry.Name())); err != nil {
				return err
			}
			continue
		}
		return fmt.Errorf("release output %s is not empty; remove %s or use -clean-owned", absolute, entry.Name())
	}
	return nil
}

func checksumFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func ResolveGitMetadata(ctx context.Context, root string) (commit string, date time.Time, dirty bool, err error) {
	run := func(arguments ...string) (string, error) {
		command := exec.CommandContext(ctx, "git", arguments...)
		command.Dir = root
		output, commandErr := command.CombinedOutput()
		if commandErr != nil {
			return "", fmt.Errorf("git %s: %w: %s", strings.Join(arguments, " "), commandErr, strings.TrimSpace(string(output)))
		}
		return strings.TrimSpace(string(output)), nil
	}
	commit, err = run("rev-parse", "HEAD")
	if err != nil {
		return "", time.Time{}, false, err
	}
	dateText, err := run("show", "-s", "--format=%cI", "HEAD")
	if err != nil {
		return "", time.Time{}, false, err
	}
	date, err = time.Parse(time.RFC3339, dateText)
	if err != nil {
		return "", time.Time{}, false, fmt.Errorf("parse commit date %q: %w", dateText, err)
	}
	status, err := run("status", "--porcelain")
	if err != nil {
		return "", time.Time{}, false, err
	}
	return commit, date, status != "", nil
}

func ValidateReleaseVersion(version string) error {
	_, err := NormalizeVersion(version)
	return err
}

var errChecksumMismatch = errors.New("checksum mismatch")
