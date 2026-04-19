// Package metric provides per-method rate and error tracking over a sliding window.
package metric

import (
	"sync"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Window holds aggregated counters for a fixed time bucket.
type Window struct {
	Start    time.Time
	Total    int
	Errors   int
	Latency  []float64
}

// Tracker accumulates per-method metrics.
type Tracker struct {
	mu      sync.Mutex
	buckets map[string][]*Window
	size    time.Duration
}

// New returns a Tracker that uses bucketSize as the window granularity.
func New(bucketSize time.Duration) *Tracker {
	if bucketSize <= 0 {
		bucketSize = time.Minute
	}
	return &Tracker{buckets: make(map[string][]*Window), size: bucketSize}
}

// Record adds an entry to the appropriate time bucket.
func (t *Tracker) Record(e capture.Entry) {
	t.mu.Lock()
	defer t.mu.Unlock()

	key := e.Method
	start := e.Timestamp.Truncate(t.size)

	windows := t.buckets[key]
	var w *Window
	if len(windows) > 0 && windows[len(windows)-1].Start.Equal(start) {
		w = windows[len(windows)-1]
	} else {
		w = &Window{Start: start}
		t.buckets[key] = append(windows, w)
	}

	w.Total++
	if e.StatusCode != 0 {
		w.Errors++
	}
	w.Latency = append(w.Latency, float64(e.LatencyMS))
}

// Windows returns a copy of the recorded windows for a method.
func (t *Tracker) Windows(method string) []Window {
	t.mu.Lock()
	defer t.mu.Unlock()

	src := t.buckets[method]
	out := make([]Window, len(src))
	for i, w := range src {
		out[i] = *w
	}
	return out
}

// Methods returns all tracked method names.
func (t *Tracker) Methods() []string {
	t.mu.Lock()
	defer t.mu.Unlock()

	keys := make([]string, 0, len(t.buckets))
	for k := range t.buckets {
		keys = append(keys, k)
	}
	return keys
}
