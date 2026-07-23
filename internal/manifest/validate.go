package manifest

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/setup-env/app/internal/compatibility"
)

type ValidationError struct {
	Problems []string `json:"problems"`
}

func (e *ValidationError) Error() string {
	return "manifest validation failed: " + strings.Join(e.Problems, "; ")
}

var (
	identifierPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	ownerPattern      = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,37}[A-Za-z0-9])?$`)
	repositoryPattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
)

var validOperatingSystems = setOf("windows", "darwin", "linux")
var validArchitectures = setOf("386", "amd64", "arm", "arm64", "loong64", "mips", "mipsle", "mips64", "mips64le", "ppc64", "ppc64le", "riscv64", "s390x", "wasm")
var validVersionSources = setOf("github-release", "git-tag", "manifest")

func ValidateSchema(value Manifest) error {
	if value.SchemaVersion != SchemaVersion {
		return &ValidationError{Problems: []string{
			fmt.Sprintf("schema_version must be %d, got %d", SchemaVersion, value.SchemaVersion),
		}}
	}
	return nil
}

func ValidateSemantics(value Manifest) error {
	var problems []string
	requireIdentifier(&problems, "id", value.ID)
	requireText(&problems, "name", value.Name)
	requireText(&problems, "description", value.Description)
	requireText(&problems, "publisher", value.Publisher)
	requireText(&problems, "license", value.License)

	if !ownerPattern.MatchString(value.Repository.Owner) {
		problems = append(problems, "repository.owner must be a valid GitHub owner")
	}
	if !repositoryPattern.MatchString(value.Repository.Name) {
		problems = append(problems, "repository.name must be a valid repository name")
	}
	validateURL(&problems, "repository.issues_url", value.Repository.IssuesURL)
	validateURL(&problems, "homepage", value.Homepage)
	validateURL(&problems, "documentation", value.Documentation)

	if _, ok := validVersionSources[value.Version.Source]; !ok {
		problems = append(problems, "version.source must be one of github-release, git-tag, or manifest")
	}
	if value.Version.Source == "manifest" {
		if _, err := compatibility.Parse(value.Version.Value); err != nil {
			problems = append(problems, "version.value must be a semantic version when version.source is manifest")
		}
	} else if value.Version.Value != "" {
		problems = append(problems, "version.value must be omitted unless version.source is manifest")
	}
	if _, err := compatibility.Parse(value.MinimumAppVersion); err != nil {
		problems = append(problems, "minimum_app_version must be a semantic version in MAJOR.MINOR.PATCH form")
	}

	validateValues(&problems, "platforms.operating_systems", value.Platforms.OperatingSystems, validOperatingSystems, true)
	validateValues(&problems, "platforms.architectures", value.Platforms.Architectures, validArchitectures, true)
	validateIdentifiers(&problems, "categories", value.Categories, true)
	validateIdentifiers(&problems, "tags", value.Tags, false)
	validateIdentifiers(&problems, "security.secret_inputs", value.Security.SecretInputs, false)

	if len(value.Workflows) == 0 {
		problems = append(problems, "workflows must contain at least one declaration")
	}
	workflowIDs := make(map[string]struct{}, len(value.Workflows))
	for index, workflow := range value.Workflows {
		prefix := fmt.Sprintf("workflows[%d]", index)
		if !identifierPattern.MatchString(workflow.ID) {
			problems = append(problems, prefix+".id must be a lowercase CLI-safe identifier")
		} else if _, exists := workflowIDs[workflow.ID]; exists {
			problems = append(problems, "workflow id "+workflow.ID+" is duplicated")
		} else {
			workflowIDs[workflow.ID] = struct{}{}
		}
		requireText(&problems, prefix+".name", workflow.Name)
		requireText(&problems, prefix+".description", workflow.Description)
		if !validRelativeYAMLPath(workflow.Entrypoint) || !strings.HasPrefix(workflow.Entrypoint, "workflows/") {
			problems = append(problems, prefix+".entrypoint must be a relative YAML path under workflows/")
		}
	}

	if value.Replacement != "" {
		if !value.Deprecated {
			problems = append(problems, "replacement may only be set when deprecated is true")
		}
		if !identifierPattern.MatchString(value.Replacement) {
			problems = append(problems, "replacement must be a lowercase CLI-safe module id")
		}
		if value.Replacement == value.ID {
			problems = append(problems, "replacement must not refer to the same module")
		}
	}
	if value.DeprecationNotice != "" && !value.Deprecated {
		problems = append(problems, "deprecation_notice may only be set when deprecated is true")
	}

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}
	return nil
}

func Validate(value Manifest) error {
	var problems []string
	if err := ValidateSchema(value); err != nil {
		problems = append(problems, err.(*ValidationError).Problems...)
	}
	if err := ValidateSemantics(value); err != nil {
		problems = append(problems, err.(*ValidationError).Problems...)
	}
	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}
	return nil
}

func requireIdentifier(problems *[]string, field, value string) {
	if !identifierPattern.MatchString(value) {
		*problems = append(*problems, field+" must be lowercase and contain only letters, numbers, and single hyphens")
	}
}

func requireText(problems *[]string, field, value string) {
	if strings.TrimSpace(value) == "" {
		*problems = append(*problems, field+" must not be empty")
	}
}

func validateURL(problems *[]string, field, value string) {
	if value == "" {
		return
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Host == "" {
		*problems = append(*problems, field+" must be an absolute HTTP or HTTPS URL")
	}
}

func validateIdentifiers(problems *[]string, field string, values []string, required bool) {
	if required && len(values) == 0 {
		*problems = append(*problems, field+" must contain at least one value")
		return
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

func validateValues(problems *[]string, field string, values []string, allowed map[string]struct{}, required bool) {
	validateIdentifiers(problems, field, values, required)
	for _, value := range values {
		if _, ok := allowed[value]; !ok {
			*problems = append(*problems, field+" contains unsupported value "+value)
		}
	}
}

func validRelativeYAMLPath(value string) bool {
	if value == "" || strings.Contains(value, "\\") || path.IsAbs(value) {
		return false
	}
	clean := path.Clean(value)
	return clean == value && clean != "." && !strings.HasPrefix(clean, "../") && (strings.HasSuffix(clean, ".yaml") || strings.HasSuffix(clean, ".yml"))
}

func setOf(values ...string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
