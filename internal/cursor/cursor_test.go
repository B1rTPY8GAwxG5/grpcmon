package cursor_test

import (
	"testing"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/cursor"
)

func makeEntries(methods ...string) []capture.Entry {
	out := make([]capture.Entry, len(methods))
	for i, m := range methods {
		out[i] = capture.Entry{Method: m}
	}
	return out
}

func TestCurrent_ReturnsFirstEntry(t *testing.T) {
	c := cursor.New(makeEntries("/svc/A", "/svc/B"))
	e, err := c.Current()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Method != "/svc/A" {
		t.Errorf("expected /svc/A, got %s", e.Method)
	}
}

func TestNext_AdvancesCursor(t *testing.T) {
	c := cursor.New(makeEntries("/svc/A", "/svc/B", "/svc/C"))
	e, err := c.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Method != "/svc/B" {
		t.Errorf("expected /svc/B, got %s", e.Method)
	}
	if c.Pos() != 1 {
		t.Errorf("expected pos 1, got %d", c.Pos())
	}
}

func TestNext_DoesNotExceedBounds(t *testing.T) {
	c := cursor.New(makeEntries("/svc/A"))
	c.Next() // already at last
	e, err := c.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Method != "/svc/A" {
		t.Errorf("expected /svc/A, got %s", e.Method)
	}
}

func TestPrev_MovesCursorBack(t *testing.T) {
	c := cursor.New(makeEntries("/svc/A", "/svc/B"))
	c.Next()
	e, err := c.Prev()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Method != "/svc/A" {
		t.Errorf("expected /svc/A, got %s", e.Method)
	}
}

func TestPrev_DoesNotGoBelowZero(t *testing.T) {
	c := cursor.New(makeEntries("/svc/A", "/svc/B"))
	c.Prev()
	if c.Pos() != 0 {
		t.Errorf("expected pos 0, got %d", c.Pos())
	}
}

func TestCurrent_EmptyEntries_ReturnsError(t *testing.T) {
	c := cursor.New(nil)
	_, err := c.Current()
	if err != cursor.ErrEmpty {
		t.Errorf("expected ErrEmpty, got %v", err)
	}
}

func TestReset_ReturnsCursorToStart(t *testing.T) {
	c := cursor.New(makeEntries("/svc/A", "/svc/B", "/svc/C"))
	c.Next()
	c.Next()
	c.Reset()
	if c.Pos() != 0 {
		t.Errorf("expected pos 0 after reset, got %d", c.Pos())
	}
}

func TestLen_ReturnsCorrectCount(t *testing.T) {
	c := cursor.New(makeEntries("/a", "/b", "/c"))
	if c.Len() != 3 {
		t.Errorf("expected 3, got %d", c.Len())
	}
}
