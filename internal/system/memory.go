package system

import (
	"context"
	"fmt"
)

type MemoryCollector struct {
	Source MemorySource
}

func (MemoryCollector) Name() string { return "memory" }

func (c MemoryCollector) Collect(ctx context.Context, snapshot *Snapshot) error {
	source := c.Source
	if source == nil {
		source = GopsutilMemorySource{}
	}
	memory, err := source.VirtualMemory(ctx)
	if err != nil {
		return fmt.Errorf("virtual memory: %w", err)
	}
	snapshot.Memory.TotalBytes = uint64Pointer(memory.Total)
	snapshot.Memory.AvailableBytes = uint64Pointer(memory.Available)
	snapshot.Memory.UsedBytes = uint64Pointer(memory.Used)
	snapshot.Memory.UtilizationPercent = floatPointer(Utilization(memory.Used, memory.Total))
	return nil
}
