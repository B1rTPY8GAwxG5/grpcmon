package format_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/format"
)

func sampleEntry() capture.Entry {
	return capture.Entry{
		ID:         "abc123",
		Method:     "/hello.Greeter/SayHello",
		StatusCode: 0,
		Duration:   42 * time.Millisecond,
		Timestamp:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Request:    []byte(`{"name":"world"}`),
		Response:   []byte(`{"message":"hello world"}`),
	}
}

func TestDefaultFormatter_Values(t *testing.T) {
	f := format.DefaultFormatter()
	if f.Verbose {
		t.Error("expected Verbose to be false by default")
	}
	if f.TimeFormat == "" {
		t.Error("expected TimeFormat to be non-empty")
	}
}

func TestFprint_ContainsMethod(t *testing.T) {
	var buf bytes.Buffer
	f := format.DefaultFormatter()
	e := sampleEntry()

	if err := f.Fprint(&buf, e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), e.Method) {
		t.Errorf("output missing method, got: %s", buf.String())
	}
}

func TestFprint_OKStatus(t *testing.T) {
	var buf bytes.Buffer
	f := format.DefaultFormatter()
	e := sampleEntry()
	e.StatusCode = 0

	_ = f.Fprint(&buf, e)

	if !strings.Contains(buf.String(), "OK") {
		t.Errorf("expected OK status in output, got: %s", buf.String())
	}
}

func TestFprint_ErrorStatus(t *testing.T) {
	var buf bytes.Buffer
	f := format.DefaultFormatter()
	e := sampleEntry()
	e.StatusCode = 14

	_ = f.Fprint(&buf, e)

	if !strings.Contains(buf.String(), "ERR(14)") {
		t.Errorf("expected ERR(14) in output, got: %s", buf.String())
	}
}

func TestFprint_VerboseIncludesBodies(t *testing.T) {
	var buf bytes.Buffer
	f := format.DefaultFormatter()
	f.Verbose = true
	e := sampleEntry()

	_ = f.Fprint(&buf, e)
	out := buf.String()

	if !strings.Contains(out, "world") {
		t.Errorf("expected request body in verbose output, got: %s", out)
	}
	if !strings.Contains(out, "hello world") {
		t.Errorf("expected response body in verbose output, got: %s", out)
	}
}

func TestFprintAll_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	f := format.DefaultFormatter()
	entries := []capture.Entry{sampleEntry(), sampleEntry()}
	entries[1].Method = "/other.Service/Call"

	if err := f.FprintAll(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
