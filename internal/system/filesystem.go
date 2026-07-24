package system

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strings"
)

var excludedFilesystemTypes = map[string]struct{}{
	"autofs":      {},
	"binfmt_misc": {},
	"cgroup":      {},
	"cgroup2":     {},
	"configfs":    {},
	"debugfs":     {},
	"devfs":       {},
	"devpts":      {},
	"devtmpfs":    {},
	"efivarfs":    {},
	"fusectl":     {},
	"hugetlbfs":   {},
	"mqueue":      {},
	"nsfs":        {},
	"overlay":     {},
	"proc":        {},
	"procfs":      {},
	"pstore":      {},
	"rpc_pipefs":  {},
	"securityfs":  {},
	"squashfs":    {},
	"sysfs":       {},
	"tmpfs":       {},
	"tracefs":     {},
}

var excludedMountPrefixes = []string{
	"/dev",
	"/proc",
	"/run",
	"/snap",
	"/sys",
	"/var/lib/containers",
	"/var/lib/docker",
	"/System/Volumes",
}

type FilesystemCollector struct {
	Source DiskSource
	GOOS   string
}

func (FilesystemCollector) Name() string { return "filesystems" }

func (c FilesystemCollector) Collect(ctx context.Context, snapshot *Snapshot) error {
	source := c.Source
	if source == nil {
		source = GopsutilDiskSource{}
	}
	goos := c.GOOS
	if goos == "" {
		goos = runtime.GOOS
	}
	partitions, err := source.Partitions(ctx, false)
	var problems []error
	if err != nil {
		if len(partitions) == 0 {
			return fmt.Errorf("list partitions: %w", err)
		}
		problems = append(problems, fmt.Errorf("list partitions: %w", err))
	}
	seen := make(map[string]struct{}, len(partitions))
	for _, partition := range partitions {
		if !ShouldIncludeFilesystem(goos, partition.Mountpoint, partition.Fstype, partition.Device) {
			continue
		}
		key := strings.ToLower(partition.Mountpoint)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		usage, err := source.Usage(ctx, partition.Mountpoint)
		if err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", partition.Mountpoint, err))
			continue
		}
		snapshot.Filesystems = append(snapshot.Filesystems, Filesystem{
			Device:             partition.Device,
			Mountpoint:         partition.Mountpoint,
			Type:               partition.Fstype,
			TotalBytes:         uint64Pointer(usage.Total),
			AvailableBytes:     uint64Pointer(usage.Free),
			UsedBytes:          uint64Pointer(usage.Used),
			UtilizationPercent: floatPointer(Utilization(usage.Used, usage.Total)),
		})
	}
	sort.Slice(snapshot.Filesystems, func(left, right int) bool {
		return strings.ToLower(snapshot.Filesystems[left].Mountpoint) < strings.ToLower(snapshot.Filesystems[right].Mountpoint)
	})
	return errors.Join(problems...)
}

func ShouldIncludeFilesystem(goos, mountpoint, filesystemType, device string) bool {
	if strings.TrimSpace(mountpoint) == "" {
		return false
	}
	if _, excluded := excludedFilesystemTypes[strings.ToLower(filesystemType)]; excluded {
		return false
	}
	normalized := strings.ReplaceAll(mountpoint, "\\", "/")
	if goos == "windows" {
		return len(normalized) >= 2 &&
			normalized[1] == ':' &&
			(len(normalized) == 2 || normalized[2] == '/')
	}
	for _, prefix := range excludedMountPrefixes {
		if normalized == prefix || strings.HasPrefix(normalized, prefix+"/") {
			return false
		}
	}
	if device == "" && filesystemType == "" {
		return false
	}
	return strings.HasPrefix(normalized, "/")
}
