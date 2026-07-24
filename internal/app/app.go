package app

import (
	"context"
	"fmt"
	"os"

	"github.com/setup-env/app/internal/catalog"
	"github.com/setup-env/app/internal/config"
	"github.com/setup-env/app/internal/dashboard"
	"github.com/setup-env/app/internal/diagnostics"
	"github.com/setup-env/app/internal/directory"
	gitinspect "github.com/setup-env/app/internal/git"
	"github.com/setup-env/app/internal/paths"
	"github.com/setup-env/app/internal/platform"
	"github.com/setup-env/app/internal/system"
)

type Info struct {
	Platform      platform.Info        `json:"platform"`
	Directory     directory.Context    `json:"directory"`
	ConfigPath    string               `json:"config_path"`
	ConfigLoaded  bool                 `json:"config_loaded"`
	CloneProtocol config.CloneProtocol `json:"clone_protocol"`
}

type Service struct {
	PlatformDetector platform.Detector
	PathResolver     paths.Resolver
	ConfigLocation   config.LocationResolver
	GitInspector     gitinspect.Inspector
	Getwd            func() (string, error)
	Commands         diagnostics.CommandRunner
	CatalogSource    catalog.Source
	SystemCollector  system.SnapshotCollector
	DashboardSource  dashboard.Source
}

func DefaultService() Service {
	return Service{
		PlatformDetector: platform.DefaultDetector(),
		PathResolver:     paths.DefaultResolver(),
		ConfigLocation:   config.DefaultLocationResolver(),
		GitInspector:     gitinspect.DefaultInspector(),
		Getwd:            os.Getwd,
		Commands:         diagnostics.OSCommandRunner{},
		CatalogSource:    catalog.EmbeddedSource{},
	}
}

func (s Service) Info(ctx context.Context) (Info, error) {
	platformInfo, err := s.PlatformDetector.Detect()
	if err != nil {
		return Info{}, err
	}
	configPath, err := s.ConfigLocation.Path()
	if err != nil {
		return Info{}, err
	}
	settings, loaded, err := config.Load(configPath)
	if err != nil {
		return Info{}, err
	}
	developmentRoot, err := s.PathResolver.DevelopmentRoot(settings.DevelopmentRoot)
	if err != nil {
		return Info{}, err
	}
	current, err := s.Getwd()
	if err != nil {
		return Info{}, fmt.Errorf("resolve current working directory: %w", err)
	}
	directoryContext, err := directory.Detect(ctx, developmentRoot, current, s.GitInspector)
	if err != nil {
		return Info{}, err
	}
	return Info{
		Platform:      platformInfo,
		Directory:     directoryContext,
		ConfigPath:    configPath,
		ConfigLoaded:  loaded,
		CloneProtocol: settings.CloneProtocol,
	}, nil
}

func (s Service) Doctor(ctx context.Context) (diagnostics.Report, error) {
	configPath, err := s.ConfigLocation.Path()
	if err != nil {
		return diagnostics.Report{}, err
	}
	settings, _, err := config.Load(configPath)
	if err != nil {
		return diagnostics.Report{}, err
	}
	developmentRoot, err := s.PathResolver.DevelopmentRoot(settings.DevelopmentRoot)
	if err != nil {
		return diagnostics.Report{}, err
	}
	return diagnostics.Run(ctx, developmentRoot, s.Commands), nil
}
