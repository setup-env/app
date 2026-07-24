package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/setup-env/app/internal/diagnostics"
	"github.com/setup-env/app/internal/system"
	"github.com/setup-env/app/internal/version"
)

func (s Service) Status(ctx context.Context) (system.Snapshot, error) {
	if s.SystemCollector != nil {
		return s.SystemCollector.Collect(ctx)
	}
	collector := system.Collector{
		Sections: system.DefaultSections(DevelopmentCollector{Service: s}),
	}
	return collector.Collect(ctx)
}

type DevelopmentCollector struct {
	Service Service
}

func (DevelopmentCollector) Name() string { return "development" }

func (c DevelopmentCollector) Collect(ctx context.Context, snapshot *system.Snapshot) error {
	info, infoErr := c.Service.Info(ctx)
	report, doctorErr := c.Service.Doctor(ctx)

	snapshot.Development.ApplicationVersion = version.Current()
	if infoErr == nil {
		snapshot.Development.Directory = info.Directory
	}
	if doctorErr == nil {
		snapshot.Development.Root = report.Details.DevelopmentRoot
		snapshot.Development.Git = report.Details.Git
		snapshot.Development.GitHubCLI = report.Details.GitHubCLI
		snapshot.Development.GitHubAuthenticated = report.Details.GitHubAuthenticated
		snapshot.Diagnostics.Checks = append([]diagnostics.Check(nil), report.Checks...)
	}

	var problems []error
	if infoErr != nil {
		problems = append(problems, fmt.Errorf("environment information: %w", infoErr))
	}
	if doctorErr != nil {
		problems = append(problems, fmt.Errorf("diagnostics: %w", doctorErr))
	}
	return errors.Join(problems...)
}
