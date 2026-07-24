package system

import (
	"context"
	"net"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type HostSource interface {
	Info(context.Context) (*host.InfoStat, error)
}

type GopsutilHostSource struct{}

func (GopsutilHostSource) Info(ctx context.Context) (*host.InfoStat, error) {
	return host.InfoWithContext(ctx)
}

type CPUSource interface {
	Info(context.Context) ([]cpu.InfoStat, error)
	Counts(context.Context, bool) (int, error)
	Percent(context.Context, time.Duration, bool) ([]float64, error)
}

type GopsutilCPUSource struct{}

func (GopsutilCPUSource) Info(ctx context.Context) ([]cpu.InfoStat, error) {
	return cpu.InfoWithContext(ctx)
}

func (GopsutilCPUSource) Counts(ctx context.Context, logical bool) (int, error) {
	return cpu.CountsWithContext(ctx, logical)
}

func (GopsutilCPUSource) Percent(ctx context.Context, interval time.Duration, perCPU bool) ([]float64, error) {
	return cpu.PercentWithContext(ctx, interval, perCPU)
}

type MemorySource interface {
	VirtualMemory(context.Context) (*mem.VirtualMemoryStat, error)
}

type GopsutilMemorySource struct{}

func (GopsutilMemorySource) VirtualMemory(ctx context.Context) (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemoryWithContext(ctx)
}

type DiskSource interface {
	Partitions(context.Context, bool) ([]disk.PartitionStat, error)
	Usage(context.Context, string) (*disk.UsageStat, error)
}

type GopsutilDiskSource struct{}

func (GopsutilDiskSource) Partitions(ctx context.Context, all bool) ([]disk.PartitionStat, error) {
	return disk.PartitionsWithContext(ctx, all)
}

func (GopsutilDiskSource) Usage(ctx context.Context, path string) (*disk.UsageStat, error) {
	return disk.UsageWithContext(ctx, path)
}

type NetworkSource interface {
	Interfaces() ([]net.Interface, error)
	Addrs(net.Interface) ([]net.Addr, error)
}

type StandardNetworkSource struct{}

func (StandardNetworkSource) Interfaces() ([]net.Interface, error) {
	return net.Interfaces()
}

func (StandardNetworkSource) Addrs(value net.Interface) ([]net.Addr, error) {
	return value.Addrs()
}
