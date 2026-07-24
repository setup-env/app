package release

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Target struct {
	OS   string
	Arch string
}

var SupportedTargets = []Target{
	{OS: "windows", Arch: "amd64"},
	{OS: "windows", Arch: "arm64"},
	{OS: "darwin", Arch: "amd64"},
	{OS: "darwin", Arch: "arm64"},
	{OS: "linux", Arch: "amd64"},
	{OS: "linux", Arch: "arm64"},
}

var semanticVersion = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?$`)

func NormalizeVersion(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("release version is required")
	}
	if !strings.HasPrefix(value, "v") {
		value = "v" + value
	}
	if !semanticVersion.MatchString(value) {
		return "", fmt.Errorf("version %q is not semantic vMAJOR.MINOR.PATCH", value)
	}
	return value, nil
}

func ParseTargets(value string) ([]Target, error) {
	if strings.TrimSpace(value) == "" || value == "all" {
		return append([]Target(nil), SupportedTargets...), nil
	}
	allowed := make(map[string]Target, len(SupportedTargets))
	for _, target := range SupportedTargets {
		allowed[target.String()] = target
	}
	seen := map[string]bool{}
	var result []Target
	for _, item := range strings.Split(value, ",") {
		key := strings.TrimSpace(strings.ToLower(item))
		target, ok := allowed[key]
		if !ok {
			return nil, fmt.Errorf("unsupported release target %q", item)
		}
		if !seen[key] {
			seen[key] = true
			result = append(result, target)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].String() < result[j].String()
	})
	return result, nil
}

func (t Target) String() string {
	return t.OS + "/" + t.Arch
}

func (t Target) BinaryName() string {
	if t.OS == "windows" {
		return "setup-env.exe"
	}
	return "setup-env"
}

func (t Target) ArchiveName(version string) string {
	version = strings.TrimPrefix(version, "v")
	extension := ".tar.gz"
	if t.OS == "windows" {
		extension = ".zip"
	}
	return fmt.Sprintf("setup-env_%s_%s_%s%s", version, t.OS, t.Arch, extension)
}
