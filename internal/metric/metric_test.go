package metric

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
)

func makeEntry(method string, code int, latency int64, ts time.Time) capture.Entry {
	return capture.Entry{
		Method:     method,
		StatusCode: code,
		LatencyMS:  latency,
		Timestamp:  ts,
	}
}

func TestRecord_BucketsEntry(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now().Truncate(time.Minute)
	tr.Record(makeEntry("/svc/Method", 0, 10, now))

	win := tr.Windows("/svc/Method")
	if len(win) != 1 {
		t.Fatalf("expected 1 window, got %d", len(win))
	}
	if win[0].Total != 1 {
		t.Errorf("expected total 1, got %d", win[0].Total)
	}
}

func TestRecord_ErrorCount(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now().Truncate(time.Minute)
	tr.Record(makeEntry("/svc/M", 0, 5, now))
	tr.Record(makeEntry("/svc/M", 2, 5, now))

	win := tr.Windows("/svc/M")
	if win[0].Errors != 1 {
		t.Errorf("expected 1 error, got %d", win[0].Errors)
	}
}

func TestRecord_SeparateBuckets(t *testing.T) {
	tr := New(time.Minute)
	t1 := time.Now().Truncate(time.Minute)
	t2 := t1.Add(time.Minute)
	tr.Record(makeEntry("/svc/M", 0, 5, t1))
	tr.Record(makeEntry("/svc/M", 0, 5, t2))

	win := tr.Windows("/svc/M")
	if len(win) != 2 {
		t.Errorf("expected 2 windows, got %d", len(win))
	}
}

func TestMethods_ReturnsKeys(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.Record(makeEntry("/a", 0, 1, now))
	tr.Record(makeEntry("/b", 0, 1, now))

	ms := tr.Methods()
	if len(ms) != 2 {
		t.Errorf("expected 2 methods, got %d", len(ms))
	}
}

func TestNew_ZeroBucketDefaultsToMinute(t *testing.T) {
	tr := New(0)
	if tr.size != time.Minute {
		t.Errorf("expected minute, got %v", tr.size)
	}
}
