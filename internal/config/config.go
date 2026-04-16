// Package config provides configuration loading and validation for grpcmon.
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level grpcmon configuration.
type Config struct {
	Capture CaptureConfig `yaml:"capture"`
	Replay  ReplayConfig  `yaml:"replay"`
	Export  ExportConfig  `yaml:"export"`
}

// CaptureConfig controls traffic capture behaviour.
type CaptureConfig struct {
	MaxEntries int `yaml:"max_entries"`
}

// ReplayConfig controls replay behaviour.
type ReplayConfig struct {
	Target  string        `yaml:"target"`
	Timeout time.Duration `yaml:"timeout"`
}

// ExportConfig controls export behaviour.
type ExportConfig struct {
	Format string `yaml:"format"`
	Path   string `yaml:"path"`
}

// Defaults returns a Config populated with sensible defaults.
func Defaults() Config {
	return Config{
		Capture: CaptureConfig{MaxEntries: 1000},
		Replay:  ReplayConfig{Timeout: 10 * time.Second},
		Export:  ExportConfig{Format: "json"},
	}
}

// Load reads a YAML config file from path, merging over defaults.
func Load(path string) (Config, error) {
	cfg := Defaults()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, Validate(cfg)
}

// Validate returns an error if cfg contains invalid values.
func Validate(cfg Config) error {
	if cfg.Capture.MaxEntries <= 0 {
		return errors.New("config: capture.max_entries must be greater than zero")
	}
	if cfg.Replay.Timeout <= 0 {
		return errors.New("config: replay.timeout must be greater than zero")
	}
	if cfg.Export.Format != "" && cfg.Export.Format != "json" {
		return errors.New("config: export.format must be \"json\" or empty")
	}
	return nil
}
