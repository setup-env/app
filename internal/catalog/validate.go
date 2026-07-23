package catalog

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/setup-env/app/internal/compatibility"
)

type ValidationError struct {
	Problems []string `json:"problems"`
}

func (e *ValidationError) Error() string {
	return "catalog validation failed: " + strings.Join(e.Problems, "; ")
}

var identifierPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
var repositoryPattern = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,37}[A-Za-z0-9])?/[A-Za-z0-9._-]+$`)

var validTrust = map[Trust]struct{}{
	TrustOfficial:  {},
	TrustVerified:  {},
	TrustCommunity: {},
}

var validStatus = map[Status]struct{}{
	StatusActive:       {},
	StatusPlanned:      {},
	StatusExperimental: {},
	StatusDeprecated:   {},
	StatusUnavailable:  {},
}

var validVersionPolicies = map[string]struct{}{
	"":               {},
	"github-release": {},
	"git-tag":        {},
}

func Validate(value Catalog) error {
	var problems []string
	if value.SchemaVersion != SchemaVersion {
		problems = append(problems, fmt.Sprintf("schema_version must be %d, got %d", SchemaVersion, value.SchemaVersion))
	}
	if len(value.Modules) == 0 {
		problems = append(problems, "modules must contain at least one entry")
	}

	ids := make(map[string]struct{}, len(value.Modules))
	repositories := make(map[string]struct{}, len(value.Modules))
	previousID := ""
	for index, entry := range value.Modules {
		prefix := fmt.Sprintf("modules[%d]", index)
		if !identifierPattern.MatchString(entry.ID) {
			problems = append(problems, prefix+".id must be a lowercase CLI-safe identifier")
		}
		if _, exists := ids[entry.ID]; exists {
			problems = append(problems, "module id "+entry.ID+" is duplicated")
		}
		ids[entry.ID] = struct{}{}
		if previousID != "" && previousID >= entry.ID {
			problems = append(problems, "modules must be sorted by id")
		}
		previousID = entry.ID

		if strings.TrimSpace(entry.Name) == "" {
			problems = append(problems, prefix+".name must not be empty")
		}
		if strings.TrimSpace(entry.Description) == "" {
			problems = append(problems, prefix+".description must not be empty")
		}
		if !repositoryPattern.MatchString(entry.Repository) {
			problems = append(problems, prefix+".repository must use owner/name form")
		}
		normalizedRepository := strings.ToLower(entry.Repository)
		if _, exists := repositories[normalizedRepository]; exists {
			problems = append(problems, "repository "+entry.Repository+" is duplicated")
		}
		repositories[normalizedRepository] = struct{}{}
		owner, repository := entry.RepositoryParts()
		if entry.Trust == TrustOfficial && owner != "setup-env" {
			problems = append(problems, "official trust is restricted to setup-env-owned repositories")
		}
		if entry.ID == "app" || repository == "app" {
			problems = append(problems, "setup-env/app must not be listed as a module")
		}
		if entry.ID == "awesome-setup-env" || repository == "awesome-setup-env" {
			problems = append(problems, "awesome-setup-env must not be listed as a module")
		}

		if _, ok := validTrust[entry.Trust]; !ok {
			problems = append(problems, prefix+".trust is not recognized")
		}
		if _, ok := validStatus[entry.Status]; !ok {
			problems = append(problems, prefix+".status is not recognized")
		}
		if !validManifestPath(entry.Manifest) {
			problems = append(problems, prefix+".manifest must be a relative path ending in setup-env.yaml")
		}
		validateIdentifiers(&problems, prefix+".categories", entry.Categories, true)
		validateIdentifiers(&problems, prefix+".tags", entry.Tags, false)
		if !sort.StringsAreSorted(entry.Categories) {
			problems = append(problems, prefix+".categories must be sorted")
		}
		if !sort.StringsAreSorted(entry.Tags) {
			problems = append(problems, prefix+".tags must be sorted")
		}
		if entry.PinnedVersion != "" {
			if _, err := compatibility.Parse(entry.PinnedVersion); err != nil {
				problems = append(problems, prefix+".pinned_version must be a semantic version")
			}
			if entry.VersionPolicy != "" {
				problems = append(problems, prefix+" must not set both pinned_version and version_policy")
			}
		}
		if _, ok := validVersionPolicies[entry.VersionPolicy]; !ok {
			problems = append(problems, prefix+".version_policy is not recognized")
		}
	}
	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}
	return nil
}

func Filter(entries []Entry, trust Trust, status Status, category string) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if trust != "" && entry.Trust != trust {
			continue
		}
		if status != "" && entry.Status != status {
			continue
		}
		if category != "" && !contains(entry.Categories, category) {
			continue
		}
		result = append(result, entry)
	}
	return result
}

func Find(entries []Entry, id string) (Entry, bool) {
	index := sort.Search(len(entries), func(index int) bool {
		return entries[index].ID >= id
	})
	if index < len(entries) && entries[index].ID == id {
		return entries[index], true
	}
	return Entry{}, false
}

func validManifestPath(value string) bool {
	if value == "" || strings.Contains(value, "\\") || path.IsAbs(value) {
		return false
	}
	clean := path.Clean(value)
	return clean == value && !strings.HasPrefix(clean, "../") && path.Base(clean) == "setup-env.yaml"
}

func validateIdentifiers(problems *[]string, field string, values []string, required bool) {
	if required && len(values) == 0 {
		*problems = append(*problems, field+" must contain at least one value")
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !identifierPattern.MatchString(value) {
			*problems = append(*problems, field+" contains invalid identifier "+value)
		}
		if _, exists := seen[value]; exists {
			*problems = append(*problems, field+" contains duplicate value "+value)
		}
		seen[value] = struct{}{}
	}
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
