package manifest

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func Parse(data []byte) (Manifest, error) {
	var result Manifest
	decoder := yaml.NewDecoder(bytes.NewReader(data), yaml.DisallowUnknownField())
	if err := decoder.Decode(&result); err != nil {
		return Manifest{}, fmt.Errorf("parse setup-env.yaml: %w", err)
	}
	return result, nil
}

func ParseFile(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest %q: %w", path, err)
	}
	result, err := Parse(data)
	if err != nil {
		return Manifest{}, fmt.Errorf("parse manifest %q: %w", path, err)
	}
	return result, nil
}
