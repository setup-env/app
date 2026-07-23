package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/setup-env/app/internal/app"
	"github.com/setup-env/app/internal/diagnostics"
	"github.com/setup-env/app/internal/version"
)

const helpText = `Setup Env is a universal setup and scaffolding platform.

Usage:
  setup-env <command> [options]

Commands:
  version   Print application version information
  info      Report platform and directory context
  doctor    Check local readiness for Setup Env operations
  help      Show this command overview

Options for info and doctor:
  --json    Emit machine-readable JSON

Module, workflow, and run commands are planned but are not implemented yet.
`

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	return run(ctx, args, stdout, stderr, app.DefaultService())
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer, service app.Service) error {
	_ = stderr
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		_, err := io.WriteString(stdout, helpText)
		return err
	}

	switch args[0] {
	case "version":
		if len(args) != 1 {
			return fmt.Errorf("version does not accept arguments; run 'setup-env help'")
		}
		current := version.Current()
		_, err := fmt.Fprintf(stdout, "setup-env %s (commit %s, built %s, %s)\n", current.Version, current.Commit, current.BuildDate, current.GoVersion)
		return err
	case "info":
		jsonOutput, err := parseOutputArguments(args[1:])
		if err != nil {
			return fmt.Errorf("info: %w", err)
		}
		info, err := service.Info(ctx)
		if err != nil {
			return fmt.Errorf("collect environment information: %w", err)
		}
		if jsonOutput {
			return writeJSON(stdout, info)
		}
		return writeInfo(stdout, info)
	case "doctor":
		jsonOutput, err := parseOutputArguments(args[1:])
		if err != nil {
			return fmt.Errorf("doctor: %w", err)
		}
		report, err := service.Doctor(ctx)
		if err != nil {
			return fmt.Errorf("run diagnostics: %w", err)
		}
		if jsonOutput {
			return writeJSON(stdout, report)
		}
		return writeDoctor(stdout, report)
	default:
		return fmt.Errorf("unknown command %q; run 'setup-env help' for available commands", args[0])
	}
}

func parseOutputArguments(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	if len(args) == 1 && args[0] == "--json" {
		return true, nil
	}
	return false, fmt.Errorf("unsupported option %q; the only supported option is --json", strings.Join(args, " "))
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func writeInfo(writer io.Writer, info app.Info) error {
	rows := [][2]string{
		{"Operating system", info.Platform.OS},
		{"Distribution", valueOrNone(info.Platform.Distribution)},
		{"Architecture", info.Platform.Architecture},
		{"User", info.Platform.User},
		{"Home", info.Platform.Home},
		{"Shell", valueOrNone(info.Platform.Shell)},
		{"Current directory", info.Directory.CurrentDirectory},
		{"Development root", info.Directory.DevelopmentRoot},
		{"Directory type", string(info.Directory.Type)},
		{"Organization", valueOrNone(info.Directory.Organization)},
		{"Repository", valueOrNone(info.Directory.Repository)},
		{"Git repository", fmt.Sprintf("%t", info.Directory.IsGitRepository)},
		{"Remote organization", valueOrNone(info.Directory.RemoteOrganization)},
		{"Remote repository", valueOrNone(info.Directory.RemoteRepository)},
		{"Configuration", info.ConfigPath},
		{"Configuration loaded", fmt.Sprintf("%t", info.ConfigLoaded)},
		{"Clone protocol", string(info.CloneProtocol)},
	}
	for _, row := range rows {
		if _, err := fmt.Fprintf(writer, "%-22s %s\n", row[0]+":", row[1]); err != nil {
			return err
		}
	}
	return nil
}

func writeDoctor(writer io.Writer, report diagnostics.Report) error {
	for _, check := range report.Checks {
		if _, err := fmt.Fprintf(writer, "[%-4s] %-18s %s\n", strings.ToUpper(string(check.Status)), check.Name, check.Message); err != nil {
			return err
		}
	}
	if report.Ready {
		_, err := io.WriteString(writer, "\nSetup Env is ready for core local operations.\n")
		return err
	}
	_, err := io.WriteString(writer, "\nSetup Env needs attention before core local operations are ready.\n")
	return err
}

func valueOrNone(value string) string {
	if value == "" {
		return "none"
	}
	return value
}
