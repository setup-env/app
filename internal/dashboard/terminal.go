package dashboard

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
)

type fileDescriptor interface {
	Fd() uintptr
}

func IsInteractive(input io.Reader, output io.Writer) bool {
	in, inputOK := input.(fileDescriptor)
	out, outputOK := output.(fileDescriptor)
	if !inputOK || !outputOK {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("TERM")), "dumb") {
		return false
	}
	return term.IsTerminal(in.Fd()) && term.IsTerminal(out.Fd())
}
