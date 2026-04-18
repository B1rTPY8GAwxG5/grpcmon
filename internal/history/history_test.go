package history

import (
	"errors"
	"testing"

	"grpcmon/internal/capture"
)

func entry(id, method string) capture.Entry {
	return capture.Entry{ID: id, Method: method}
}

func TestAdd_RecordsSuccess(t *testing.T) {
	h := New(10)
	h.Add(entry("1", "/svc/Method"), nil)
	records := h.List()
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if !records[0].Success {
		t.Error("expected success=true")
	}
	if records[0].Method != "/svc/Method" {
		t.Errorf("unexpected method %q", records[0].Method)
	}
}

func TestAdd_RecordsError(t *testing.T) {
	h := New(10)
	h.Add(entry("2", "/svc/Fail"), errors.New("deadline exceeded"))
	records := h.List()
	if records[0].Success {
		t.Error("expected success=false")
	}
	if records[0].Error != "deadline exceeded" {
		t.Errorf("unexpected error string %q", records[0].Error)
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	h := New(3)
	for i := 0; i < 4; i++ {
		h.Add(entry(string(rune('a'+i)), "/m"), nil)
	}
	records := h.List()
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
	if records[0].EntryID != "b" {
		t.Errorf("expected oldest evicted, got %q", records[0].EntryID)
	}
}

func TestClear_RemovesAll(t *testing.T) {
	h := New(10)
	h.Add(entry("1", "/m"), nil)
	h.Clear()
	if len(h.List()) != 0 {
		t.Error("expected empty history after Clear")
	}
}

func TestNew_ZeroMaxDefaultsTo100(t *testing.T) {
	h := New(0)
	if h.maxSize != 100 {
		t.Errorf("expected maxSize=100, got %d", h.maxSize)
	}
}
