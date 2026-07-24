package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/setup-env/app/internal/dashboard"
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

func (s Service) Dashboard() dashboard.Source {
	if s.DashboardSource != nil {
		return s.DashboardSource
	}
	return LiveCollector{Service: s}
}

type LiveCollector struct {
	Service Service
}

func (c LiveCollector) Initial(ctx context.Context) (system.Snapshot, error) {
	return c.Service.Status(ctx)
}

func (c LiveCollector) Refresh(ctx context.Context, request dashboard.RefreshRequest) (system.Snapshot, error) {
	sections := []system.SectionCollector{
		system.CPUCollector{
			Source:         system.GopsutilCPUSource{},
			SampleDuration: system.DefaultCPUSampleDuration,
		},
		system.MemoryCollector{Source: system.GopsutilMemorySource{}},
		system.NetworkCollector{
			Source:        system.StandardNetworkSource{},
			CounterSource: system.GopsutilNetworkCounterSource{},
		},
	}
	if request.IncludeFilesystems {
		sections = append(sections, system.FilesystemCollector{Source: system.GopsutilDiskSource{}})
	}
	if request.IncludeDiagnostics {
		sections = append(sections, DevelopmentCollector{Service: c.Service})
	}
	fresh, collectErr := (system.Collector{
		Now:      func() time.Time { return request.CollectedAt },
		Sections: sections,
	}).Collect(ctx)

	result := request.Previous
	result.Timestamp = fresh.Timestamp
	result.TimeZone = fresh.TimeZone
	result.Warnings = fresh.Warnings
	result.CPU = mergeCPU(result.CPU, fresh.CPU)
	result.Memory = mergeMemory(result.Memory, fresh.Memory)
	result.Networks = fresh.Networks
	if request.IncludeFilesystems {
		result.Filesystems = fresh.Filesystems
	}
	if request.IncludeDiagnostics {
		result.Development = fresh.Development
		result.Diagnostics.Checks = fresh.Diagnostics.Checks
	}
	if result.Host.UptimeSeconds != nil {
		uptime := *result.Host.UptimeSeconds
		if elapsed := request.CollectedAt.Sub(request.Previous.Timestamp); elapsed > 0 {
			uptime += uint64(elapsed / time.Second)
		}
		result.Host.UptimeSeconds = &uptime
	}
	system.FinalizeHealth(&result)
	return result, collectErr
}

func mergeCPU(previous, current system.CPU) system.CPU {
	if current.Model != "" {
		previous.Model = current.Model
	}
	if current.PhysicalCores != nil {
		previous.PhysicalCores = current.PhysicalCores
	}
	if current.LogicalCPUs != nil {
		previous.LogicalCPUs = current.LogicalCPUs
	}
	if current.UtilizationPercent != nil {
		previous.UtilizationPercent = current.UtilizationPercent
	}
	if current.SampleDurationMillis > 0 {
		previous.SampleDurationMillis = current.SampleDurationMillis
	}
	return previous
}

func mergeMemory(previous, current system.Memory) system.Memory {
	if current.TotalBytes != nil {
		previous.TotalBytes = current.TotalBytes
	}
	if current.AvailableBytes != nil {
		previous.AvailableBytes = current.AvailableBytes
	}
	if current.UsedBytes != nil {
		previous.UsedBytes = current.UsedBytes
	}
	if current.UtilizationPercent != nil {
		previous.UtilizationPercent = current.UtilizationPercent
	}
	return previous
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
