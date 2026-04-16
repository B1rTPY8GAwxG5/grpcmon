package config_test

import (
	"testing"
	"time"

	"github.com/grpcmon/grpcmon/internal/config"
)

func TestFromEnv_MaxEntries(t *testing.T) {
	t.Setenv("GRPCMON_MAX_ENTRIES", "250")
	cfg := config.FromEnv(config.Defaults())
	if cfg.Capture.MaxEntries != 250 {
		t.Errorf("expected 250, got %d", cfg.Capture.MaxEntries)
	}
}

func TestFromEnv_Target(t *testing.T) {
	t.Setenv("GRPCMON_TARGET", "localhost:9090")
	cfg := config.FromEnv(config.Defaults())
	if cfg.Replay.Target != "localhost:9090" {
		t.Errorf("unexpected target: %s", cfg.Replay.Target)
	}
}

func TestFromEnv_Timeout(t *testing.T) {
	t.Setenv("GRPCMON_TIMEOUT", "30s")
	cfg := config.FromEnv(config.Defaults())
	if cfg.Replay.Timeout != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.Replay.Timeout)
	}
}

func TestFromEnv_ExportFormat(t *testing.T) {
	t.Setenv("GRPCMON_EXPORT_FORMAT", "json")
	cfg := config.FromEnv(config.Defaults())
	if cfg.Export.Format != "json" {
		t.Errorf("unexpected format: %s", cfg.Export.Format)
	}
}

func TestFromEnv_InvalidMaxEntries_Ignored(t *testing.T) {
	t.Setenv("GRPCMON_MAX_ENTRIES", "notanumber")
	cfg := config.FromEnv(config.Defaults())
	if cfg.Capture.MaxEntries != 1000 {
		t.Errorf("expected default 1000, got %d", cfg.Capture.MaxEntries)
	}
}

func TestFromEnv_NoEnv_PreservesDefaults(t *testing.T) {
	defaults := config.Defaults()
	cfg := config.FromEnv(defaults)
	if cfg.Capture.MaxEntries != defaults.Capture.MaxEntries {
		t.Error("defaults should be preserved when no env vars set")
	}
}
