package system

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

const DefaultCPUSampleDuration = 500 * time.Millisecond

type CPUCollector struct {
	Source         CPUSource
	SampleDuration time.Duration
}

func (CPUCollector) Name() string { return "cpu" }

func (c CPUCollector) Collect(ctx context.Context, snapshot *Snapshot) error {
	source := c.Source
	if source == nil {
		source = GopsutilCPUSource{}
	}
	sampleDuration := c.SampleDuration
	if sampleDuration <= 0 {
		sampleDuration = DefaultCPUSampleDuration
	}
	snapshot.CPU.SampleDurationMillis = sampleDuration.Milliseconds()

	var problems []error
	info, err := source.Info(ctx)
	if err != nil {
		problems = append(problems, fmt.Errorf("model: %w", err))
	} else {
		for _, item := range info {
			if model := strings.TrimSpace(item.ModelName); model != "" {
				snapshot.CPU.Model = model
				break
			}
		}
	}
	physical, err := source.Counts(ctx, false)
	if err != nil {
		problems = append(problems, fmt.Errorf("physical cores: %w", err))
	} else {
		snapshot.CPU.PhysicalCores = intPointer(physical)
	}
	logical, err := source.Counts(ctx, true)
	if err != nil {
		problems = append(problems, fmt.Errorf("logical CPUs: %w", err))
	} else {
		snapshot.CPU.LogicalCPUs = intPointer(logical)
	}
	percentages, err := source.Percent(ctx, sampleDuration, false)
	if err != nil {
		problems = append(problems, fmt.Errorf("utilization: %w", err))
	} else if len(percentages) == 0 {
		problems = append(problems, fmt.Errorf("utilization: no samples returned"))
	} else {
		snapshot.CPU.UtilizationPercent = floatPointer(clampPercentage(percentages[0]))
	}
	return errors.Join(problems...)
}
