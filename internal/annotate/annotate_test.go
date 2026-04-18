package annotate_test

import (
	"testing"

	"grpcmon/internal/annotate"
)

func TestSet_And_Get(t *testing.T) {
	s := annotate.New()
	if err := s.Set("id1", "first note"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	note, err := s.Get("id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note != "first note" {
		t.Errorf("got %q, want %q", note, "first note")
	}
}

func TestGet_NotFound(t *testing.T) {
	s := annotate.New()
	_, err := s.Get("missing")
	if err != annotate.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSet_EmptyID(t *testing.T) {
	s := annotate.New()
	if err := s.Set("", "note"); err == nil {
		t.Error("expected error for empty id")
	}
}

func TestDelete_RemovesAnnotation(t *testing.T) {
	s := annotate.New()
	_ = s.Set("id2", "to delete")
	s.Delete("id2")
	_, err := s.Get("id2")
	if err != annotate.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestList_ReturnsCopy(t *testing.T) {
	s := annotate.New()
	_ = s.Set("a", "alpha")
	_ = s.Set("b", "beta")
	list := s.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
	// mutating the copy must not affect the store
	delete(list, "a")
	if _, err := s.Get("a"); err != nil {
		t.Error("store should still contain 'a' after mutating copy")
	}
}

func TestSet_OverwritesExisting(t *testing.T) {
	s := annotate.New()
	_ = s.Set("id3", "old")
	_ = s.Set("id3", "new")
	note, _ := s.Get("id3")
	if note != "new" {
		t.Errorf("expected 'new', got %q", note)
	}
}
