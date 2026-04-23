package transform

import (
	"strings"

	"github.com/example/grpcmon/internal/capture"
)

// RedactMetadataKey returns a Func that replaces the value of the given
// metadata key (case-insensitive) with "REDACTED".
func RedactMetadataKey(key string) Func {
	lower := strings.ToLower(key)
	return func(e capture.Entry) (capture.Entry, bool) {
		if e.Metadata == nil {
			return e, true
		}
		copy := make(map[string]string, len(e.Metadata))
		for k, v := range e.Metadata {
			if strings.ToLower(k) == lower {
				copy[k] = "REDACTED"
			} else {
				copy[k] = v
			}
		}
		e.Metadata = copy
		return e, true
	}
}

// KeepMethods returns a Func that drops entries whose Method is not in the
// provided allow-list.
func KeepMethods(methods ...string) Func {
	allowed := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		allowed[m] = struct{}{}
	}
	return func(e capture.Entry) (capture.Entry, bool) {
		_, ok := allowed[e.Method]
		if !ok {
			return capture.Entry{}, false
		}
		return e, true
	}
}

// NormaliseMethod returns a Func that trims leading slashes and lower-cases
// the method name for consistent comparison.
func NormaliseMethod() Func {
	return func(e capture.Entry) (capture.Entry, bool) {
		e.Method = strings.ToLower(strings.TrimLeft(e.Method, "/"))
		return e, true
	}
}
