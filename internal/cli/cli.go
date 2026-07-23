package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/setup-env/app/internal/app"
	"github.com/setup-env/app/internal/catalog"
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
  module    Inspect and validate Setup Env modules
  help      Show this command overview

Options for info and doctor:
  --json    Emit machine-readable JSON

Module discovery and manifest validation are implemented. Downloading,
installation, updates, and workflow execution are not implemented.
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
	case "module":
		return runModule(ctx, args[1:], stdout, service)
	default:
		return fmt.Errorf("unknown command %q; run 'setup-env help' for available commands", args[0])
	}
}

func runModule(ctx context.Context, args []string, stdout io.Writer, service app.Service) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		_, err := io.WriteString(stdout, moduleHelpText)
		return err
	}
	switch args[0] {
	case "list":
		options, jsonOutput, err := parseModuleListArguments(args[1:])
		if err != nil {
			return fmt.Errorf("module list: %w", err)
		}
		entries, err := service.ModuleList(ctx, options)
		if err != nil {
			return fmt.Errorf("load official module catalog: %w", err)
		}
		if jsonOutput {
			return writeJSON(stdout, struct {
				Modules []catalog.Entry `json:"modules"`
			}{Modules: entries})
		}
		return writeModuleList(stdout, entries)
	case "info":
		id, jsonOutput, err := parseModuleTargetArguments(args[1:], "module id")
		if err != nil {
			return fmt.Errorf("module info: %w", err)
		}
		info, err := service.ModuleInfo(ctx, id)
		if err != nil {
			return fmt.Errorf("inspect module %q: %w", id, err)
		}
		if jsonOutput {
			return writeJSON(stdout, info)
		}
		return writeModuleInfo(stdout, info)
	case "validate":
		target, jsonOutput, err := parseModuleTargetArguments(args[1:], "manifest path or module directory")
		if err != nil {
			return fmt.Errorf("module validate: %w", err)
		}
		report := service.ValidateManifest(target)
		if jsonOutput {
			if err := writeJSON(stdout, report); err != nil {
				return err
			}
		} else if err := writeManifestValidation(stdout, report); err != nil {
			return err
		}
		if !report.Valid {
			return fmt.Errorf("manifest validation failed for %q", report.Path)
		}
		return nil
	case "validate-catalog":
		jsonOutput, err := parseOutputArguments(args[1:])
		if err != nil {
			return fmt.Errorf("module validate-catalog: %w", err)
		}
		report := service.ValidateCatalog(ctx)
		if jsonOutput {
			if err := writeJSON(stdout, report); err != nil {
				return err
			}
		} else if err := writeCatalogValidation(stdout, report); err != nil {
			return err
		}
		if !report.Valid {
			return fmt.Errorf("official catalog validation failed")
		}
		return nil
	default:
		return fmt.Errorf("unknown module command %q; run 'setup-env module help'", args[0])
	}
}

const moduleHelpText = `Usage:
  setup-env module list [--json] [--trust <level>] [--status <status>] [--category <category>]
  setup-env module info <module> [--json]
  setup-env module validate <path> [--json]
  setup-env module validate-catalog [--json]

The catalog is embedded and local-only. These commands do not download,
install, update, or execute modules.
`

func parseModuleListArguments(args []string) (app.ModuleListOptions, bool, error) {
	var options app.ModuleListOptions
	jsonOutput := false
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--json":
			jsonOutput = true
		case "--trust":
			index++
			if index >= len(args) {
				return options, false, fmt.Errorf("--trust requires a value")
			}
			options.Trust = catalog.Trust(args[index])
			if options.Trust != catalog.TrustOfficial && options.Trust != catalog.TrustVerified && options.Trust != catalog.TrustCommunity {
				return options, false, fmt.Errorf("--trust must be official, verified, or community")
			}
		case "--status":
			index++
			if index >= len(args) {
				return options, false, fmt.Errorf("--status requires a value")
			}
			options.Status = catalog.Status(args[index])
			switch options.Status {
			case catalog.StatusActive, catalog.StatusPlanned, catalog.StatusExperimental, catalog.StatusDeprecated, catalog.StatusUnavailable:
			default:
				return options, false, fmt.Errorf("--status is not recognized")
			}
		case "--category":
			index++
			if index >= len(args) || args[index] == "" {
				return options, false, fmt.Errorf("--category requires a value")
			}
			options.Category = args[index]
		default:
			return options, false, fmt.Errorf("unsupported option %q", args[index])
		}
	}
	return options, jsonOutput, nil
}

func parseModuleTargetArguments(args []string, targetName string) (string, bool, error) {
	if len(args) == 0 {
		return "", false, fmt.Errorf("%s is required", targetName)
	}
	target := args[0]
	if target == "--json" {
		return "", false, fmt.Errorf("%s is required before --json", targetName)
	}
	jsonOutput, err := parseOutputArguments(args[1:])
	if err != nil {
		return "", false, err
	}
	return target, jsonOutput, nil
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

func writeModuleList(writer io.Writer, entries []catalog.Entry) error {
	table := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "ID\tNAME\tTRUST\tSTATUS\tREPOSITORY\tDESCRIPTION"); err != nil {
		return err
	}
	for _, entry := range entries {
		if _, err := fmt.Fprintf(table, "%s\t%s\t%s\t%s\t%s\t%s\n", entry.ID, entry.Name, entry.Trust, entry.Status, entry.Repository, entry.Description); err != nil {
			return err
		}
	}
	return table.Flush()
}

func writeModuleInfo(writer io.Writer, info app.ModuleInfo) error {
	entry := info.CatalogEntry
	rows := [][2]string{
		{"Module", entry.ID},
		{"Name", entry.Name},
		{"Description", entry.Description},
		{"Repository", entry.Repository},
		{"Manifest", entry.Manifest},
		{"Trust", string(entry.Trust)},
		{"Status", string(entry.Status)},
		{"Categories", strings.Join(entry.Categories, ", ")},
		{"Local manifest", fmt.Sprintf("%t", info.ManifestAvailable)},
		{"Compatibility", string(info.Compatibility.State)},
		{"Compatibility reason", info.Compatibility.Reason},
	}
	for _, row := range rows {
		if _, err := fmt.Fprintf(writer, "%-22s %s\n", row[0]+":", row[1]); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(writer, "Message:               "+info.Message); err != nil {
		return err
	}
	if info.Manifest != nil {
		if _, err := fmt.Fprintln(writer, "Workflows:"); err != nil {
			return err
		}
		for _, workflow := range info.Manifest.Workflows {
			if _, err := fmt.Fprintf(writer, "  - %s: %s (%s)\n", workflow.ID, workflow.Name, workflow.Entrypoint); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeManifestValidation(writer io.Writer, report app.ManifestValidationReport) error {
	if report.Valid {
		_, err := fmt.Fprintf(writer, "[PASS] %s is a valid setup-env.yaml manifest\nCompatibility: %s — %s\n", report.Path, report.Compatibility.State, report.Compatibility.Reason)
		return err
	}
	if _, err := fmt.Fprintf(writer, "[FAIL] %s is not a valid setup-env.yaml manifest\n", report.Path); err != nil {
		return err
	}
	for _, problem := range report.Problems {
		if _, err := fmt.Fprintln(writer, "  - "+problem); err != nil {
			return err
		}
	}
	return nil
}

func writeCatalogValidation(writer io.Writer, report app.CatalogValidationReport) error {
	if report.Valid {
		_, err := fmt.Fprintf(writer, "[PASS] embedded official catalog is valid (%d modules)\n", report.ModuleCount)
		return err
	}
	if _, err := io.WriteString(writer, "[FAIL] embedded official catalog is invalid\n"); err != nil {
		return err
	}
	for _, problem := range report.Problems {
		if _, err := fmt.Fprintln(writer, "  - "+problem); err != nil {
			return err
		}
	}
	return nil
}

func valueOrNone(value string) string {
	if value == "" {
		return "none"
	}
	return value
}
