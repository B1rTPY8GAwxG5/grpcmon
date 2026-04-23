package window_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/window"
)

func makeStore(entries []capture.Entry) *capture.Store {
	s := capture.NewStore(100)
	for _, e := range entries {
		s.Add(e)
	}
	return s
}

func TestNew_DefaultsToOneMinute(t *testing.T) {
	s := capture.NewStore(10)
	w := window.New(s, 0)
	if w.Duration() != time.Minute {
		t.Fatalf("expected 1m, got %v", w.Duration())
	}
}

func TestEntries_ReturnsOnlyWithinWindow(t *testing.T) {
	now := time.Now()
	old := capture.Entry{ID: "old", Timestamp: now.Add(-5 * time.Minute)}
	recent := capture.Entry{ID: "recent", Timestamp: now.Add(-30 * time.Second)}

	s := makeStore([]capture.Entry{old, recent})
	w := window.New(s, time.Minute)

	got := w.Entries()
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].ID != "recent" {
		t.Errorf("expected 'recent', got %q", got[0].ID)
	}
}

func TestEntries_EmptyWhenAllOutsideWindow(t *testing.T) {
	now := time.Now()
	s := makeStore([]capture.Entry{
		{ID: "a", Timestamp: now.Add(-10 * time.Minute)},
		{ID: "b", Timestamp: now.Add(-2 * time.Minute)},
	})
	w := window.New(s, time.Minute)

	if got := w.Entries(); len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestEntries_AllWithinWindow(t *testing.T) {
	now := time.Now()
	s := makeStore([]capture.Entry{
		{ID: "x", Timestamp: now.Add(-10 * time.Second)},
		{ID: "y", Timestamp: now.Add(-20 * time.Second)},
	})
	w := window.New(s, time.Minute)

	if got := w.Entries(); len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestEntries_ExactlyAtBoundaryIncluded(t *testing.T) {
	now := time.Now()
	boundary := capture.Entry{ID: "boundary", Timestamp: now.Add(-time.Minute)}
	s := makeStore([]capture.Entry{boundary})
	w := window.New(s, time.Minute)

	// The cutoff is now-duration; an entry AT that time is not Before it.
	got := w.Entries()
	if len(got) != 1 {
		t.Fatalf("expected boundary entry to be included, got %d entries", len(got))
	}
}
