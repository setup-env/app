package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const SchemaVersion = 1

type CloneProtocol string

const (
	CloneProtocolHTTPS CloneProtocol = "https"
	CloneProtocolSSH   CloneProtocol = "ssh"
)

type Settings struct {
	CheckUpdates bool `json:"check_updates"`
}

type Config struct {
	SchemaVersion   int           `json:"schema_version"`
	DevelopmentRoot string        `json:"development_root,omitempty"`
	CloneProtocol   CloneProtocol `json:"clone_protocol"`
	Organizations   []string      `json:"known_organizations,omitempty"`
	Settings        Settings      `json:"settings"`
}

type LocationResolver struct {
	UserConfigDir func() (string, error)
}

func Default() Config {
	return Config{
		SchemaVersion: SchemaVersion,
		CloneProtocol: CloneProtocolHTTPS,
		Organizations: []string{},
		Settings: Settings{
			CheckUpdates: true,
		},
	}
}

func DefaultLocationResolver() LocationResolver {
	return LocationResolver{UserConfigDir: os.UserConfigDir}
}

func (r LocationResolver) Path() (string, error) {
	base, err := r.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve operating system configuration directory: %w", err)
	}
	if base == "" {
		return "", fmt.Errorf("resolve operating system configuration directory: empty path")
	}
	return filepath.Join(base, "setup-env", "config.json"), nil
}

func Load(path string) (Config, bool, error) {
	result := Default()
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return result, false, nil
	}
	if err != nil {
		return Config{}, false, fmt.Errorf("read configuration %q: %w", path, err)
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return Config{}, true, fmt.Errorf("parse configuration %q: %w", path, err)
	}
	if err := result.Validate(); err != nil {
		return Config{}, true, fmt.Errorf("validate configuration %q: %w", path, err)
	}
	return result, true, nil
}

func (c Config) Validate() error {
	if c.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema version %d; supported version is %d", c.SchemaVersion, SchemaVersion)
	}
	if c.CloneProtocol != CloneProtocolHTTPS && c.CloneProtocol != CloneProtocolSSH {
		return fmt.Errorf("clone_protocol must be %q or %q", CloneProtocolHTTPS, CloneProtocolSSH)
	}
	return nil
}
