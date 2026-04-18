package annotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"grpcmon/internal/annotate"
)

func TestSave_And_LoadInto_RoundTrip(t *testing.T) {
	src := annotate.New()
	_ = src.Set("x", "note x")
	_ = src.Set("y", "note y")

	tmp := filepath.Join(t.TempDir(), "annotations.json")
	if err := annotate.Save(src, tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	dst := annotate.New()
	if err := annotate.LoadInto(dst, tmp); err != nil {
		t.Fatalf("LoadInto: %v", err)
	}

	for _, id := range []string{"x", "y"} {
		note, err := dst.Get(id)
		if err != nil {
			t.Errorf("Get(%q): %v", id, err)
		}
		orig, _ := src.Get(id)
		if note != orig {
			t.Errorf("id %q: got %q want %q", id, note, orig)
		}
	}
}

func TestLoadInto_MissingFile(t *testing.T) {
	s := annotate.New()
	err := annotate.LoadInto(s, "/nonexistent/path/annotations.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	s := annotate.New()
	_ = s.Set("z", "note z")
	tmp := filepath.Join(t.TempDir(), "out.json")
	if err := annotate.Save(s, tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(tmp); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
