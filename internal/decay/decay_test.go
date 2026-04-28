package decay_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/decay"
)

func makeEntry(method string, age time.Duration) capture.Entry {
	return capture.Entry{
		Method:    method,
		Timestamp: time.Now().Add(-age),
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := decay.DefaultOptions()
	if opts.HalfLife != 5*time.Minute {
		t.Fatalf("expected 5m half-life, got %v", opts.HalfLife)
	}
}

func TestScore_RecentEntryIsNearOne(t *testing.T) {
	s := decay.New(decay.DefaultOptions())
	e := makeEntry("/svc/Method", 0)
	sc := s.Score(e)
	if sc < 0.99 {
		t.Fatalf("expected score near 1 for brand-new entry, got %f", sc)
	}
}

func TestScore_ZeroTimestamp_ReturnsZero(t *testing.T) {
	s := decay.New(decay.DefaultOptions())
	sc := s.Score(capture.Entry{})
	if sc != 0 {
		t.Fatalf("expected 0 for zero timestamp, got %f", sc)
	}
}

func TestScore_HalfAfterHalfLife(t *testing.T) {
	opts := decay.Options{HalfLife: time.Second}
	s := decay.New(opts)
	e := makeEntry("/svc/Old", time.Second)
	sc := s.Score(e)
	// allow ±5% tolerance for timing
	if sc < 0.45 || sc > 0.55 {
		t.Fatalf("expected ~0.5 after one half-life, got %f", sc)
	}
}

func TestApply_DropsEntriesBelowThreshold(t *testing.T) {
	opts := decay.Options{HalfLife: time.Millisecond}
	s := decay.New(opts)

	entries := []capture.Entry{
		makeEntry("/svc/New", 0),
		makeEntry("/svc/Old", 10*time.Second),
	}

	out := s.Apply(entries, 0.5)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry above threshold, got %d", len(out))
	}
	if out[0].Method != "/svc/New" {
		t.Fatalf("expected /svc/New, got %s", out[0].Method)
	}
}

func TestApply_SortsByScoreDescending(t *testing.T) {
	s := decay.New(decay.DefaultOptions())

	entries := []capture.Entry{
		makeEntry("/svc/Old", 4*time.Minute),
		makeEntry("/svc/New", 0),
		makeEntry("/svc/Mid", 2*time.Minute),
	}

	out := s.Apply(entries, 0)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0].Method != "/svc/New" {
		t.Fatalf("expected newest first, got %s", out[0].Method)
	}
	if out[2].Method != "/svc/Old" {
		t.Fatalf("expected oldest last, got %s", out[2].Method)
	}
}

func TestNew_ZeroHalfLifeDefaultsToFiveMinutes(t *testing.T) {
	s := decay.New(decay.Options{HalfLife: 0})
	e := makeEntry("/svc/Method", 0)
	if s.Score(e) < 0.99 {
		t.Fatal("expected valid score with defaulted half-life")
	}
}
