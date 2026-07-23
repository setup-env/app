package catalog

import (
	"context"
	"fmt"
	"os"

	catalogdata "github.com/setup-env/app/catalog"
)

type Source interface {
	Read(context.Context) ([]byte, error)
}

type EmbeddedSource struct{}

func (EmbeddedSource) Read(context.Context) ([]byte, error) {
	return append([]byte(nil), catalogdata.Modules...), nil
}

type FileSource struct {
	Path string
}

func (s FileSource) Read(context.Context) ([]byte, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, fmt.Errorf("read catalog %q: %w", s.Path, err)
	}
	return data, nil
}

type BytesSource []byte

func (s BytesSource) Read(context.Context) ([]byte, error) {
	return append([]byte(nil), s...), nil
}
