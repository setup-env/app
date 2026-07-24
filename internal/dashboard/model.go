package dashboard

import (
	"context"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/setup-env/app/internal/system"
)

type tickMsg struct {
	at time.Time
}

type refreshMsg struct {
	snapshot system.Snapshot
	request  RefreshRequest
	err      error
}

type Model struct {
	ctx                context.Context
	source             Source
	snapshot           system.Snapshot
	clock              time.Time
	cpuHistory         History
	memoryHistory      History
	networkRates       []NetworkRate
	width              int
	height             int
	help               bool
	paused             bool
	collecting         bool
	lastError          error
	lastFilesystem     time.Time
	lastDiagnostics    time.Time
	refreshInterval    time.Duration
	filesystemInterval time.Duration
	diagnosticInterval time.Duration
	tick               func(time.Duration, func(time.Time) tea.Msg) tea.Cmd
}

type ModelOptions struct {
	Now                func() time.Time
	RefreshInterval    time.Duration
	FilesystemInterval time.Duration
	DiagnosticInterval time.Duration
	HistoryLimit       int
	Tick               func(time.Duration, func(time.Time) tea.Msg) tea.Cmd
}

func NewModel(ctx context.Context, source Source, snapshot system.Snapshot, options ModelOptions) Model {
	now := time.Now
	if options.Now != nil {
		now = options.Now
	}
	refreshInterval := options.RefreshInterval
	if refreshInterval <= 0 {
		refreshInterval = DefaultRefreshInterval
	}
	filesystemInterval := options.FilesystemInterval
	if filesystemInterval <= 0 {
		filesystemInterval = DefaultFilesystemInterval
	}
	diagnosticInterval := options.DiagnosticInterval
	if diagnosticInterval <= 0 {
		diagnosticInterval = DefaultDiagnosticInterval
	}
	historyLimit := options.HistoryLimit
	if historyLimit <= 0 {
		historyLimit = DefaultHistoryLimit
	}
	tick := options.Tick
	if tick == nil {
		tick = tea.Tick
	}
	cpuHistory := NewHistory(historyLimit)
	if snapshot.CPU.UtilizationPercent != nil {
		cpuHistory = cpuHistory.Add(*snapshot.CPU.UtilizationPercent)
	}
	memoryHistory := NewHistory(historyLimit)
	if snapshot.Memory.UtilizationPercent != nil {
		memoryHistory = memoryHistory.Add(*snapshot.Memory.UtilizationPercent)
	}
	clock := now()
	return Model{
		ctx:                ctx,
		source:             source,
		snapshot:           snapshot,
		clock:              clock,
		cpuHistory:         cpuHistory,
		memoryHistory:      memoryHistory,
		lastFilesystem:     snapshot.Timestamp,
		lastDiagnostics:    snapshot.Timestamp,
		refreshInterval:    refreshInterval,
		filesystemInterval: filesystemInterval,
		diagnosticInterval: diagnosticInterval,
		tick:               tick,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.nextTick(), requestWindowSize())
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width = max(0, message.Width)
		m.height = max(0, message.Height)
		return m, nil
	case tea.KeyPressMsg:
		return m.handleKey(message.String())
	case tea.InterruptMsg:
		return m, tea.Quit
	case tickMsg:
		m.clock = message.at
		next := tea.Batch(m.nextTick(), requestWindowSize())
		if m.paused || m.collecting {
			return m, next
		}
		var refresh tea.Cmd
		m, refresh = m.startRefresh(false)
		return m, tea.Batch(next, refresh)
	case refreshMsg:
		m.collecting = false
		m.lastError = message.err
		if message.snapshot.SchemaVersion == system.SnapshotSchemaVersion {
			previous := m.snapshot
			m.snapshot = message.snapshot
			m.networkRates = CalculateNetworkRates(previous.Networks, message.snapshot.Networks, message.snapshot.Timestamp.Sub(previous.Timestamp))
			if value := message.snapshot.CPU.UtilizationPercent; value != nil {
				m.cpuHistory = m.cpuHistory.Add(*value)
			}
			if value := message.snapshot.Memory.UtilizationPercent; value != nil {
				m.memoryHistory = m.memoryHistory.Add(*value)
			}
			if message.request.IncludeFilesystems {
				m.lastFilesystem = message.snapshot.Timestamp
			}
			if message.request.IncludeDiagnostics {
				m.lastDiagnostics = message.snapshot.Timestamp
			}
		}
		return m, nil
	default:
		return m, nil
	}
}

func requestWindowSize() tea.Cmd {
	return func() tea.Msg {
		return tea.RequestWindowSize()
	}
}

func (m Model) View() tea.View {
	view := tea.NewView(Render(ViewState{
		Snapshot:      m.snapshot,
		Clock:         m.clock,
		CPUHistory:    m.cpuHistory,
		MemoryHistory: m.memoryHistory,
		NetworkRates:  m.networkRates,
		Width:         m.width,
		Height:        m.height,
		Help:          m.help,
		Paused:        m.paused,
		Collecting:    m.collecting,
		LastError:     m.lastError,
	}))
	view.AltScreen = true
	view.WindowTitle = "Setup Env"
	return view
}

func (m Model) handleKey(key string) (tea.Model, tea.Cmd) {
	switch KeyCommand(key) {
	case CommandQuit:
		return m, tea.Quit
	case CommandRefresh:
		if m.collecting {
			return m, nil
		}
		return m.startRefresh(true)
	case CommandHelp:
		m.help = !m.help
		return m, nil
	case CommandPause:
		m.paused = !m.paused
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) startRefresh(force bool) (Model, tea.Cmd) {
	if m.source == nil {
		m.lastError = context.Canceled
		return m, nil
	}
	m.collecting = true
	collectedAt := m.clock
	request := RefreshRequest{
		Previous:           m.snapshot,
		CollectedAt:        collectedAt,
		IncludeFilesystems: force || collectedAt.Sub(m.lastFilesystem) >= m.filesystemInterval,
		IncludeDiagnostics: force || collectedAt.Sub(m.lastDiagnostics) >= m.diagnosticInterval,
	}
	return m, func() tea.Msg {
		snapshot, err := m.source.Refresh(m.ctx, request)
		return refreshMsg{snapshot: snapshot, request: request, err: err}
	}
}

func (m Model) nextTick() tea.Cmd {
	return m.tick(m.refreshInterval, func(at time.Time) tea.Msg {
		return tickMsg{at: at}
	})
}

type Command int

const (
	CommandNone Command = iota
	CommandQuit
	CommandRefresh
	CommandHelp
	CommandPause
)

func KeyCommand(key string) Command {
	switch key {
	case "q", "ctrl+c":
		return CommandQuit
	case "r":
		return CommandRefresh
	case "?":
		return CommandHelp
	case "p", "space", " ":
		return CommandPause
	default:
		return CommandNone
	}
}
