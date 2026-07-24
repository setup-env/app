package system

import (
	"context"
	"io"
	"os/user"
	"strings"
	"testing"

	"github.com/setup-env/app/internal/platform"
	"github.com/shirou/gopsutil/v4/host"
)

type fakeHostSource struct {
	info *host.InfoStat
}

func (f fakeHostSource) Info(context.Context) (*host.InfoStat, error) {
	return f.info, nil
}

func TestHostCollectorCombinesPlatformAndHostInformation(t *testing.T) {
	detector := platform.Detector{
		GOOS:   "linux",
		GOARCH: "arm64",
		CurrentUser: func() (*user.User, error) {
			return &user.User{Username: "tester"}, nil
		},
		UserHomeDir: func() (string, error) { return "/home/tester", nil },
		LookupEnv:   func(name string) (string, bool) { return "/bin/zsh", name == "SHELL" },
		Open: func(string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(`PRETTY_NAME="Example Linux 1"`)), nil
		},
	}
	snapshot := Snapshot{}
	err := (HostCollector{
		Source: fakeHostSource{info: &host.InfoStat{
			Hostname:        "example",
			Uptime:          3600,
			OS:              "linux",
			Platform:        "example",
			PlatformVersion: "1",
			KernelVersion:   "6.0",
		}},
		PlatformDetector: detector,
	}).Collect(context.Background(), &snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Host.Hostname != "example" || *snapshot.Host.UptimeSeconds != 3600 {
		t.Fatalf("host = %#v", snapshot.Host)
	}
	if snapshot.OperatingSystem.DisplayName != "Example Linux 1" ||
		snapshot.OperatingSystem.Architecture != "arm64" ||
		snapshot.OperatingSystem.KernelVersion != "6.0" {
		t.Fatalf("operating system = %#v", snapshot.OperatingSystem)
	}
	if snapshot.User.Username != "tester" ||
		snapshot.User.Home != "/home/tester" ||
		snapshot.User.Shell != "/bin/zsh" {
		t.Fatalf("user = %#v", snapshot.User)
	}
}
