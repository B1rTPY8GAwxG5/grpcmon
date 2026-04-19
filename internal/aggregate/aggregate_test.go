package aggregate_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/aggregate"
	"github.com/grpcmon/internal/capture"
)

func makeEntry(ts time.Time, statusCode int, dur time.Duration) capture.Entry {
	return capture.Entry{
		Timestamp:  ts,
		StatusCode: statusCode,
		Duration:   dur,
	}
}

func TestAdd_BucketsEntry(t *testing.T) {
	a := aggregate.New(time.Minute)
	now := time.Now().Truncate(time.Minute)
	a.Add(makeEntry(now, 0, 10*time.Millisecond))
	ws := a.Windows()
	if len(ws) != 1 {
		t.Fatalf("expected 1 window, got %d", len(ws))
	}
	if ws[0].Count != 1 {
		t.Errorf("expected count 1, got %d", ws[0].Count)
	}
}

func TestAdd_ErrorCount(t *testing.T) {
	a := aggregate.New(time.Minute)
	now := time.Now().Truncate(time.Minute)
	a.Add(makeEntry(now, 0, 5*time.Millisecond))
	a.Add(makeEntry(now, 2, 5*time.Millisecond))
	ws := a.Windows()
	if ws[0].ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", ws[0].ErrorCount)
	}
}

func TestAvgLatencyMS(t *testing.T) {
	a := aggregate.New(time.Minute)
	now := time.Now().Truncate(time.Minute)
	a.Add(makeEntry(now, 0, 20*time.Millisecond))
	a.Add(makeEntry(now, 0, 40*time.Millisecond))
	ws := a.Windows()
	if ws[0].AvgLatencyMS() != 30 {
		t.Errorf("expected avg 30ms, got %.2f", ws[0].AvgLatencyMS())
	}
}

func TestWindows_SortedByStart(t *testing.T) {
	a := aggregate.New(time.Minute)
	base := time.Now().Truncate(time.Minute)
	a.Add(makeEntry(base.Add(2*time.Minute), 0, 1*time.Millisecond))
	a.Add(makeEntry(base, 0, 1*time.Millisecond))
	a.Add(makeEntry(base.Add(time.Minute), 0, 1*time.Millisecond))
	ws := a.Windows()
	for i := 1; i < len(ws); i++ {
		if ws[i].Start.Before(ws[i-1].Start) {
			t.Errorf("windows not sorted at index %d", i)
		}
	}
}

func TestReset_ClearsBuckets(t *testing.T) {
	a := aggregate.New(time.Minute)
	a.Add(makeEntry(time.Now(), 0, 5*time.Millisecond))
	a.Reset()
	if len(a.Windows()) != 0 {
		t.Error("expected empty windows after reset")
	}
}

func TestNew_ZeroDurationDefaultsToMinute(t *testing.T) {
	a := aggregate.New(0)
	now := time.Now()
	a.Add(makeEntry(now, 0, 1*time.Millisecond))
	ws := a.Windows()
	dur := ws[0].End.Sub(ws[0].Start)
	if dur != time.Minute {
		t.Errorf("expected 1m window, got %s", dur)
	}
}
