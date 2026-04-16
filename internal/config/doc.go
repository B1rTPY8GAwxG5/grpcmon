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
package config
