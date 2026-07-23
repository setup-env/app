package compatibility

import "testing"

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name    string
		app     string
		minimum string
		state   State
	}{
		{name: "compatible equal", app: "1.2.3", minimum: "1.2.3", state: StateCompatible},
		{name: "compatible newer", app: "1.3.0", minimum: "1.2.3", state: StateCompatible},
		{name: "incompatible", app: "1.2.2", minimum: "1.2.3", state: StateIncompatible},
		{name: "missing minimum", app: "1.2.3", minimum: "", state: StateUnknown},
		{name: "malformed minimum", app: "1.2.3", minimum: "latest", state: StateUnknown},
		{name: "development app", app: "dev", minimum: "1.2.3", state: StateUnknown},
		{name: "prerelease below stable", app: "1.2.3-rc.1", minimum: "1.2.3", state: StateIncompatible},
		{name: "stable above prerelease", app: "1.2.3", minimum: "1.2.3-rc.1", state: StateCompatible},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Evaluate(test.app, test.minimum)
			if result.State != test.state {
				t.Fatalf("Evaluate(%q, %q).State = %q, want %q", test.app, test.minimum, result.State, test.state)
			}
		})
	}
}

func TestParseRejectsMalformedVersion(t *testing.T) {
	for _, value := range []string{"", "1", "1.2", "1.2.x", "1.2.3-01"} {
		if _, err := Parse(value); err == nil {
			t.Fatalf("Parse(%q) error = nil, want error", value)
		}
	}
}
