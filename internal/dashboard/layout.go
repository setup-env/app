package dashboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/setup-env/app/internal/diagnostics"
	statusformat "github.com/setup-env/app/internal/status"
	"github.com/setup-env/app/internal/system"
)

const (
	MinimumWidth  = 40
	MinimumHeight = 12
	WideWidth     = 100
)

type ViewState struct {
	Snapshot      system.Snapshot
	Clock         time.Time
	CPUHistory    History
	MemoryHistory History
	NetworkRates  []NetworkRate
	Width         int
	Height        int
	Help          bool
	Paused        bool
	Collecting    bool
	LastError     error
}

func Render(state ViewState) string {
	if state.Width < MinimumWidth || state.Height < MinimumHeight {
		return renderTooSmall(state.Width, state.Height)
	}
	if state.Help {
		return renderHelp(state.Width, state.Height)
	}
	if state.Width >= WideWidth {
		return renderWide(state)
	}
	return renderCompact(state)
}

func renderWide(state ViewState) string {
	snapshot := state.Snapshot
	header := fmt.Sprintf(
		"%s | %s | %s | uptime %s | %s",
		statusformat.Text(snapshot.Host.Hostname),
		osSummary(snapshot.OperatingSystem),
		statusformat.Text(snapshot.OperatingSystem.Architecture),
		statusformat.Duration(effectiveUptime(snapshot, state.Clock)),
		statusformat.Timestamp(state.Clock),
	)
	cpuWidth := state.Width / 2
	memoryWidth := state.Width - cpuWidth - 1
	cpu := panel("CPU", []string{
		fmt.Sprintf("%s  physical %s  logical %s",
			statusformat.Percentage(snapshot.CPU.UtilizationPercent),
			optionalInt(snapshot.CPU.PhysicalCores),
			optionalInt(snapshot.CPU.LogicalCPUs)),
		Truncate(snapshot.CPU.Model, cpuWidth-4),
		Sparkline(state.CPUHistory.Values(), cpuWidth-4),
	}, cpuWidth)
	memory := panel("Memory", []string{
		fmt.Sprintf("%s / %s  %s",
			statusformat.Bytes(snapshot.Memory.UsedBytes),
			statusformat.Bytes(snapshot.Memory.TotalBytes),
			statusformat.Percentage(snapshot.Memory.UtilizationPercent)),
		fmt.Sprintf("available %s", statusformat.Bytes(snapshot.Memory.AvailableBytes)),
		Sparkline(state.MemoryHistory.Values(), memoryWidth-4),
	}, memoryWidth)

	sections := []string{
		panel("Setup Env", []string{Truncate(header, state.Width-4)}, state.Width),
		joinColumns(cpu, memory, 1),
		panel("Filesystems", filesystemLines(snapshot.Filesystems, state.Width-4, availableRows(state.Height, 2)), state.Width),
		panel("Network", networkLines(snapshot.Networks, state.NetworkRates, state.Width-4, availableRows(state.Height, 2)), state.Width),
		panel("Development and health", developmentLines(snapshot, state.Width-4), state.Width),
		footer(state),
	}
	return fitHeight(strings.Join(sections, "\n"), state.Width, state.Height)
}

func renderCompact(state ViewState) string {
	snapshot := state.Snapshot
	lines := []string{
		fmt.Sprintf("%s | %s", statusformat.Text(snapshot.Host.Hostname), statusformat.Timestamp(state.Clock)),
		fmt.Sprintf("%s | uptime %s", osSummary(snapshot.OperatingSystem), statusformat.Duration(effectiveUptime(snapshot, state.Clock))),
		fmt.Sprintf("CPU %s %s", statusformat.Percentage(snapshot.CPU.UtilizationPercent), Sparkline(state.CPUHistory.Values(), max(8, state.Width-18))),
		fmt.Sprintf("MEM %s %s/%s", statusformat.Percentage(snapshot.Memory.UtilizationPercent), statusformat.Bytes(snapshot.Memory.UsedBytes), statusformat.Bytes(snapshot.Memory.TotalBytes)),
	}
	if len(snapshot.Filesystems) > 0 {
		filesystem := snapshot.Filesystems[0]
		lines = append(lines, fmt.Sprintf("DISK %s %s %s",
			Truncate(filesystem.Mountpoint, 12),
			statusformat.Percentage(filesystem.UtilizationPercent),
			UsageBar(filesystem.UtilizationPercent, max(8, state.Width-28))))
	} else {
		lines = append(lines, "DISK unavailable")
	}
	if len(snapshot.Networks) > 0 {
		network := snapshot.Networks[0]
		rate := findRate(state.NetworkRates, network.Name)
		lines = append(lines, fmt.Sprintf("NET  %s down %s up %s",
			Truncate(network.Name, 12),
			ByteRate(rate.BytesReceivedPerSec),
			ByteRate(rate.BytesSentPerSec)))
	} else {
		lines = append(lines, "NET  unavailable")
	}
	lines = append(lines,
		fmt.Sprintf("DEV  %s | Git %s | gh %s",
			Truncate(snapshot.Development.Root.Path, max(8, state.Width-31)),
			toolSummary(snapshot.Development.Git),
			toolSummary(snapshot.Development.GitHubCLI)),
		fmt.Sprintf("HEALTH %s | warnings %d | failures %d",
			snapshot.Diagnostics.Health,
			snapshot.Diagnostics.WarningCount,
			snapshot.Diagnostics.FailureCount),
	)
	if warning := firstWarning(snapshot); warning != "" {
		lines = append(lines, "WARN "+Truncate(warning, state.Width-7))
	}
	content := panel("Setup Env", lines, state.Width) + "\n" + footer(state)
	return fitHeight(content, state.Width, state.Height)
}

func renderHelp(width, height int) string {
	lines := []string{
		"q / Ctrl+C  quit and restore the terminal",
		"r           refresh all metrics now",
		"p / Space   pause or resume metric refresh",
		"?           return to the dashboard",
		"",
		"CPU, memory, and network refresh every second.",
		"Filesystems refresh every five seconds.",
		"Development diagnostics refresh every sixty seconds.",
		"Unavailable metrics remain visible and do not stop the dashboard.",
		"No process, public-IP, credential, or Wi-Fi secret data is collected.",
	}
	return fitHeight(panel("Setup Env help", lines, width), width, height)
}

func renderTooSmall(width, height int) string {
	return fmt.Sprintf(
		"Setup Env dashboard\n\nTerminal too small (%dx%d).\nResize to at least %dx%d.\n\nq quit | ? help\n",
		width, height, MinimumWidth, MinimumHeight,
	)
}

func panel(title string, lines []string, width int) string {
	if width < 4 {
		return ""
	}
	inside := width - 2
	label := " " + title + " "
	if len(label) > inside {
		label = Truncate(label, inside)
	}
	top := "+" + label + strings.Repeat("-", max(0, inside-len(label))) + "+"
	var result strings.Builder
	result.WriteString(top)
	for _, line := range lines {
		result.WriteByte('\n')
		line = Truncate(line, inside-2)
		result.WriteString("| ")
		result.WriteString(line)
		result.WriteString(strings.Repeat(" ", max(0, inside-2-len([]rune(line)))))
		result.WriteString(" |")
	}
	result.WriteByte('\n')
	result.WriteString("+")
	result.WriteString(strings.Repeat("-", inside))
	result.WriteString("+")
	return result.String()
}

func joinColumns(left, right string, gap int) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")
	count := max(len(leftLines), len(rightLines))
	leftWidth := 0
	for _, line := range leftLines {
		leftWidth = max(leftWidth, len([]rune(line)))
	}
	var result strings.Builder
	for index := 0; index < count; index++ {
		if index > 0 {
			result.WriteByte('\n')
		}
		leftLine := ""
		if index < len(leftLines) {
			leftLine = leftLines[index]
		}
		rightLine := ""
		if index < len(rightLines) {
			rightLine = rightLines[index]
		}
		result.WriteString(leftLine)
		result.WriteString(strings.Repeat(" ", leftWidth-len([]rune(leftLine))+gap))
		result.WriteString(rightLine)
	}
	return result.String()
}

func filesystemLines(filesystems []system.Filesystem, width, limit int) []string {
	if len(filesystems) == 0 {
		return []string{"unavailable"}
	}
	limit = max(1, min(limit, len(filesystems)))
	lines := make([]string, 0, limit+1)
	for _, filesystem := range filesystems[:limit] {
		barWidth := max(8, min(24, width/4))
		valueWidth := max(8, width-barWidth-34)
		lines = append(lines, fmt.Sprintf("%-*s %9s / %-9s %6s %s",
			valueWidth,
			Truncate(filesystem.Mountpoint, valueWidth),
			statusformat.Bytes(filesystem.UsedBytes),
			statusformat.Bytes(filesystem.TotalBytes),
			statusformat.Percentage(filesystem.UtilizationPercent),
			UsageBar(filesystem.UtilizationPercent, barWidth)))
	}
	if hidden := len(filesystems) - limit; hidden > 0 {
		lines = append(lines, fmt.Sprintf("... %d more filesystem(s) hidden", hidden))
	}
	return lines
}

func networkLines(networks []system.NetworkInterface, rates []NetworkRate, width, limit int) []string {
	if len(networks) == 0 {
		return []string{"unavailable"}
	}
	limit = max(1, min(limit, len(networks)))
	lines := make([]string, 0, limit+1)
	for _, network := range networks[:limit] {
		rate := findRate(rates, network.Name)
		addresses := make([]string, 0, len(network.Addresses))
		for _, address := range network.Addresses {
			addresses = append(addresses, address.Address)
		}
		lines = append(lines, fmt.Sprintf("%-18s %-*s down %-12s up %-12s",
			Truncate(network.Name, 18),
			max(8, width-65),
			Truncate(strings.Join(addresses, ", "), max(8, width-65)),
			ByteRate(rate.BytesReceivedPerSec),
			ByteRate(rate.BytesSentPerSec)))
	}
	if hidden := len(networks) - limit; hidden > 0 {
		lines = append(lines, fmt.Sprintf("... %d more interface(s) hidden", hidden))
	}
	return lines
}

func developmentLines(snapshot system.Snapshot, width int) []string {
	lines := []string{
		fmt.Sprintf("root %s | writable %s | context %s",
			Truncate(snapshot.Development.Root.Path, max(8, width-43)),
			statusformat.YesNo(snapshot.Development.Root.Writable),
			snapshot.Development.Directory.Type),
		fmt.Sprintf("Git %s | GitHub CLI %s | authenticated %s",
			toolSummary(snapshot.Development.Git),
			toolSummary(snapshot.Development.GitHubCLI),
			statusformat.OptionalBool(snapshot.Development.GitHubAuthenticated)),
		fmt.Sprintf("health %s | warnings %d | failures %d",
			snapshot.Diagnostics.Health,
			snapshot.Diagnostics.WarningCount,
			snapshot.Diagnostics.FailureCount),
	}
	if warning := firstWarning(snapshot); warning != "" {
		lines = append(lines, "warning: "+Truncate(warning, width-9))
	}
	return lines
}

func footer(state ViewState) string {
	mode := "live"
	if state.Paused {
		mode = "paused"
	} else if state.Collecting {
		mode = "refreshing"
	}
	message := fmt.Sprintf(" q quit | r refresh | p pause | ? help | %s ", mode)
	if state.LastError != nil {
		message += "| last refresh: " + Truncate(state.LastError.Error(), max(8, state.Width-len(message)-3)) + " "
	}
	return Truncate(message, state.Width)
}

func fitHeight(content string, width, height int) string {
	lines := strings.Split(content, "\n")
	if len(lines) > height {
		lines = lines[:height]
		if height > 0 {
			lines[height-1] = Truncate("... terminal height hides additional dashboard rows", width)
		}
	}
	return strings.Join(lines, "\n")
}

func effectiveUptime(snapshot system.Snapshot, now time.Time) *uint64 {
	if snapshot.Host.UptimeSeconds == nil {
		return nil
	}
	seconds := *snapshot.Host.UptimeSeconds
	if delta := now.Sub(snapshot.Timestamp); delta > 0 {
		seconds += uint64(delta / time.Second)
	}
	return &seconds
}

func osSummary(value system.OperatingSystem) string {
	name := value.DisplayName
	if name == "" {
		name = value.OS
	}
	if value.Version != "" && !strings.Contains(name, value.Version) {
		name += " " + value.Version
	}
	return statusformat.Text(name)
}

func optionalInt(value *int) string {
	if value == nil {
		return "unavailable"
	}
	return fmt.Sprintf("%d", *value)
}

func toolSummary(value diagnostics.ToolStatus) string {
	if !value.Available {
		return "unavailable"
	}
	return "available"
}

func firstWarning(snapshot system.Snapshot) string {
	if len(snapshot.Warnings) > 0 {
		return snapshot.Warnings[0].Section + ": " + snapshot.Warnings[0].Message
	}
	for _, check := range snapshot.Diagnostics.Checks {
		if check.Status != diagnostics.StatusPass {
			return check.Name + ": " + check.Message
		}
	}
	return ""
}

func findRate(rates []NetworkRate, name string) NetworkRate {
	for _, rate := range rates {
		if strings.EqualFold(rate.Name, name) {
			return rate
		}
	}
	return NetworkRate{Name: name}
}

func availableRows(height, divisor int) int {
	return max(1, (height-22)/divisor)
}
