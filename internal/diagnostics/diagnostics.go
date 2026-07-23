package diagnostics

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Status string

const (
	StatusPass Status = "pass"
	StatusWarn Status = "warn"
	StatusFail Status = "fail"
)

type Check struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Message string `json:"message"`
}

type Report struct {
	Ready  bool    `json:"ready"`
	Checks []Check `json:"checks"`
}

type CommandRunner interface {
	LookPath(string) (string, error)
	Run(context.Context, string, ...string) ([]byte, error)
}

type OSCommandRunner struct{}

func (OSCommandRunner) LookPath(name string) (string, error) {
	return exec.LookPath(name)
}

func (OSCommandRunner) Run(ctx context.Context, name string, args ...string) ([]byte, error) {
	commandCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	return exec.CommandContext(commandCtx, name, args...).CombinedOutput()
}

func Run(ctx context.Context, developmentRoot string, runner CommandRunner) Report {
	if runner == nil {
		runner = OSCommandRunner{}
	}
	checks := []Check{
		checkCommand(ctx, runner, "git", "Git", true),
		checkCommand(ctx, runner, "gh", "GitHub CLI", false),
		checkGitHubAuthentication(ctx, runner),
		checkDevelopmentRoot(developmentRoot),
	}
	ready := true
	for _, check := range checks {
		if check.Status == StatusFail {
			ready = false
		}
	}
	return Report{Ready: ready, Checks: checks}
}

func checkCommand(ctx context.Context, runner CommandRunner, executable, label string, required bool) Check {
	path, err := runner.LookPath(executable)
	if err != nil {
		status := StatusWarn
		if required {
			status = StatusFail
		}
		return Check{Name: executable, Status: status, Message: label + " is not available on PATH"}
	}
	output, err := runner.Run(ctx, path, "--version")
	if err != nil {
		return Check{Name: executable, Status: StatusWarn, Message: label + " is installed but its version could not be determined"}
	}
	version := strings.TrimSpace(string(output))
	if first, _, ok := strings.Cut(version, "\n"); ok {
		version = strings.TrimSpace(first)
	}
	return Check{Name: executable, Status: StatusPass, Message: version}
}

func checkGitHubAuthentication(ctx context.Context, runner CommandRunner) Check {
	path, err := runner.LookPath("gh")
	if err != nil {
		return Check{Name: "github-auth", Status: StatusWarn, Message: "GitHub CLI is unavailable; authenticated GitHub operations are not ready"}
	}
	if _, err := runner.Run(ctx, path, "auth", "status", "--hostname", "github.com"); err != nil {
		return Check{Name: "github-auth", Status: StatusWarn, Message: "GitHub CLI authentication is not configured or requires renewal; run 'gh auth login'"}
	}
	return Check{Name: "github-auth", Status: StatusPass, Message: "GitHub CLI authentication appears configured"}
}

func checkDevelopmentRoot(root string) Check {
	info, err := os.Stat(root)
	if errors.Is(err, os.ErrNotExist) {
		parent := filepath.Dir(root)
		if writeErr := probeWrite(parent); writeErr != nil {
			return Check{Name: "development-root", Status: StatusFail, Message: fmt.Sprintf("%s does not exist and its parent is not writable: %v", root, writeErr)}
		}
		return Check{Name: "development-root", Status: StatusWarn, Message: root + " does not exist yet; its parent is writable"}
	}
	if err != nil {
		return Check{Name: "development-root", Status: StatusFail, Message: fmt.Sprintf("cannot inspect %s: %v", root, err)}
	}
	if !info.IsDir() {
		return Check{Name: "development-root", Status: StatusFail, Message: root + " exists but is not a directory"}
	}
	if err := probeWrite(root); err != nil {
		return Check{Name: "development-root", Status: StatusFail, Message: fmt.Sprintf("%s is not writable: %v", root, err)}
	}
	return Check{Name: "development-root", Status: StatusPass, Message: root + " exists and is writable"}
}

func probeWrite(directory string) error {
	file, err := os.CreateTemp(directory, ".setup-env-write-check-*")
	if err != nil {
		return err
	}
	name := file.Name()
	closeErr := file.Close()
	removeErr := os.Remove(name)
	if closeErr != nil {
		return closeErr
	}
	return removeErr
}
