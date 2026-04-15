package diff_test

import (
	"strings"
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/diff"
)

func baseEntry() capture.Entry {
	return capture.Entry{
		Method:     "/pkg.Service/Method",
		StatusCode: 0,
		Response:   `{"ok":true}`,
		Error:      "",
	}
}

func TestCompare_IdenticalEntries(t *testing.T) {
	a := baseEntry()
	b := baseEntry()

	result := diff.Compare(a, b)

	if !result.Match {
		t.Errorf("expected match, got differences: %v", result.Differences)
	}
}

func TestCompare_DifferentMethod(t *testing.T) {
	a := baseEntry()
	b := baseEntry()
	b.Method = "/pkg.Service/Other"

	result := diff.Compare(a, b)

	if result.Match {
		t.Fatal("expected mismatch for different methods")
	}
	if len(result.Differences) != 1 || !strings.Contains(result.Differences[0], "method") {
		t.Errorf("unexpected differences: %v", result.Differences)
	}
}

func TestCompare_DifferentStatusAndResponse(t *testing.T) {
	a := baseEntry()
	b := baseEntry()
	b.StatusCode = 2
	b.Response = `{"ok":false}`

	result := diff.Compare(a, b)

	if result.Match {
		t.Fatal("expected mismatch")
	}
	if len(result.Differences) != 2 {
		t.Errorf("expected 2 differences, got %d: %v", len(result.Differences), result.Differences)
	}
}

func TestResult_String_Match(t *testing.T) {
	r := diff.Result{Match: true}
	if r.String() != "entries match" {
		t.Errorf("unexpected string: %q", r.String())
	}
}

func TestResult_String_Mismatch(t *testing.T) {
	r := diff.Result{Match: false, Differences: []string{"method: \"a\" != \"b\""}}
	if !strings.Contains(r.String(), "differences found") {
		t.Errorf("unexpected string: %q", r.String())
	}
}

func TestCompareAll_LengthMismatch(t *testing.T) {
	as := []capture.Entry{baseEntry(), baseEntry()}
	bs := []capture.Entry{baseEntry()}

	results := diff.CompareAll(as, bs)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Match != true {
		t.Error("first pair should match")
	}
	if results[1].Match != false {
		t.Error("second pair should not match (length mismatch)")
	}
}
