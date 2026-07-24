package system

import (
	"context"
	"errors"
	"testing"

	"github.com/shirou/gopsutil/v4/disk"
)

type fakeDiskSource struct {
	partitions     []disk.PartitionStat
	partitionError error
	usage          map[string]*disk.UsageStat
	errors         map[string]error
}

func (f fakeDiskSource) Partitions(context.Context, bool) ([]disk.PartitionStat, error) {
	return f.partitions, f.partitionError
}

func (f fakeDiskSource) Usage(_ context.Context, path string) (*disk.UsageStat, error) {
	return f.usage[path], f.errors[path]
}

func TestFilesystemFiltering(t *testing.T) {
	tests := []struct {
		name       string
		goos       string
		mountpoint string
		fsType     string
		device     string
		want       bool
	}{
		{name: "Linux root", goos: "linux", mountpoint: "/", fsType: "ext4", device: "/dev/sda1", want: true},
		{name: "Linux home", goos: "linux", mountpoint: "/home", fsType: "xfs", device: "/dev/sdb1", want: true},
		{name: "proc", goos: "linux", mountpoint: "/proc", fsType: "proc", want: false},
		{name: "docker", goos: "linux", mountpoint: "/var/lib/docker/overlay2", fsType: "ext4", device: "/dev/sda1", want: false},
		{name: "macOS internal", goos: "darwin", mountpoint: "/System/Volumes/Data", fsType: "apfs", device: "/dev/disk3s1", want: false},
		{name: "Windows drive", goos: "windows", mountpoint: `C:\`, fsType: "NTFS", device: `C:\`, want: true},
		{name: "Windows drive without slash", goos: "windows", mountpoint: `C:`, fsType: "NTFS", device: `C:`, want: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ShouldIncludeFilesystem(test.goos, test.mountpoint, test.fsType, test.device); got != test.want {
				t.Fatalf("ShouldIncludeFilesystem() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestFilesystemCollectorKeepsPartitionsReturnedWithWarning(t *testing.T) {
	source := fakeDiskSource{
		partitions:     []disk.PartitionStat{{Mountpoint: `C:\`, Fstype: "NTFS", Device: `C:\`}},
		partitionError: errors.New("one volume could not be inspected"),
		usage: map[string]*disk.UsageStat{
			`C:\`: {Total: 100, Used: 20, Free: 80},
		},
		errors: map[string]error{},
	}
	snapshot := Snapshot{}
	err := (FilesystemCollector{Source: source, GOOS: "windows"}).Collect(context.Background(), &snapshot)
	if err == nil {
		t.Fatal("Collect() error = nil, want partial warning")
	}
	if len(snapshot.Filesystems) != 1 || snapshot.Filesystems[0].Mountpoint != `C:\` {
		t.Fatalf("filesystems = %#v", snapshot.Filesystems)
	}
}

func TestFilesystemCollectorSortsAndKeepsSuccessfulMounts(t *testing.T) {
	source := fakeDiskSource{
		partitions: []disk.PartitionStat{
			{Mountpoint: "/home", Fstype: "ext4", Device: "/dev/sdb1"},
			{Mountpoint: "/", Fstype: "ext4", Device: "/dev/sda1"},
			{Mountpoint: "/data", Fstype: "xfs", Device: "/dev/sdc1"},
		},
		usage: map[string]*disk.UsageStat{
			"/home": {Total: 200, Used: 50, Free: 150},
			"/":     {Total: 100, Used: 25, Free: 75},
		},
		errors: map[string]error{"/data": errors.New("permission denied")},
	}
	snapshot := Snapshot{}
	err := (FilesystemCollector{Source: source, GOOS: "linux"}).Collect(context.Background(), &snapshot)
	if err == nil {
		t.Fatal("Collect() error = nil, want partial warning")
	}
	if len(snapshot.Filesystems) != 2 || snapshot.Filesystems[0].Mountpoint != "/" || snapshot.Filesystems[1].Mountpoint != "/home" {
		t.Fatalf("filesystems = %#v", snapshot.Filesystems)
	}
	if *snapshot.Filesystems[0].UtilizationPercent != 25 {
		t.Fatalf("utilization = %v", *snapshot.Filesystems[0].UtilizationPercent)
	}
}
