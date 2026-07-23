package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/setup-env/app/internal/app"
)

func TestNoArgumentsDisplaysHelp(t *testing.T) {
	var output bytes.Buffer
	if err := run(context.Background(), nil, &output, &output, app.Service{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "setup-env <command>") {
		t.Fatalf("help output = %q", output.String())
	}
}

func TestVersionCommand(t *testing.T) {
	var output bytes.Buffer
	if err := run(context.Background(), []string{"version"}, &output, &output, app.Service{}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(output.String(), "setup-env ") {
		t.Fatalf("version output = %q", output.String())
	}
}

func TestUnknownCommandIsActionable(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), []string{"module"}, &output, &output, app.Service{})
	if err == nil || !strings.Contains(err.Error(), "setup-env help") {
		t.Fatalf("run() error = %v", err)
	}
}

func TestUnsupportedInfoOption(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), []string{"info", "--yaml"}, &output, &output, app.Service{})
	if err == nil || !strings.Contains(err.Error(), "--json") {
		t.Fatalf("run() error = %v", err)
	}
}
