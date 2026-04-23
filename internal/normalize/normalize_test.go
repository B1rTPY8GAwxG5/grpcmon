package normalize_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/normalize"
)

func baseEntry() capture.Entry {
	return capture.Entry{
		ID:        "abc-123",
		Method:    "/pkg.Service/DoThing",
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"authorization": "Bearer token",
			"x-request-id":  "req-999",
			"content-type":  "application/grpc",
		},
	}
}

func TestApply_ClearTimestamp(t *testing.T) {
	e := baseEntry()
	out := normalize.Apply(e, normalize.ClearTimestamp())
	if !out.Timestamp.IsZero() {
		t.Errorf("expected zero timestamp, got %v", out.Timestamp)
	}
	if e.Timestamp.IsZero() {
		t.Error("original entry timestamp was mutated")
	}
}

func TestApply_LowerMethod(t *testing.T) {
	e := baseEntry()
	out := normalize.Apply(e, normalize.LowerMethod())
	want := "/pkg.service/dothing"
	if out.Method != want {
		t.Errorf("got method %q, want %q", out.Method, want)
	}
	if e.Method == want {
		t.Error("original entry method was mutated")
	}
}

func TestApply_StripMetadataKeys(t *testing.T) {
	e := baseEntry()
	out := normalize.Apply(e, normalize.StripMetadataKeys("authorization", "x-request-id"))

	if _, ok := out.Metadata["authorization"]; ok {
		t.Error("authorization key should have been stripped")
	}
	if _, ok := out.Metadata["x-request-id"]; ok {
		t.Error("x-request-id key should have been stripped")
	}
	if _, ok := out.Metadata["content-type"]; !ok {
		t.Error("content-type key should have been preserved")
	}
	// original must be untouched
	if _, ok := e.Metadata["authorization"]; !ok {
		t.Error("original metadata was mutated")
	}
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	e := baseEntry()
	out := normalize.Apply(e)
	if out.Method != e.Method {
		t.Errorf("method changed unexpectedly: %q", out.Method)
	}
	if out.Timestamp != e.Timestamp {
		t.Error("timestamp changed unexpectedly")
	}
}

func TestApplyAll_NormalisesEveryEntry(t *testing.T) {
	entries := []capture.Entry{baseEntry(), baseEntry(), baseEntry()}
	out := normalize.ApplyAll(entries, normalize.ClearTimestamp(), normalize.LowerMethod())

	if len(out) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(out))
	}
	for i, e := range out {
		if !e.Timestamp.IsZero() {
			t.Errorf("entry %d: timestamp not cleared", i)
		}
		if e.Method != "/pkg.service/dothing" {
			t.Errorf("entry %d: method not lowercased: %q", i, e.Method)
		}
	}
}
