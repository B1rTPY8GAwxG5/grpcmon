package bookmark_test

import (
	"testing"

	"github.com/grpcmon/internal/bookmark"
	"github.com/grpcmon/internal/capture"
)

func entry(id string) capture.Entry {
	return capture.Entry{ID: id, Method: "/svc/Method"}
}

func TestAdd_And_Get(t *testing.T) {
	s := bookmark.New()
	if err := s.Add("first", entry("abc")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id, err := s.Get("first")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if id != "abc" {
		t.Errorf("got id %q, want %q", id, "abc")
	}
}

func TestAdd_Duplicate(t *testing.T) {
	s := bookmark.New()
	_ = s.Add("dup", entry("1"))
	if err := s.Add("dup", entry("2")); err != bookmark.ErrDuplicate {
		t.Errorf("expected ErrDuplicate, got %v", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := bookmark.New()
	_, err := s.Get("missing")
	if err != bookmark.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRemove_Disassociates(t *testing.T) {
	s := bookmark.New()
	_ = s.Add("r", entry("xyz"))
	if err := s.Remove("r"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if _, err := s.Get("r"); err != bookmark.ErrNotFound {
		t.Errorf("expected ErrNotFound after remove, got %v", err)
	}
}

func TestRemove_NotFound(t *testing.T) {
	s := bookmark.New()
	if err := s.Remove("ghost"); err != bookmark.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList_ReturnsAllNames(t *testing.T) {
	s := bookmark.New()
	_ = s.Add("a", entry("1"))
	_ = s.Add("b", entry("2"))
	_ = s.Add("c", entry("3"))
	names := s.List()
	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}
}
