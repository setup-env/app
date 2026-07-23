package git

import "testing"

func TestParseRemote(t *testing.T) {
	tests := []struct {
		name   string
		remote string
		host   string
		owner  string
		repo   string
	}{
		{name: "https", remote: "https://github.com/setup-env/app.git", host: "github.com", owner: "setup-env", repo: "app"},
		{name: "ssh", remote: "git@github.com:setup-env/app.git", host: "github.com", owner: "setup-env", repo: "app"},
		{name: "ssh URL", remote: "ssh://git@github.com/setup-env/app.git", host: "github.com", owner: "setup-env", repo: "app"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			host, owner, repo := ParseRemote(test.remote)
			if host != test.host || owner != test.owner || repo != test.repo {
				t.Fatalf("ParseRemote() = %q, %q, %q", host, owner, repo)
			}
		})
	}
}
