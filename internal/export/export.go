// Package export serialises captured gRPC entries to common interchange
// formats so that they can be stored, shared, or loaded back for replay.
package export

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/grpcmon/internal/capture"
)

// Format identifies the serialisation format used when writing entries.
type Format string

const (
	FormatJSON Format = "json"
)

// Write encodes entries into the chosen format and writes the result to w.
// Currently only FormatJSON is supported.
func Write(w io.Writer, entries []capture.Entry, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, entries)
	default:
		return fmt.Errorf("export: unsupported format %q", format)
	}
}

// Read decodes entries from r, inferring the format from the format argument.
func Read(r io.Reader, format Format) ([]capture.Entry, error) {
	switch format {
	case FormatJSON:
		return readJSON(r)
	default:
		return nil, fmt.Errorf("export: unsupported format %q", format)
	}
}

func writeJSON(w io.Writer, entries []capture.Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func readJSON(r io.Reader) ([]capture.Entry, error) {
	var entries []capture.Entry
	if err := json.NewDecoder(r).Decode(&entries); err != nil {
		return nil, fmt.Errorf("export: decode JSON: %w", err)
	}
	return entries, nil
}
