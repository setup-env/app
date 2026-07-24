package release

import "testing"

func TestNormalizeVersion(t *testing.T) {
	for input, want := range map[string]string{
		"0.1.0":       "v0.1.0",
		"v1.2.3":      "v1.2.3",
		"v1.2.3-rc.1": "v1.2.3-rc.1",
	} {
		got, err := NormalizeVersion(input)
		if err != nil || got != want {
			t.Fatalf("NormalizeVersion(%q) = %q, %v; want %q", input, got, err, want)
		}
	}
	for _, input := range []string{"", "v1", "v1.2", "v01.2.3", "latest", "v1.2.3+build.4", "v1.2.3/evil"} {
		if _, err := NormalizeVersion(input); err == nil {
			t.Fatalf("NormalizeVersion(%q) error = nil", input)
		}
	}
}

func TestArtifactNames(t *testing.T) {
	tests := []struct {
		target Target
		want   string
	}{
		{Target{OS: "windows", Arch: "amd64"}, "setup-env_0.1.0_windows_amd64.zip"},
		{Target{OS: "darwin", Arch: "arm64"}, "setup-env_0.1.0_darwin_arm64.tar.gz"},
		{Target{OS: "linux", Arch: "amd64"}, "setup-env_0.1.0_linux_amd64.tar.gz"},
	}
	for _, test := range tests {
		if got := test.target.ArchiveName("v0.1.0"); got != test.want {
			t.Fatalf("%s ArchiveName() = %q, want %q", test.target, got, test.want)
		}
	}
}

func TestParseTargetsDeduplicatesAndRejectsUnknown(t *testing.T) {
	targets, err := ParseTargets("linux/arm64,windows/amd64,linux/arm64")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 2 || targets[0].String() != "linux/arm64" || targets[1].String() != "windows/amd64" {
		t.Fatalf("targets = %+v", targets)
	}
	if _, err := ParseTargets("freebsd/amd64"); err == nil {
		t.Fatal("ParseTargets() accepted unsupported target")
	}
}
