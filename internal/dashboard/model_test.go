package dashboard

import (
	"context"
	"errors"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/setup-env/app/internal/system"
)

type fakeSource struct {
	initial system.Snapshot
	refresh func(context.Context, RefreshRequest) (system.Snapshot, error)
}

func (f fakeSource) Initial(context.Context) (system.Snapshot, error) {
	return f.initial, nil
}

func (f fakeSource) Refresh(ctx context.Context, request RefreshRequest) (system.Snapshot, error) {
	return f.refresh(ctx, request)
}

func TestModelForcedRefreshUpdatesHistoriesAndRates(t *testing.T) {
	initial := dashboardTestSnapshot()
	oldReceived := uint64(100)
	initial.Networks[0].BytesReceived = &oldReceived
	current := initial
	current.Timestamp = current.Timestamp.Add(time.Second)
	cpu, memory := 75.0, 60.0
	current.CPU.UtilizationPercent = &cpu
	current.Memory.UtilizationPercent = &memory
	newReceived := uint64(1124)
	current.Networks = append([]system.NetworkInterface(nil), current.Networks...)
	current.Networks[0].BytesReceived = &newReceived

	var request RefreshRequest
	source := fakeSource{initial: initial, refresh: func(_ context.Context, value RefreshRequest) (system.Snapshot, error) {
		request = value
		return current, nil
	}}
	model := NewModel(context.Background(), source, initial, ModelOptions{
		Now:  func() time.Time { return current.Timestamp },
		Tick: func(time.Duration, func(time.Time) tea.Msg) tea.Cmd { return nil },
	})
	updatedValue, command := model.handleKey("r")
	updated := updatedValue.(Model)
	if !updated.collecting || command == nil {
		t.Fatalf("collecting = %t, command = %v", updated.collecting, command)
	}
	message := command()
	finalValue, _ := updated.Update(message)
	final := finalValue.(Model)
	if !request.IncludeFilesystems || !request.IncludeDiagnostics {
		t.Fatalf("forced request = %#v", request)
	}
	if final.collecting || final.cpuHistory.Len() != 2 || final.memoryHistory.Len() != 2 {
		t.Fatalf("model = %#v", final)
	}
	if len(final.networkRates) != 1 || *final.networkRates[0].BytesReceivedPerSec != 1024 {
		t.Fatalf("rates = %#v", final.networkRates)
	}
}

func TestModelPreventsOverlappingRefreshAndSupportsPause(t *testing.T) {
	initial := dashboardTestSnapshot()
	source := fakeSource{initial: initial, refresh: func(context.Context, RefreshRequest) (system.Snapshot, error) {
		return initial, nil
	}}
	model := NewModel(context.Background(), source, initial, ModelOptions{
		Tick: func(time.Duration, func(time.Time) tea.Msg) tea.Cmd { return nil },
	})
	model.collecting = true
	updatedValue, command := model.handleKey("r")
	if command != nil || !updatedValue.(Model).collecting {
		t.Fatal("forced refresh overlapped active collection")
	}
	model.collecting = false
	pausedValue, _ := model.handleKey("p")
	paused := pausedValue.(Model)
	if !paused.paused {
		t.Fatal("pause key did not pause")
	}
	tickValue, _ := paused.Update(tickMsg{at: initial.Timestamp.Add(time.Second)})
	if tickValue.(Model).collecting {
		t.Fatal("paused tick started collection")
	}
}

func TestPartialRefreshUpdatesSnapshotAndKeepsWarning(t *testing.T) {
	initial := dashboardTestSnapshot()
	current := initial
	current.Timestamp = current.Timestamp.Add(time.Second)
	partialErr := errors.New("network counters unavailable")
	source := fakeSource{initial: initial, refresh: func(context.Context, RefreshRequest) (system.Snapshot, error) {
		return current, partialErr
	}}
	model := NewModel(context.Background(), source, initial, ModelOptions{
		Now:  func() time.Time { return current.Timestamp },
		Tick: func(time.Duration, func(time.Time) tea.Msg) tea.Cmd { return nil },
	})
	started, command := model.startRefresh(false)
	finalValue, _ := started.Update(command())
	final := finalValue.(Model)
	if final.snapshot.Timestamp != current.Timestamp || !errors.Is(final.lastError, partialErr) {
		t.Fatalf("model = %#v", final)
	}
}

func TestRefreshHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	initial := dashboardTestSnapshot()
	source := fakeSource{initial: initial, refresh: func(ctx context.Context, _ RefreshRequest) (system.Snapshot, error) {
		<-ctx.Done()
		return system.Snapshot{}, ctx.Err()
	}}
	model := NewModel(ctx, source, initial, ModelOptions{
		Now:  func() time.Time { return initial.Timestamp.Add(time.Second) },
		Tick: func(time.Duration, func(time.Time) tea.Msg) tea.Cmd { return nil },
	})
	_, command := model.startRefresh(false)
	cancel()
	message := command().(refreshMsg)
	if message.err != context.Canceled {
		t.Fatalf("error = %v", message.err)
	}
}

func TestKeyCommandsAndAltScreenLifecycle(t *testing.T) {
	tests := map[string]Command{
		"q":      CommandQuit,
		"ctrl+c": CommandQuit,
		"r":      CommandRefresh,
		"?":      CommandHelp,
		"p":      CommandPause,
		"space":  CommandPause,
		"x":      CommandNone,
	}
	for key, want := range tests {
		if got := KeyCommand(key); got != want {
			t.Fatalf("KeyCommand(%q) = %v, want %v", key, got, want)
		}
	}
	model := NewModel(context.Background(), nil, dashboardTestSnapshot(), ModelOptions{})
	model.width, model.height = 80, 24
	if view := model.View(); !view.AltScreen {
		t.Fatal("dashboard view did not request alternate-screen lifecycle")
	}
}
