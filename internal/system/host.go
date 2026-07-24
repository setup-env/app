package system

import (
	"context"
	"fmt"
	"runtime"

	"github.com/setup-env/app/internal/platform"
)

type HostCollector struct {
	Source           HostSource
	PlatformDetector platform.Detector
}

func (HostCollector) Name() string { return "host" }

func (c HostCollector) Collect(ctx context.Context, snapshot *Snapshot) error {
	source := c.Source
	if source == nil {
		source = GopsutilHostSource{}
	}
	detector := c.PlatformDetector
	if detector.CurrentUser == nil {
		detector = platform.DefaultDetector()
	}
	identity, identityErr := detector.Detect()
	info, hostErr := source.Info(ctx)

	if identityErr == nil {
		snapshot.User = User{Username: identity.User, Home: identity.Home, Shell: identity.Shell}
		snapshot.OperatingSystem.OS = identity.OS
		snapshot.OperatingSystem.Architecture = identity.Architecture
		snapshot.OperatingSystem.DisplayName = identity.Distribution
	}
	if info != nil {
		snapshot.Host.Hostname = info.Hostname
		snapshot.Host.UptimeSeconds = uint64Pointer(info.Uptime)
		if snapshot.OperatingSystem.OS == "" {
			snapshot.OperatingSystem.OS = info.OS
		}
		if snapshot.OperatingSystem.Architecture == "" {
			snapshot.OperatingSystem.Architecture = runtime.GOARCH
		}
		snapshot.OperatingSystem.Version = info.PlatformVersion
		snapshot.OperatingSystem.Distribution = info.Platform
		snapshot.OperatingSystem.DistributionVersion = info.PlatformVersion
		snapshot.OperatingSystem.KernelVersion = info.KernelVersion
		if info.OS == "windows" {
			snapshot.OperatingSystem.BuildInformation = info.KernelVersion
		}
		if snapshot.OperatingSystem.DisplayName == "" {
			snapshot.OperatingSystem.DisplayName = info.Platform
		}
	}
	switch {
	case identityErr != nil && hostErr != nil:
		return fmt.Errorf("platform identity: %v; host information: %v", identityErr, hostErr)
	case identityErr != nil:
		return fmt.Errorf("platform identity: %w", identityErr)
	case hostErr != nil:
		return fmt.Errorf("host information: %w", hostErr)
	default:
		return nil
	}
}
