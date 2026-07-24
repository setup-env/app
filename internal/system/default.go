package system

import "github.com/setup-env/app/internal/platform"

func DefaultSections(additional ...SectionCollector) []SectionCollector {
	sections := []SectionCollector{
		HostCollector{
			Source:           GopsutilHostSource{},
			PlatformDetector: platform.DefaultDetector(),
		},
		CPUCollector{
			Source:         GopsutilCPUSource{},
			SampleDuration: DefaultCPUSampleDuration,
		},
		MemoryCollector{Source: GopsutilMemorySource{}},
		FilesystemCollector{Source: GopsutilDiskSource{}},
		NetworkCollector{Source: StandardNetworkSource{}},
	}
	return append(sections, additional...)
}
