package sampler_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/sampler"
)

func makeEntries(n int) []capture.Entry {
	entries := make([]capture.Entry, n)
	for i := range entries {
		entries[i] = capture.Entry{Method: "/svc/Method"}
	}
	return entries
}

func TestNew_ClampsRateBelow(t *testing.T) {
	s := sampler.New(-0.5, nil)
	if s.Rate() != 0 {
		t.Fatalf("expected rate 0, got %f", s.Rate())
	}
}

func TestNew_ClampsRateAbove(t *testing.T) {
	s := sampler.New(1.5, nil)
	if s.Rate() != 1 {
		t.Fatalf("expected rate 1, got %f", s.Rate())
	}
}

func TestKeep_RateZero_RejectsAll(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	s := sampler.New(0, src)
	for i := 0; i < 100; i++ {
		if s.Keep(capture.Entry{}) {
			t.Fatal("expected Keep to return false with rate 0")
		}
	}
}

func TestKeep_RateOne_AcceptsAll(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	s := sampler.New(1.0, src)
	for i := 0; i < 100; i++ {
		if !s.Keep(capture.Entry{}) {
			t.Fatal("expected Keep to return true with rate 1")
		}
	}
}

func TestFilter_ReducesEntries(t *testing.T) {
	src := rand.NewSource(99)
	s := sampler.New(0.5, src)
	entries := makeEntries(1000)
	result := s.Filter(entries)
	if len(result) == 0 || len(result) == 1000 {
		t.Fatalf("expected partial retention, got %d", len(result))
	}
}

func TestFilter_RateZero_ReturnsEmpty(t *testing.T) {
	s := sampler.New(0, nil)
	result := s.Filter(makeEntries(50))
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestFilter_RateOne_ReturnsAll(t *testing.T) {
	s := sampler.New(1.0, nil)
	entries := makeEntries(50)
	result := s.Filter(entries)
	if len(result) != 50 {
		t.Fatalf("expected 50 entries, got %d", len(result))
	}
}

func TestSetRate_UpdatesRate(t *testing.T) {
	s := sampler.New(0.5, nil)
	s.SetRate(0.9)
	if s.Rate() != 0.9 {
		t.Fatalf("expected rate 0.9, got %f", s.Rate())
	}
}

func TestSetRate_ClampsValue(t *testing.T) {
	s := sampler.New(0.5, nil)
	s.SetRate(2.0)
	if s.Rate() != 1.0 {
		t.Fatalf("expected clamped rate 1.0, got %f", s.Rate())
	}
}
