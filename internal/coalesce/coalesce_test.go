package coalesce

import (
	"testing"
	"time"

	"github.com/example/grpcmon/internal/capture"
)

func entry(method string, status uint32, latency time.Duration, ts time.Time) capture.Entry {
	return capture.Entry{
		ID:         "id-" + method,
		Method:     method,
		StatusCode: status,
		Latency:    latency,
		Timestamp:  ts,
	}
}

func TestFlush_EmptyCoalescer_ReturnsEmpty(t *testing.T) {
	c := New()
	got := c.Flush()
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestFlush_SingleEntry_ReturnedUnchanged(t *testing.T) {
	c := New()
	now := time.Now()
	e := entry("/svc/Method", 0, 10*time.Millisecond, now)
	c.Add(e)
	got := c.Flush()
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].Latency != 10*time.Millisecond {
		t.Errorf("latency: want 10ms, got %v", got[0].Latency)
	}
}

func TestFlush_AveragesLatency(t *testing.T) {
	c := New()
	now := time.Now()
	c.Add(entry("/svc/A", 0, 20*time.Millisecond, now))
	c.Add(entry("/svc/A", 0, 40*time.Millisecond, now))
	got := c.Flush()
	if len(got) != 1 {
		t.Fatalf("expected 1 merged entry, got %d", len(got))
	}
	if got[0].Latency != 30*time.Millisecond {
		t.Errorf("latency: want 30ms, got %v", got[0].Latency)
	}
}

func TestFlush_KeepsLatestTimestamp(t *testing.T) {
	c := New()
	earlier := time.Now().Add(-5 * time.Second)
	later := time.Now()
	c.Add(entry("/svc/B", 0, 10*time.Millisecond, earlier))
	c.Add(entry("/svc/B", 0, 10*time.Millisecond, later))
	got := c.Flush()
	if !got[0].Timestamp.Equal(later) {
		t.Errorf("timestamp: want %v, got %v", later, got[0].Timestamp)
	}
}

func TestFlush_SeparatesDistinctMethodStatusPairs(t *testing.T) {
	c := New()
	now := time.Now()
	c.Add(entry("/svc/C", 0, 10*time.Millisecond, now))
	c.Add(entry("/svc/C", 2, 10*time.Millisecond, now))
	c.Add(entry("/svc/D", 0, 10*time.Millisecond, now))
	got := c.Flush()
	if len(got) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(got))
	}
}

func TestFlush_ResetsStateAfterFlush(t *testing.T) {
	c := New()
	now := time.Now()
	c.Add(entry("/svc/E", 0, 5*time.Millisecond, now))
	c.Flush()
	second := c.Flush()
	if len(second) != 0 {
		t.Errorf("expected empty flush after reset, got %d entries", len(second))
	}
}
