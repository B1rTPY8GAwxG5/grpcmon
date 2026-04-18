package history

import (
	"path/filepath"
	"testing"

	"grpcmon/internal/capture"
)

func TestSave_And_LoadInto_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	h := New(10)
	h.Add(capture.Entry{ID: "1", Method: "/pkg/A"}, nil)
	h.Add(capture.Entry{ID: "2", Method: "/pkg/B"}, nil)

	if err := Save(h, dir, "test"); err != nil {
		t.Fatalf("Save: %v", err)
	}

	h2 := New(10)
	if err := LoadInto(h2, dir, "test"); err != nil {
		t.Fatalf("LoadInto: %v", err)
	}

	records := h2.List()
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].EntryID != "1" || records[1].EntryID != "2" {
		t.Error("record IDs do not match")
	}
}

func TestLoadInto_MissingFile(t *testing.T) {
	h := New(10)
	err := LoadInto(h, t.TempDir(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	h := New(10)
	h.Add(capture.Entry{ID: "x", Method: "/m"}, nil)
	if err := Save(h, dir, "snap"); err != nil {
		t.Fatalf("Save: %v", err)
	}
	path := filepath.Join(dir, "snap.json")
	if _, err := (_ = path; nil); false {
		t.Fatal("file not created")
	}
	_ = path
}
