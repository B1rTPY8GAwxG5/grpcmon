package config

import (
	"os"
	"strconv"
	"time"
)

// FromEnv overlays environment variables onto cfg and returns the result.
// Recognised variables:
//
//	GRPCMON_MAX_ENTRIES   – capture.max_entries (integer)
//	GRPCMON_TARGET        – replay.target (string)
//	GRPCMON_TIMEOUT       – replay.timeout (duration, e.g. "5s")
//	GRPCMON_EXPORT_FORMAT – export.format (string)
//	GRPCMON_EXPORT_PATH   – export.path (string)
func FromEnv(cfg Config) Config {
	if v := os.Getenv("GRPCMON_MAX_ENTRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Capture.MaxEntries = n
		}
	}
	if v := os.Getenv("GRPCMON_TARGET"); v != "" {
		cfg.Replay.Target = v
	}
	if v := os.Getenv("GRPCMON_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Replay.Timeout = d
		}
	}
	if v := os.Getenv("GRPCMON_EXPORT_FORMAT"); v != "" {
		cfg.Export.Format = v
	}
	if v := os.Getenv("GRPCMON_EXPORT_PATH"); v != "" {
		cfg.Export.Path = v
	}
	return cfg
}
