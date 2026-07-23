package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/setup-env/app/internal/catalog"
	"github.com/setup-env/app/internal/compatibility"
	"github.com/setup-env/app/internal/config"
	"github.com/setup-env/app/internal/manifest"
	"github.com/setup-env/app/internal/version"
)

type ModuleListOptions struct {
	Trust    catalog.Trust
	Status   catalog.Status
	Category string
}

type ModuleInfo struct {
	CatalogEntry      catalog.Entry        `json:"catalog_entry"`
	ManifestPath      string               `json:"manifest_path"`
	ManifestAvailable bool                 `json:"manifest_available"`
	Manifest          *manifest.Manifest   `json:"manifest,omitempty"`
	Compatibility     compatibility.Result `json:"compatibility"`
	Message           string               `json:"message"`
}

type ManifestValidationReport struct {
	Path          string               `json:"path"`
	Valid         bool                 `json:"valid"`
	Problems      []string             `json:"problems,omitempty"`
	Manifest      *manifest.Manifest   `json:"manifest,omitempty"`
	Compatibility compatibility.Result `json:"compatibility"`
}

type CatalogValidationReport struct {
	Valid       bool     `json:"valid"`
	ModuleCount int      `json:"module_count"`
	Problems    []string `json:"problems,omitempty"`
}

func (s Service) ModuleList(ctx context.Context, options ModuleListOptions) ([]catalog.Entry, error) {
	current, err := catalog.Load(ctx, s.CatalogSource)
	if err != nil {
		return nil, err
	}
	return catalog.Filter(current.Modules, options.Trust, options.Status, options.Category), nil
}

func (s Service) ModuleInfo(ctx context.Context, id string) (ModuleInfo, error) {
	current, err := catalog.Load(ctx, s.CatalogSource)
	if err != nil {
		return ModuleInfo{}, err
	}
	entry, ok := catalog.Find(current.Modules, id)
	if !ok {
		return ModuleInfo{}, fmt.Errorf("module %q is not listed in the official catalog", id)
	}
	developmentRoot, err := s.developmentRoot()
	if err != nil {
		return ModuleInfo{}, err
	}
	owner, repository := entry.RepositoryParts()
	manifestPath := filepath.Join(developmentRoot, owner, repository, filepath.FromSlash(entry.Manifest))
	result := ModuleInfo{
		CatalogEntry: entry,
		ManifestPath: manifestPath,
		Compatibility: compatibility.Result{
			State:      compatibility.StateUnknown,
			AppVersion: version.Current().Version,
			Reason:     "module manifest is not available locally",
		},
		Message: "The module repository or manifest has not been fetched locally. Downloading is not implemented.",
	}
	if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
		return result, nil
	} else if err != nil {
		return ModuleInfo{}, fmt.Errorf("inspect local manifest %q: %w", manifestPath, err)
	}

	moduleManifest, err := manifest.ParseFile(manifestPath)
	if err != nil {
		return ModuleInfo{}, err
	}
	if err := manifest.Validate(moduleManifest); err != nil {
		return ModuleInfo{}, fmt.Errorf("validate local manifest %q: %w", manifestPath, err)
	}
	if err := matchCatalogEntry(entry, moduleManifest); err != nil {
		return ModuleInfo{}, err
	}
	result.ManifestAvailable = true
	result.Manifest = &moduleManifest
	result.Compatibility = compatibility.Evaluate(version.Current().Version, moduleManifest.MinimumAppVersion)
	result.Message = "Local manifest loaded. Downloading, installation, and workflow execution are not implemented."
	return result, nil
}

func (s Service) ValidateManifest(target string) ManifestValidationReport {
	path := target
	if info, err := os.Stat(target); err == nil && info.IsDir() {
		path = filepath.Join(target, "setup-env.yaml")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return ManifestValidationReport{Path: path, Problems: []string{"resolve manifest path: " + err.Error()}}
	}
	result := ManifestValidationReport{
		Path: absolute,
		Compatibility: compatibility.Result{
			State:      compatibility.StateUnknown,
			AppVersion: version.Current().Version,
			Reason:     "manifest has not passed validation",
		},
	}
	moduleManifest, err := manifest.ParseFile(absolute)
	if err != nil {
		result.Problems = []string{err.Error()}
		return result
	}
	result.Manifest = &moduleManifest
	if err := manifest.Validate(moduleManifest); err != nil {
		var validationError *manifest.ValidationError
		if errors.As(err, &validationError) {
			result.Problems = append(result.Problems, validationError.Problems...)
		} else {
			result.Problems = []string{err.Error()}
		}
		return result
	}
	result.Valid = true
	result.Compatibility = compatibility.Evaluate(version.Current().Version, moduleManifest.MinimumAppVersion)
	return result
}

func (s Service) ValidateCatalog(ctx context.Context) CatalogValidationReport {
	current, err := catalog.Load(ctx, s.CatalogSource)
	if err == nil {
		return CatalogValidationReport{Valid: true, ModuleCount: len(current.Modules)}
	}
	result := CatalogValidationReport{}
	var validationError *catalog.ValidationError
	if errors.As(err, &validationError) {
		result.Problems = append(result.Problems, validationError.Problems...)
	} else {
		result.Problems = []string{err.Error()}
	}
	return result
}

func (s Service) developmentRoot() (string, error) {
	configPath, err := s.ConfigLocation.Path()
	if err != nil {
		return "", err
	}
	settings, _, err := config.Load(configPath)
	if err != nil {
		return "", err
	}
	return s.PathResolver.DevelopmentRoot(settings.DevelopmentRoot)
}

func matchCatalogEntry(entry catalog.Entry, moduleManifest manifest.Manifest) error {
	owner, repository := entry.RepositoryParts()
	var problems []string
	if entry.ID != moduleManifest.ID {
		problems = append(problems, fmt.Sprintf("catalog id %q does not match manifest id %q", entry.ID, moduleManifest.ID))
	}
	if owner != moduleManifest.Repository.Owner || repository != moduleManifest.Repository.Name {
		problems = append(problems, fmt.Sprintf("catalog repository %q does not match manifest repository %s/%s", entry.Repository, moduleManifest.Repository.Owner, moduleManifest.Repository.Name))
	}
	if len(problems) > 0 {
		return fmt.Errorf("catalog and manifest disagree: %v", problems)
	}
	return nil
}
