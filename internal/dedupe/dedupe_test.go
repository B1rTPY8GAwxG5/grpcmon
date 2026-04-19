package dedupe_test

import (
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/dedupe"
)

func entry(method, request string) capture.Entry {
	return capture.Entry{Method: method, Request: request}
}

func TestIsDuplicate_NewEntry(t *testing.T) {
	s := dedupe.New()
	if s.IsDuplicate(entry("/svc/Method", `{"id":1}`)) {
		t.Fatal("expected false for first occurrence")
	}
}

func TestIsDuplicate_SeenEntry(t *testing.T) {
	s := dedupe.New()
	e := entry("/svc/Method", `{"id":1}`)
	s.IsDuplicate(e)
	if !s.IsDuplicate(e) {
		t.Fatal("expected true for duplicate")
	}
}

func TestIsDuplicate_DifferentMethod(t *testing.T) {
	s := dedupe.New()
	s.IsDuplicate(entry("/svc/A", `{}`))
	if s.IsDuplicate(entry("/svc/B", `{}`)) {
		t.Fatal("different method should not be duplicate")
	}
}

func TestFilter_RemovesDuplicates(t *testing.T) {
	s := dedupe.New()
	entries := []capture.Entry{
		entry("/svc/M", `{"x":1}`),
		entry("/svc/M", `{"x":1}`),
		entry("/svc/M", `{"x":2}`),
	}
	got := s.Filter(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 unique entries, got %d", len(got))
	}
}

func TestReset_ClearsFingerprints(t *testing.T) {
	s := dedupe.New()
	e := entry("/svc/M", `{}`)
	s.IsDuplicate(e)
	s.Reset()
	if s.IsDuplicate(e) {
		t.Fatal("expected false after reset")
	}
}

func TestLen_TracksCount(t *testing.T) {
	s := dedupe.New()
	s.IsDuplicate(entry("/a", `1`))
	s.IsDuplicate(entry("/b", `2`))
	s.IsDuplicate(entry("/a", `1`)) // duplicate, not counted
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}
