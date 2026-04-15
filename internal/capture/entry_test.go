package capture_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
)

func TestEntry_ZeroValue(t *testing.T) {
	var e capture.Entry
	if e.ID != "" {
		t.Errorf("expected empty ID, got %q", e.ID)
	}
	if e.Duration != 0 {
		t.Errorf("expected zero duration, got %v", e.Duration)
	}
	if !e.Timestamp.IsZero() {
		t.Error("expected zero timestamp")
	}
}

func TestEntry_FieldAssignment(t *testing.T) {
	now := time.Now()
	e := capture.Entry{
		ID:           "abc-123",
		Method:       "/pkg.Svc/Do",
		StatusCode:   "OK",
		Duration:     42 * time.Millisecond,
		Timestamp:    now,
		RequestBody:  []byte(`{"name":"world"}`),
		ResponseBody: []byte(`{"message":"hello world"}`),
		Metadata:     map[string][]string{"x-request-id": {"req-1"}},
	}

	if e.ID != "abc-123" {
		t.Errorf("unexpected ID: %s", e.ID)
	}
	if e.Method != "/pkg.Svc/Do" {
		t.Errorf("unexpected Method: %s", e.Method)
	}
	if e.StatusCode != "OK" {
		t.Errorf("unexpected StatusCode: %s", e.StatusCode)
	}
	if e.Duration != 42*time.Millisecond {
		t.Errorf("unexpected Duration: %v", e.Duration)
	}
	if !e.Timestamp.Equal(now) {
		t.Errorf("unexpected Timestamp: %v", e.Timestamp)
	}
	if string(e.RequestBody) != `{"name":"world"}` {
		t.Errorf("unexpected RequestBody: %s", e.RequestBody)
	}
	if vals, ok := e.Metadata["x-request-id"]; !ok || vals[0] != "req-1" {
		t.Errorf("unexpected Metadata: %v", e.Metadata)
	}
}
