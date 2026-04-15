package capture

import (
	"context"
	"testing"
	"time"
)

func TestStore_AddAndList(t *testing.T) {
	s := NewStore(10)
	if s.Len() != 0 {
		t.Fatalf("expected empty store, got %d", s.Len())
	}

	s.Add(&Entry{ID: "1", Method: "/pkg.Svc/Method"})
	s.Add(&Entry{ID: "2", Method: "/pkg.Svc/Other"})

	if s.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", s.Len())
	}

	list := s.List()
	if list[0].ID != "1" || list[1].ID != "2" {
		t.Errorf("unexpected entry order: %v", list)
	}
}

func TestStore_EvictsOldestWhenFull(t *testing.T) {
	s := NewStore(3)
	for i, id := range []string{"a", "b", "c", "d"} {
		s.Add(&Entry{ID: id, Method: "/svc/m", Timestamp: time.Now().Add(time.Duration(i) * time.Second)})
	}

	if s.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", s.Len())
	}

	list := s.List()
	if list[0].ID != "b" {
		t.Errorf("expected oldest entry to be evicted, first entry is %s", list[0].ID)
	}
}

func TestStore_Clear(t *testing.T) {
	s := NewStore(5)
	s.Add(&Entry{ID: "x"})
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("expected 0 entries after clear, got %d", s.Len())
	}
}

func TestRecorder_SetsTimestamp(t *testing.T) {
	s := NewStore(5)
	r := NewRecorder(s)

	e := &Entry{ID: "ts-test", Method: "/svc/m"}
	before := time.Now()
	r.Record(context.Background(), e)
	after := time.Now()

	list := s.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(list))
	}
	ts := list[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestRecorder_PreservesExistingTimestamp(t *testing.T) {
	s := NewStore(5)
	r := NewRecorder(s)

	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	e := &Entry{ID: "fixed", Method: "/svc/m", Timestamp: fixed}
	r.Record(context.Background(), e)

	if s.List()[0].Timestamp != fixed {
		t.Errorf("expected preserved timestamp %v", fixed)
	}
}
