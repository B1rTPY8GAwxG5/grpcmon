package format_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/format"
)

func TestWriteTable_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := format.WriteTable(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"STATUS", "METHOD", "LATENCY", "TIMESTAMP", "ID"} {
		if !strings.Contains(out, col) {
			t.Errorf("header missing column %q", col)
		}
	}
}

func TestWriteTable_ContainsSeparator(t *testing.T) {
	var buf bytes.Buffer
	_ = format.WriteTable(&buf, nil)
	if !strings.Contains(buf.String(), "---") {
		t.Error("expected separator line in table output")
	}
}

func TestWriteTable_RowPerEntry(t *testing.T) {
	entries := []capture.Entry{
		{
			ID:        "id1",
			Method:    "/svc.A/MethodOne",
			Duration:  5 * time.Millisecond,
			Timestamp: time.Now(),
		},
		{
			ID:        "id2",
			Method:    "/svc.B/MethodTwo",
			Duration:  12 * time.Millisecond,
			Timestamp: time.Now(),
		},
	}

	var buf bytes.Buffer
	if err := format.WriteTable(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "id1") {
		t.Error("expected id1 in output")
	}
	if !strings.Contains(out, "id2") {
		t.Error("expected id2 in output")
	}
}

func TestWriteTable_TruncatesLongMethod(t *testing.T) {
	long := "/" + strings.Repeat("x", 60) + "/Method"
	entries := []capture.Entry{
		{ID: "id1", Method: long, Timestamp: time.Now()},
	}

	var buf bytes.Buffer
	_ = format.WriteTable(&buf, entries)

	// The full method should not appear verbatim; it must be truncated.
	if strings.Contains(buf.String(), long) {
		t.Error("expected long method to be truncated")
	}
}
