package status

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/setup-env/app/internal/diagnostics"
	"github.com/setup-env/app/internal/system"
)

func WriteHuman(writer io.Writer, snapshot system.Snapshot) error {
	if _, err := io.WriteString(writer, "SETUP ENV SYSTEM STATUS\n\n"); err != nil {
		return err
	}
	if err := writeRows(writer, "Date and time", [][2]string{
		{"Local time", Timestamp(snapshot.Timestamp)},
		{"Time zone", Text(snapshot.TimeZone.Name) + " (" + UTCOffset(snapshot.TimeZone.UTCOffsetSeconds) + ")"},
	}); err != nil {
		return err
	}
	osName := snapshot.OperatingSystem.DisplayName
	if osName == "" {
		osName = snapshot.OperatingSystem.OS
	}
	if err := writeRows(writer, "System", [][2]string{
		{"Hostname", Text(snapshot.Host.Hostname)},
		{"Operating system", Text(osName)},
		{"OS version", Text(snapshot.OperatingSystem.Version)},
		{"Kernel version", Text(snapshot.OperatingSystem.KernelVersion)},
		{"Architecture", Text(snapshot.OperatingSystem.Architecture)},
		{"Current user", Text(snapshot.User.Username)},
		{"Uptime", Duration(snapshot.Host.UptimeSeconds)},
	}); err != nil {
		return err
	}
	if err := writeRows(writer, "CPU", [][2]string{
		{"Model", Text(snapshot.CPU.Model)},
		{"Physical cores", optionalInt(snapshot.CPU.PhysicalCores)},
		{"Logical CPUs", optionalInt(snapshot.CPU.LogicalCPUs)},
		{"Utilization", Percentage(snapshot.CPU.UtilizationPercent)},
	}); err != nil {
		return err
	}
	if err := writeRows(writer, "Memory", [][2]string{
		{"Used", Bytes(snapshot.Memory.UsedBytes)},
		{"Available", Bytes(snapshot.Memory.AvailableBytes)},
		{"Total", Bytes(snapshot.Memory.TotalBytes)},
		{"Utilization", Percentage(snapshot.Memory.UtilizationPercent)},
	}); err != nil {
		return err
	}
	if err := writeFilesystems(writer, snapshot.Filesystems); err != nil {
		return err
	}
	if err := writeNetworks(writer, snapshot.Networks); err != nil {
		return err
	}
	if err := writeRows(writer, "Development", [][2]string{
		{"Setup Env", snapshot.Development.ApplicationVersion.Version},
		{"Root", Text(snapshot.Development.Root.Path)},
		{"Root exists", YesNo(snapshot.Development.Root.Exists)},
		{"Root writable", YesNo(snapshot.Development.Root.Writable)},
		{"Directory context", string(snapshot.Development.Directory.Type)},
		{"Git", tool(snapshot.Development.Git)},
		{"GitHub CLI", tool(snapshot.Development.GitHubCLI)},
		{"GitHub authenticated", OptionalBool(snapshot.Development.GitHubAuthenticated)},
	}); err != nil {
		return err
	}
	if err := writeRows(writer, "Health", [][2]string{
		{"Status", string(snapshot.Diagnostics.Health)},
		{"Warnings", fmt.Sprintf("%d", snapshot.Diagnostics.WarningCount)},
		{"Failures", fmt.Sprintf("%d", snapshot.Diagnostics.FailureCount)},
	}); err != nil {
		return err
	}
	return writeWarnings(writer, snapshot)
}

func writeRows(writer io.Writer, heading string, rows [][2]string) error {
	if _, err := fmt.Fprintln(writer, heading); err != nil {
		return err
	}
	table := tabwriter.NewWriter(writer, 2, 4, 2, ' ', 0)
	for _, row := range rows {
		if _, err := fmt.Fprintf(table, "%s\t%s\n", row[0], row[1]); err != nil {
			return err
		}
	}
	if err := table.Flush(); err != nil {
		return err
	}
	_, err := fmt.Fprintln(writer)
	return err
}

func writeFilesystems(writer io.Writer, filesystems []system.Filesystem) error {
	if _, err := fmt.Fprintln(writer, "Filesystems"); err != nil {
		return err
	}
	if len(filesystems) == 0 {
		if _, err := fmt.Fprintln(writer, "  unavailable"); err != nil {
			return err
		}
		_, err := fmt.Fprintln(writer)
		return err
	}
	table := tabwriter.NewWriter(writer, 2, 4, 2, ' ', 0)
	for _, filesystem := range filesystems {
		usage := fmt.Sprintf("%s / %s  %s", Bytes(filesystem.UsedBytes), Bytes(filesystem.TotalBytes), Percentage(filesystem.UtilizationPercent))
		if _, err := fmt.Fprintf(table, "%s\t%s\t%s\n", filesystem.Mountpoint, usage, Text(filesystem.Type)); err != nil {
			return err
		}
	}
	if err := table.Flush(); err != nil {
		return err
	}
	_, err := fmt.Fprintln(writer)
	return err
}

func writeNetworks(writer io.Writer, networks []system.NetworkInterface) error {
	if _, err := fmt.Fprintln(writer, "Network"); err != nil {
		return err
	}
	if len(networks) == 0 {
		if _, err := fmt.Fprintln(writer, "  unavailable"); err != nil {
			return err
		}
		_, err := fmt.Fprintln(writer)
		return err
	}
	table := tabwriter.NewWriter(writer, 2, 4, 2, ' ', 0)
	for _, network := range networks {
		addresses := make([]string, 0, len(network.Addresses))
		for _, address := range network.Addresses {
			addresses = append(addresses, address.Address)
		}
		if len(addresses) == 0 {
			addresses = append(addresses, unavailable)
		}
		if _, err := fmt.Fprintf(table, "%s\t%s\t%s\n", network.Name, strings.Join(addresses, ", "), network.Status); err != nil {
			return err
		}
	}
	if err := table.Flush(); err != nil {
		return err
	}
	_, err := fmt.Fprintln(writer)
	return err
}

func writeWarnings(writer io.Writer, snapshot system.Snapshot) error {
	var messages []string
	for _, warning := range snapshot.Warnings {
		messages = append(messages, warning.Section+": "+warning.Message)
	}
	for _, check := range snapshot.Diagnostics.Checks {
		if check.Status != diagnostics.StatusPass {
			messages = append(messages, check.Name+": "+check.Message)
		}
	}
	if len(messages) == 0 {
		return nil
	}
	if _, err := fmt.Fprintln(writer, "Warnings"); err != nil {
		return err
	}
	for _, message := range messages {
		if _, err := fmt.Fprintln(writer, "  - "+message); err != nil {
			return err
		}
	}
	return nil
}

func optionalInt(value *int) string {
	if value == nil {
		return unavailable
	}
	return fmt.Sprintf("%d", *value)
}

func tool(value diagnostics.ToolStatus) string {
	if !value.Available {
		return "unavailable"
	}
	if value.Version == "" {
		return "available"
	}
	return value.Version
}
