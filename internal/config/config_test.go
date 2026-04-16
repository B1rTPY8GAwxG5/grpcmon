package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/grpcmon/grpcmon/internal/config"
)

func TestDefaults_Values(t *testing.T) {
	cfg := config.Defaults()
	if cfg.Capture.MaxEntries != 1000 {
		t.Errorf("expected 1000, got %d", cfg.Capture.MaxEntries)
	}
	if cfg.Replay.Timeout != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.Replay.Timeout)
	}
	if cfg.Export.Format != "json" {
		t.Errorf("expected json, got %s", cfg.Export.Format)
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := config.Validate(config.Defaults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidMaxEntries(t *testing.T) {
	cfg := config.Defaults()
	cfg.Capture.MaxEntries = 0
	if err := config.Validate(cfg); err == nil {
		t.Fatal("expected error for zero max_entries")
	}
}

func TestValidate_InvalidTimeout(t *testing.T) {
	cfg := config.Defaults()
	cfg.Replay.Timeout = 0
	if err := config.Validate(cfg); err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	cfg := config.Defaults()
	cfg.Export.Format = "csv"
	if err := config.Validate(cfg); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/grpcmon.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	content := []byte("capture:\n  max_entries: 500\nreplay:\n  timeout: 5s\nexport:\n  format: json\n")
	f, err := os.CreateTemp(t.TempDir(), "grpcmon-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(content)
	f.Close()

	cfg, err := config.Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Capture.MaxEntries != 500 {
		t.Errorf("expected 500, got %d", cfg.Capture.MaxEntries)
	}
	if cfg.Replay.Timeout != 5*time.Second {
		t.Errorf("expected 5s, got %v", cfg.Replay.Timeout)
	}
}
