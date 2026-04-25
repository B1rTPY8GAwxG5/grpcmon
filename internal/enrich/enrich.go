// Package enrich provides functionality for attaching derived metadata
// to captured gRPC entries, such as resolved hostnames, environment labels,
// and request size estimates.
package enrich

import (
	"encoding/json"
	"strings"

	"github.com/user/grpcmon/internal/capture"
)

// Enricher applies a set of enrichment functions to a capture entry,
// returning a new entry with additional metadata populated.
type Enricher struct {
	steps []StepFunc
}

// StepFunc is a function that accepts an entry and returns a (possibly
// modified) copy of it. Steps must not mutate the original entry.
type StepFunc func(e capture.Entry) capture.Entry

// New creates an Enricher with the provided steps applied in order.
func New(steps ...StepFunc) *Enricher {
	return &Enricher{steps: steps}
}

// Apply runs all enrichment steps over e and returns the result.
func (en *Enricher) Apply(e capture.Entry) capture.Entry {
	for _, step := range en.steps {
		e = step(e)
	}
	return e
}

// ApplyAll enriches every entry in the slice and returns the results.
func (en *Enricher) ApplyAll(entries []capture.Entry) []capture.Entry {
	out := make([]capture.Entry, len(entries))
	for i, e := range entries {
		out[i] = en.Apply(e)
	}
	return out
}

// WithEnvLabel returns a StepFunc that sets the "env" key in the entry's
// metadata to the supplied label (e.g. "production", "staging").
func WithEnvLabel(label string) StepFunc {
	return func(e capture.Entry) capture.Entry {
		e.Metadata = cloneMetadata(e.Metadata)
		e.Metadata["env"] = label
		return e
	}
}

// WithServiceName extracts the service portion of a fully-qualified gRPC
// method (/package.Service/Method) and stores it under the "service" key.
func WithServiceName() StepFunc {
	return func(e capture.Entry) capture.Entry {
		service := parseService(e.Method)
		if service == "" {
			return e
		}
		e.Metadata = cloneMetadata(e.Metadata)
		e.Metadata["service"] = service
		return e
	}
}

// WithRequestSize estimates the serialised size of the request payload in
// bytes (via JSON marshalling) and stores it as "request_size_bytes".
func WithRequestSize() StepFunc {
	return func(e capture.Entry) capture.Entry {
		if e.Request == nil {
			return e
		}
		b, err := json.Marshal(e.Request)
		if err != nil {
			return e
		}
		e.Metadata = cloneMetadata(e.Metadata)
		e.Metadata["request_size_bytes"] = len(b)
		return e
	}
}

// parseService returns the service segment from a gRPC method string of the
// form /package.Service/Method, or an empty string when the format is unknown.
func parseService(method string) string {
	// Expected format: /package.Service/Method
	trimmed := strings.TrimPrefix(method, "/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) != 2 || parts[0] == "" {
		return ""
	}
	// parts[0] is "package.Service" — return only the Service segment.
	pkg := parts[0]
	dot := strings.LastIndex(pkg, ".")
	if dot < 0 {
		return pkg
	}
	return pkg[dot+1:]
}

// cloneMetadata returns a shallow copy of m, creating a new map when m is nil.
func cloneMetadata(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}
