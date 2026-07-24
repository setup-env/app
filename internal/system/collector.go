package system

import (
	"context"
	"fmt"
	"time"
)

type SnapshotCollector interface {
	Collect(context.Context) (Snapshot, error)
}

type SectionCollector interface {
	Name() string
	Collect(context.Context, *Snapshot) error
}

type Collector struct {
	Now      func() time.Time
	Sections []SectionCollector
}

func (c Collector) Collect(ctx context.Context) (Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return Snapshot{}, fmt.Errorf("collect system snapshot: %w", err)
	}
	now := time.Now
	if c.Now != nil {
		now = c.Now
	}
	timestamp := now()
	zoneName, offset := timestamp.Zone()
	result := Snapshot{
		SchemaVersion: SnapshotSchemaVersion,
		Timestamp:     timestamp,
		TimeZone: TimeZone{
			Name:             zoneName,
			UTCOffsetSeconds: offset,
		},
		Filesystems: []Filesystem{},
		Networks:    []NetworkInterface{},
		Warnings:    []Warning{},
	}

	successes := 0
	for _, section := range c.Sections {
		if err := ctx.Err(); err != nil {
			result.Warnings = append(result.Warnings, Warning{Section: section.Name(), Message: err.Error()})
			break
		}
		if err := section.Collect(ctx, &result); err != nil {
			result.Warnings = append(result.Warnings, Warning{Section: section.Name(), Message: err.Error()})
			continue
		}
		successes++
	}
	FinalizeHealth(&result)
	if len(c.Sections) > 0 && successes == 0 {
		return result, fmt.Errorf("collect system snapshot: every section failed")
	}
	return result, nil
}

func FinalizeHealth(snapshot *Snapshot) {
	warnings := len(snapshot.Warnings)
	failures := 0
	for _, check := range snapshot.Diagnostics.Checks {
		switch check.Status {
		case "warn":
			warnings++
		case "fail":
			failures++
		}
	}
	snapshot.Diagnostics.WarningCount = warnings
	snapshot.Diagnostics.FailureCount = failures
	switch {
	case failures > 0:
		snapshot.Diagnostics.Health = HealthUnhealthy
	case warnings > 0:
		snapshot.Diagnostics.Health = HealthWarning
	default:
		snapshot.Diagnostics.Health = HealthHealthy
	}
}
