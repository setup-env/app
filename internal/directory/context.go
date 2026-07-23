package directory

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	gitinspect "github.com/setup-env/app/internal/git"
)

type Type string

const (
	TypeDevelopmentRoot Type = "development-root"
	TypeOrganization    Type = "organization"
	TypeRepository      Type = "repository"
	TypeOther           Type = "other"
)

type Context struct {
	DevelopmentRoot    string                 `json:"development_root"`
	CurrentDirectory   string                 `json:"current_directory"`
	Type               Type                   `json:"type"`
	Organization       string                 `json:"organization,omitempty"`
	Repository         string                 `json:"repository,omitempty"`
	IsGitRepository    bool                   `json:"is_git_repository"`
	Git                *gitinspect.Repository `json:"git,omitempty"`
	RemoteOrganization string                 `json:"remote_organization,omitempty"`
	RemoteRepository   string                 `json:"remote_repository,omitempty"`
}

func Detect(ctx context.Context, developmentRoot, currentDirectory string, inspector gitinspect.Inspector) (Context, error) {
	root, err := filepath.Abs(developmentRoot)
	if err != nil {
		return Context{}, fmt.Errorf("resolve development root: %w", err)
	}
	current, err := filepath.Abs(currentDirectory)
	if err != nil {
		return Context{}, fmt.Errorf("resolve current directory: %w", err)
	}
	root = filepath.Clean(root)
	current = filepath.Clean(current)

	result := Context{
		DevelopmentRoot:  root,
		CurrentDirectory: current,
		Type:             TypeOther,
	}
	if inspector != nil {
		if repository, ok := inspector.Repository(ctx, current); ok {
			result.IsGitRepository = true
			result.Git = &repository
			result.RemoteOrganization = repository.Owner
			result.RemoteRepository = repository.Name
		}
	}

	relative, err := filepath.Rel(root, current)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return result, nil
	}
	if relative == "." {
		result.Type = TypeDevelopmentRoot
		return result, nil
	}

	parts := splitRelative(relative)
	if len(parts) == 1 {
		result.Type = TypeOrganization
		result.Organization = parts[0]
		return result, nil
	}
	result.Organization = parts[0]
	if result.IsGitRepository {
		repositoryRoot, rootErr := filepath.Rel(root, result.Git.Root)
		repositoryParts := splitRelative(repositoryRoot)
		if rootErr == nil && len(repositoryParts) >= 2 && repositoryParts[0] == parts[0] {
			result.Type = TypeRepository
			result.Repository = repositoryParts[1]
			return result, nil
		}
	}
	return result, nil
}

func splitRelative(relative string) []string {
	return strings.FieldsFunc(relative, func(r rune) bool {
		return r == '/' || r == '\\'
	})
}

// ClassifyShape provides OS-independent tests and consumers with a lexical model
// of ~/dev/<organization>/<repository>. Runtime detection still uses filepath.
func ClassifyShape(relative string, isGitRepository bool) (Type, string, string) {
	relative = strings.Trim(relative, `/\`)
	if relative == "" || relative == "." {
		return TypeDevelopmentRoot, "", ""
	}
	parts := splitRelative(relative)
	if len(parts) == 1 {
		return TypeOrganization, parts[0], ""
	}
	if isGitRepository {
		return TypeRepository, parts[0], parts[1]
	}
	return TypeOther, parts[0], ""
}
