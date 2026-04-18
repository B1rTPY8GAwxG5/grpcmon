package tag_test

import (
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/tag"
)

func TestAdd_And_Lookup(t *testing.T) {
	s := tag.New()
	s.Add("id1", "slow", "error")
	s.Add("id2", "slow")

	ids := s.Lookup("slow")
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}
	ids = s.Lookup("error")
	if len(ids) != 1 || ids[0] != "id1" {
		t.Fatalf("expected [id1], got %v", ids)
	}
}

func TestLookup_UnknownTag(t *testing.T) {
	s := tag.New()
	ids := s.Lookup("missing")
	if len(ids) != 0 {
		t.Fatalf("expected empty, got %v", ids)
	}
}

func TestRemove_DisassociatesID(t *testing.T) {
	s := tag.New()
	s.Add("id1", "slow")
	s.Remove("id1", "slow")
	ids := s.Lookup("slow")
	if len(ids) != 0 {
		t.Fatalf("expected empty after remove, got %v", ids)
	}
}

func TestFilter_ReturnsMatchingEntries(t *testing.T) {
	s := tag.New()
	s.Add("abc", "important")

	entries := []capture.Entry{
		{ID: "abc", Method: "/svc/A"},
		{ID: "xyz", Method: "/svc/B"},
	}
	out := s.Filter("important", entries)
	if len(out) != 1 || out[0].ID != "abc" {
		t.Fatalf("unexpected filter result: %v", out)
	}
}

func TestFilter_NoMatch(t *testing.T) {
	s := tag.New()
	entries := []capture.Entry{{ID: "abc"}}
	out := s.Filter("nope", entries)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %v", out)
	}
}

func TestTags_ReturnsAllTags(t *testing.T) {
	s := tag.New()
	s.Add("id1", "a", "b")
	s.Add("id2", "c")
	tags := s.Tags()
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(tags), tags)
	}
}
