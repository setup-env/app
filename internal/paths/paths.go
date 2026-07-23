package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

type Resolver struct {
	UserHomeDir func() (string, error)
}

func DefaultResolver() Resolver {
	return Resolver{UserHomeDir: os.UserHomeDir}
}

func (r Resolver) Home() (string, error) {
	home, err := r.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve current user's home directory: %w", err)
	}
	if home == "" {
		return "", fmt.Errorf("resolve current user's home directory: empty path")
	}
	return filepath.Clean(home), nil
}

func (r Resolver) DevelopmentRoot(override string) (string, error) {
	if override != "" {
		root, err := filepath.Abs(override)
		if err != nil {
			return "", fmt.Errorf("resolve development root override: %w", err)
		}
		return filepath.Clean(root), nil
	}
	home, err := r.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "dev"), nil
}
