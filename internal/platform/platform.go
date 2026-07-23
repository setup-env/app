package platform

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"runtime"
	"strings"
)

type Info struct {
	OS           string `json:"os"`
	Distribution string `json:"distribution,omitempty"`
	Architecture string `json:"architecture"`
	User         string `json:"user"`
	Home         string `json:"home"`
	Shell        string `json:"shell,omitempty"`
}

type Detector struct {
	GOOS        string
	GOARCH      string
	CurrentUser func() (*user.User, error)
	UserHomeDir func() (string, error)
	LookupEnv   func(string) (string, bool)
	Open        func(string) (io.ReadCloser, error)
}

func DefaultDetector() Detector {
	return Detector{
		GOOS:        runtime.GOOS,
		GOARCH:      runtime.GOARCH,
		CurrentUser: user.Current,
		UserHomeDir: os.UserHomeDir,
		LookupEnv:   os.LookupEnv,
		Open: func(name string) (io.ReadCloser, error) {
			return os.Open(name)
		},
	}
}

func (d Detector) Detect() (Info, error) {
	current, err := d.CurrentUser()
	if err != nil {
		return Info{}, fmt.Errorf("detect current user: %w", err)
	}
	home, err := d.UserHomeDir()
	if err != nil {
		return Info{}, fmt.Errorf("detect home directory: %w", err)
	}

	info := Info{
		OS:           d.GOOS,
		Architecture: d.GOARCH,
		User:         current.Username,
		Home:         home,
		Shell:        d.shell(),
	}
	if d.GOOS == "linux" {
		info.Distribution = d.linuxDistribution()
	}
	return info, nil
}

func (d Detector) shell() string {
	for _, name := range []string{"SHELL", "COMSPEC"} {
		if value, ok := d.LookupEnv(name); ok && value != "" {
			return value
		}
	}
	return ""
}

func (d Detector) linuxDistribution() string {
	file, err := d.Open("/etc/os-release")
	if err != nil {
		return ""
	}
	defer file.Close()
	values := parseOSRelease(file)
	if pretty := values["PRETTY_NAME"]; pretty != "" {
		return pretty
	}
	return values["ID"]
}

func parseOSRelease(reader io.Reader) map[string]string {
	values := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		values[key] = strings.Trim(value, `"'`)
	}
	return values
}
