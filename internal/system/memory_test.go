package system

import (
	"context"
	"testing"

	"github.com/shirou/gopsutil/v4/mem"
)

type fakeMemorySource struct {
	value *mem.VirtualMemoryStat
}

func (f fakeMemorySource) VirtualMemory(context.Context) (*mem.VirtualMemoryStat, error) {
	return f.value, nil
}

func TestMemoryCollectorUsesByteUnits(t *testing.T) {
	snapshot := Snapshot{}
	collector := MemoryCollector{Source: fakeMemorySource{value: &mem.VirtualMemoryStat{
		Total:     8 * 1024,
		Available: 5 * 1024,
		Used:      3 * 1024,
	}}}
	if err := collector.Collect(context.Background(), &snapshot); err != nil {
		t.Fatal(err)
	}
	if *snapshot.Memory.TotalBytes != 8*1024 || *snapshot.Memory.UsedBytes != 3*1024 {
		t.Fatalf("Memory = %#v", snapshot.Memory)
	}
	if *snapshot.Memory.UtilizationPercent != 37.5 {
		t.Fatalf("utilization = %v", *snapshot.Memory.UtilizationPercent)
	}
}
