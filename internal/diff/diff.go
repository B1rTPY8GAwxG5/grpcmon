// Package diff provides utilities for comparing gRPC capture entries,
// useful for spotting regressions between replayed and original responses.
package diff

import (
	"fmt"
	"strings"

	"github.com/grpcmon/internal/capture"
)

// Result holds the outcome of comparing two capture entries.
type Result struct {
	Match      bool
	Differences []string
}

// String returns a human-readable summary of the diff result.
func (r Result) String() string {
	if r.Match {
		return "entries match"
	}
	return fmt.Sprintf("differences found:\n  - %s", strings.Join(r.Differences, "\n  - "))
}

// Compare compares two capture entries and returns a Result describing
// any differences found in method, status code, and response payload.
func Compare(a, b capture.Entry) Result {
	var diffs []string

	if a.Method != b.Method {
		diffs = append(diffs, fmt.Sprintf("method: %q != %q", a.Method, b.Method))
	}

	if a.StatusCode != b.StatusCode {
		diffs = append(diffs, fmt.Sprintf("status_code: %v != %v", a.StatusCode, b.StatusCode))
	}

	if a.Response != b.Response {
		diffs = append(diffs, fmt.Sprintf("response: %q != %q", a.Response, b.Response))
	}

	if a.Error != b.Error {
		diffs = append(diffs, fmt.Sprintf("error: %q != %q", a.Error, b.Error))
	}

	return Result{
		Match:       len(diffs) == 0,
		Differences: diffs,
	}
}

// CompareAll compares slices of entries pairwise and returns one Result per pair.
// If the slices differ in length, a length mismatch difference is appended.
func CompareAll(as, bs []capture.Entry) []Result {
	results := make([]Result, 0, len(as))

	max := len(as)
	if len(bs) > max {
		max = len(bs)
	}

	for i := 0; i < max; i++ {
		if i >= len(as) || i >= len(bs) {
			results = append(results, Result{
				Match:       false,
				Differences: []string{fmt.Sprintf("entry %d: present in only one set", i)},
			})
			continue
		}
		results = append(results, Compare(as[i], bs[i]))
	}

	return results
}
