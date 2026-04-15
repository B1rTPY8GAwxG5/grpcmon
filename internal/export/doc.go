// Package export provides utilities for persisting and restoring captured
// gRPC traffic entries.
//
// Supported formats:
//
//	- JSON (FormatJSON): newline-terminated JSON array, human-readable and
//	  suitable for version-control or manual editing.
//
// Typical usage — save a capture session to disk:
//
//	f, _ := os.Create("capture.json")
//	defer f.Close()
//	export.Write(f, store.List(), export.FormatJSON)
//
// Load it back for replay:
//
//	f, _ := os.Open("capture.json")
//	defer f.Close()
//	entries, _ := export.Read(f, export.FormatJSON)
package export
