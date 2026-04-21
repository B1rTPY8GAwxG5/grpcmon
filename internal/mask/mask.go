// Package mask provides utilities for redacting sensitive fields from
// captured gRPC entries before display, export, or replay.
package mask

import (
	"strings"

	"github.com/user/grpcmon/internal/capture"
)

const redacted = "[REDACTED]"

// Masker redacts values whose metadata keys match a configured set of
// case-insensitive field names.
type Masker struct {
	fields map[string]struct{}
}

// New returns a Masker that will redact the supplied field names.
// Field matching is case-insensitive.
func New(fields ...string) *Masker {
	m := &Masker{fields: make(map[string]struct{}, len(fields))}
	for _, f := range fields {
		m.fields[strings.ToLower(f)] = struct{}{}
	}
	return m
}

// Apply returns a copy of e with matching metadata values redacted.
// The original entry is never modified.
func (m *Masker) Apply(e capture.Entry) capture.Entry {
	if len(e.Metadata) == 0 {
		return e
	}
	copied := make(map[string]string, len(e.Metadata))
	for k, v := range e.Metadata {
		if _, ok := m.fields[strings.ToLower(k)]; ok {
			copied[k] = redacted
		} else {
			copied[k] = v
		}
	}
	e.Metadata = copied
	return e
}

// ApplyAll returns a new slice with Apply called on every entry.
func (m *Masker) ApplyAll(entries []capture.Entry) []capture.Entry {
	out := make([]capture.Entry, len(entries))
	for i, e := range entries {
		out[i] = m.Apply(e)
	}
	return out
}
