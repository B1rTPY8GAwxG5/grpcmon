package truncate_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/truncate"
)

func makeStore(n int) *capture.Store {
	s := capture.NewStore(1000)
	for i := 0; i < n; i++ {
		s.Add(capture.Entry{
			ID:        capture.NewID(),
			Method:    "/svc/Method",
			Timestamp: time.Now(),
		})
	}
	return s
}

func TestNew_DefaultsMaxSizeToOne(t *testing.T) {
	tr := truncate.New(makeStore(0), 0)
	if tr.MaxSize() != 1 {
		t.Fatalf("expected MaxSize 1, got %d", tr.MaxSize())
	}
}

func TestTrim_NothingToRemove(t *testing.T) {
	s := makeStore(3)
	tr := truncate.New(s, 10)
	removed := tr.Trim()
	if removed != 0 {
		t.Fatalf("expected 0 removed, got %d", removed)
	}
	if len(s.List()) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(s.List()))
	}
}

func TestTrim_RemovesExcess(t *testing.T) {
	s := makeStore(10)
	tr := truncate.New(s, 4)
	removed := tr.Trim()
	if removed != 6 {
		t.Fatalf("expected 6 removed, got %d", removed)
	}
	if len(s.List()) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(s.List()))
	}
}

func TestTrim_KeepsNewest(t *testing.T) {
	s := capture.NewStore(1000)
	var ids []string
	for i := 0; i < 5; i++ {
		id := capture.NewID()
		ids = append(ids, id)
		s.Add(capture.Entry{ID: id, Method: "/svc/M", Timestamp: time.Now()})
	}
	tr := truncate.New(s, 3)
	tr.Trim()

	remaining := s.List()
	if len(remaining) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(remaining))
	}
	// Newest 3 ids should be kept.
	for i, e := range remaining {
		if e.ID != ids[2+i] {
			t.Errorf("entry %d: expected id %s, got %s", i, ids[2+i], e.ID)
		}
	}
}
