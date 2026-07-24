package dashboard

import (
	"bytes"
	"testing"
)

func TestBuffersAreNotInteractiveTerminals(t *testing.T) {
	var input bytes.Buffer
	var output bytes.Buffer
	if IsInteractive(&input, &output) {
		t.Fatal("buffers reported as interactive terminals")
	}
}
