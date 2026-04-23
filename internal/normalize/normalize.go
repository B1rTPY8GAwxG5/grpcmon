// Package normalize provides utilities for normalising captured gRPC entries
// before storage, comparison, or export. It strips volatile fields such as
// per-request timestamps and trace IDs so that structurally identical calls
// can be compared reliably.
package normalize

import (
	"strings"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Option configures the behaviour of Normalize.
type Option func(*options)

type options struct {
	clearTimestamp bool
	lowerMethod    bool
	stripMetaKeys  []string
}

// ClearTimestamp replaces the entry timestamp with the zero value.
func ClearTimestamp() Option {
	return func(o *options) { o.clearTimestamp = true }
}

// LowerMethod lowercases the method name.
func LowerMethod() Option {
	return func(o *options) { o.lowerMethod = true }
}

// StripMetadataKeys removes the specified keys (case-insensitive) from
// the entry metadata map.
func StripMetadataKeys(keys ...string) Option {
	return func(o *options) { o.stripMetaKeys = append(o.stripMetaKeys, keys...) }
}

// Apply returns a copy of e with the requested normalisations applied.
// The original entry is never mutated.
func Apply(e capture.Entry, opts ...Option) capture.Entry {
	cfg := &options{}
	for _, o := range opts {
		o(cfg)
	}

	out := e

	if cfg.clearTimestamp {
		out.Timestamp = time.Time{}
	}

	if cfg.lowerMethod {
		out.Method = strings.ToLower(strings.TrimSpace(out.Method))
	}

	if len(cfg.stripMetaKeys) > 0 && len(out.Metadata) > 0 {
		newMeta := make(map[string]string, len(out.Metadata))
		for k, v := range out.Metadata {
			newMeta[k] = v
		}
		for _, key := range cfg.stripMetaKeys {
			delete(newMeta, strings.ToLower(key))
			delete(newMeta, key)
		}
		out.Metadata = newMeta
	}

	return out
}

// ApplyAll applies the given options to every entry in the slice and returns
// a new slice of normalised copies.
func ApplyAll(entries []capture.Entry, opts ...Option) []capture.Entry {
	out := make([]capture.Entry, len(entries))
	for i, e := range entries {
		out[i] = Apply(e, opts...)
	}
	return out
}
