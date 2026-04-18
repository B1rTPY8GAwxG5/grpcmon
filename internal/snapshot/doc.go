// Package snapshot provides save, load, and list operations for named
// snapshots of captured gRPC traffic.
//
// Snapshots are persisted as JSON files in a configurable directory,
// allowing developers to capture a set of requests, save them, and
// reload them later for replay or diffing.
//
// Basic usage:
//
//	// Save current capture store to a named snapshot
//	meta, err := snapshot.Save("/tmp/snaps", "baseline", store)
//
//	// Load entries back from a snapshot
//	entries, err := snapshot.Load("/tmp/snaps", "baseline")
//
//	// List all available snapshots
//	names, err := snapshot.List("/tmp/snaps")
package snapshot
