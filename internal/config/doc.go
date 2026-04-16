// Package config handles loading, parsing, and validating grpcmon
// configuration from YAML files.
//
// A configuration file may look like:
//
//	capture:
//	  max_entries: 500
//	replay:
//	  target: localhost:50051
//	  timeout: 10s
//	export:
//	  format: json
//	  path: ./captures.json
//
// Call [Load] to read a file from disk, or [Defaults] to obtain a
// configuration with sensible default values.
//
// # Validation
//
// After loading, configuration values are validated automatically.
// Validation errors are returned as a [ValidationError] which lists
// all offending fields, allowing callers to surface actionable
// messages to users without requiring multiple load attempts.
package config
