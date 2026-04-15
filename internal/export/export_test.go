package export_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/grpcmon/internal/capture"
	"github.com/user/grpcmon/internal/export"
	"google.golang.org/grpc/codes"
)

func sampleEntries() []capture.Entry {
	return []capture.Entry{
		{
			ID:         "entry-001",
			Method:     "/helloworld.Greeter/SayHello",
			StatusCode: codes.OK,
			Duration:   42 * time.Millisecond,
		},
		{
			ID:         "entry-002",
			Method:     "/helloworld.Greeter/SayGoodbye",
			StatusCode: codes.NotFound,
			Duration:   7 * time.Millisecond,
		},
	}
}

func TestWrite_JSON_RoundTrip(t *testing.T) {
	entries := sampleEntries()

	var buf bytes.Buffer
	if err := export.Write(&buf, entries, export.FormatJSON); err != nil {
		t.Fatalf("Write: unexpected error: %v", err)
	}

	got, err := export.Read(&buf, export.FormatJSON)
	if err != nil {
		t.Fatalf("Read: unexpected error: %v", err)
	}

	if len(got) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(got))
	}
	for i, e := range entries {
		if got[i].ID != e.ID {
			t.Errorf("entry[%d].ID: want %q, got %q", i, e.ID, got[i].ID)
		}
		if got[i].Method != e.Method {
			t.Errorf("entry[%d].Method: want %q, got %q", i, e.Method, got[i].Method)
		}
		if got[i].StatusCode != e.StatusCode {
			t.Errorf("entry[%d].StatusCode: want %v, got %v", i, e.StatusCode, got[i].StatusCode)
		}
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, sampleEntries(), export.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRead_UnsupportedFormat(t *testing.T) {
	_, err := export.Read(strings.NewReader("{}"), export.Format("csv"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestRead_MalformedJSON(t *testing.T) {
	_, err := export.Read(strings.NewReader("not-json"), export.FormatJSON)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}
