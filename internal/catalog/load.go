package catalog

import (
	"bytes"
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
)

func Parse(data []byte) (Catalog, error) {
	var result Catalog
	decoder := yaml.NewDecoder(bytes.NewReader(data), yaml.DisallowUnknownField())
	if err := decoder.Decode(&result); err != nil {
		return Catalog{}, fmt.Errorf("parse module catalog: %w", err)
	}
	return result, nil
}

func Load(ctx context.Context, source Source) (Catalog, error) {
	if source == nil {
		return Catalog{}, fmt.Errorf("load module catalog: source is not configured")
	}
	data, err := source.Read(ctx)
	if err != nil {
		return Catalog{}, err
	}
	result, err := Parse(data)
	if err != nil {
		return Catalog{}, err
	}
	if err := Validate(result); err != nil {
		return Catalog{}, err
	}
	return result, nil
}
