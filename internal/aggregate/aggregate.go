// Package aggregate provides time-windowed aggregation of captured gRPC entries.
package aggregate

import (
	"sync"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Window holds aggregated metrics for a time bucket.
type Window struct {
	Start      time.Time
	End        time.Time
	Count      int
	ErrorCount int
	TotalMS    float64
}

// AvgLatencyMS returns the mean latency across entries in the window.
func (w Window) AvgLatencyMS() float64 {
	if w.Count == 0 {
		return 0
	}
	return w.TotalMS / float64(w.Count)
}

// Aggregator buckets entries into fixed-duration windows.
type Aggregator struct {
	mu       sync.Mutex
	window   time.Duration
	buckets  map[time.Time]*Window
}

// New creates an Aggregator with the given bucket duration.
func New(window time.Duration) *Aggregator {
	if window <= 0 {
		window = time.Minute
	}
	return &Aggregator{window: window, buckets: make(map[time.Time]*Window)}
}

// Add records an entry into the appropriate time bucket.
func (a *Aggregator) Add(e capture.Entry) {
	t := e.Timestamp.Truncate(a.window)
	a.mu.Lock()
	defer a.mu.Unlock()
	w, ok := a.buckets[t]
	if !ok {
		w = &Window{Start: t, End: t.Add(a.window)}
		a.buckets[t] = w
	}
	w.Count++
	if e.StatusCode != 0 {
		w.ErrorCount++
	}
	w.TotalMS += float64(e.Duration.Milliseconds())
}

// Windows returns all collected windows sorted by start time.
func (a *Aggregator) Windows() []Window {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]Window, 0, len(a.buckets))
	for _, w := range a.buckets {
		out = append(out, *w)
	}
	sortWindows(out)
	return out
}

// Reset clears all buckets.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buckets = make(map[time.Time]*Window)
}
