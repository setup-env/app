package git

import (
	"context"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Repository struct {
	Root       string `json:"root"`
	RemoteURL  string `json:"remote_url,omitempty"`
	RemoteHost string `json:"remote_host,omitempty"`
	Owner      string `json:"owner,omitempty"`
	Name       string `json:"name,omitempty"`
}

type Inspector interface {
	Repository(context.Context, string) (Repository, bool)
}

type CommandInspector struct {
	Timeout time.Duration
}

func DefaultInspector() CommandInspector {
	return CommandInspector{Timeout: 3 * time.Second}
}

func (i CommandInspector) Repository(ctx context.Context, directory string) (Repository, bool) {
	root, err := i.run(ctx, directory, "rev-parse", "--show-toplevel")
	if err != nil {
		return Repository{}, false
	}
	repository := Repository{Root: filepath.Clean(root)}
	if remote, err := i.run(ctx, directory, "remote", "get-url", "origin"); err == nil {
		repository.RemoteURL = remote
		repository.RemoteHost, repository.Owner, repository.Name = ParseRemote(remote)
	}
	return repository, true
}

func (i CommandInspector) run(parent context.Context, directory string, args ...string) (string, error) {
	timeout := i.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	command := exec.CommandContext(ctx, "git", args...)
	command.Dir = directory
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func ParseRemote(remote string) (host, owner, name string) {
	value := strings.TrimSuffix(strings.TrimSpace(remote), ".git")
	if value == "" {
		return "", "", ""
	}

	var path string
	if strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err != nil {
			return "", "", ""
		}
		host = parsed.Hostname()
		path = parsed.Path
	} else if at := strings.Index(value, "@"); at >= 0 {
		afterAt := value[at+1:]
		host, path, _ = strings.Cut(afterAt, ":")
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return host, parts[len(parts)-2], parts[len(parts)-1]
	}
	return host, "", ""
}
