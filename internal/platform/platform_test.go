package platform

import (
	"io"
	"os/user"
	"strings"
	"testing"
)

func TestDetectLinuxPlatform(t *testing.T) {
	detector := Detector{
		GOOS:   "linux",
		GOARCH: "arm64",
		CurrentUser: func() (*user.User, error) {
			return &user.User{Username: "tester"}, nil
		},
		UserHomeDir: func() (string, error) { return "/home/tester", nil },
		LookupEnv: func(key string) (string, bool) {
			if key == "SHELL" {
				return "/bin/zsh", true
			}
			return "", false
		},
		Open: func(string) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("ID=ubuntu\nPRETTY_NAME=\"Ubuntu 26.04 LTS\"\n")), nil
		},
	}
	info, err := detector.Detect()
	if err != nil {
		t.Fatal(err)
	}
	if info.OS != "linux" || info.Architecture != "arm64" || info.Distribution != "Ubuntu 26.04 LTS" {
		t.Fatalf("Detect() = %#v", info)
	}
	if info.User != "tester" || info.Home != "/home/tester" || info.Shell != "/bin/zsh" {
		t.Fatalf("Detect() user fields = %#v", info)
	}
}

func TestDetectWindowsPlatformModel(t *testing.T) {
	detector := Detector{
		GOOS:   "windows",
		GOARCH: "amd64",
		CurrentUser: func() (*user.User, error) {
			return &user.User{Username: `DOMAIN\person`}, nil
		},
		UserHomeDir: func() (string, error) { return `C:\Users\person`, nil },
		LookupEnv: func(key string) (string, bool) {
			if key == "COMSPEC" {
				return `C:\Windows\System32\cmd.exe`, true
			}
			return "", false
		},
		Open: func(string) (io.ReadCloser, error) {
			t.Fatal("Windows detection must not read /etc/os-release")
			return nil, nil
		},
	}
	info, err := detector.Detect()
	if err != nil {
		t.Fatal(err)
	}
	if info.OS != "windows" || info.Distribution != "" || info.Home != `C:\Users\person` {
		t.Fatalf("Detect() = %#v", info)
	}
}
