package system

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type fakeCPUSource struct {
	info        []cpu.InfoStat
	physical    int
	logical     int
	percentages []float64
	percentErr  error
}

func (f fakeCPUSource) Info(context.Context) ([]cpu.InfoStat, error) {
	return f.info, nil
}

func (f fakeCPUSource) Counts(_ context.Context, logical bool) (int, error) {
	if logical {
		return f.logical, nil
	}
	return f.physical, nil
}

func (f fakeCPUSource) Percent(context.Context, time.Duration, bool) ([]float64, error) {
	return f.percentages, f.percentErr
}

func TestCPUCollectorCalculatesSnapshot(t *testing.T) {
	snapshot := Snapshot{}
	collector := CPUCollector{
		Source: fakeCPUSource{
			info:        []cpu.InfoStat{{ModelName: "Test CPU"}},
			physical:    4,
			logical:     8,
			percentages: []float64{17.25},
		},
		SampleDuration: 250 * time.Millisecond,
	}
	if err := collector.Collect(context.Background(), &snapshot); err != nil {
		t.Fatal(err)
	}
	if snapshot.CPU.Model != "Test CPU" || *snapshot.CPU.PhysicalCores != 4 || *snapshot.CPU.LogicalCPUs != 8 {
		t.Fatalf("CPU = %#v", snapshot.CPU)
	}
	if *snapshot.CPU.UtilizationPercent != 17.25 || snapshot.CPU.SampleDurationMillis != 250 {
		t.Fatalf("CPU utilization = %#v", snapshot.CPU)
	}
}

func TestCPUCollectorRetainsCountsWhenUtilizationFails(t *testing.T) {
	snapshot := Snapshot{}
	err := (CPUCollector{Source: fakeCPUSource{
		physical:   2,
		logical:    4,
		percentErr: errors.New("not supported"),
	}}).Collect(context.Background(), &snapshot)
	if err == nil || snapshot.CPU.PhysicalCores == nil || snapshot.CPU.LogicalCPUs == nil {
		t.Fatalf("Collect() error = %v, CPU = %#v", err, snapshot.CPU)
	}
}
