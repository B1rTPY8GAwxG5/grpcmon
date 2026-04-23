// Package audit records replay and export actions performed during a session,
// providing a lightweight operation log for debugging and traceability.
package audit

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Kind classifies the type of audited operation.
type Kind string

const (
	KindReplay Kind = "replay"
	KindExport Kind = "export"
	KindSnapshot Kind = "snapshot"
	KindFilter Kind = "filter"
)

// Event represents a single audited operation.
type Event struct {
	Kind      Kind
	Detail    string
	OccurredAt time.Time
	Err       error
}

// String returns a human-readable representation of the event.
func (e Event) String() string {
	status := "ok"
	if e.Err != nil {
		status = "err: " + e.Err.Error()
	}
	return fmt.Sprintf("%s [%s] %s (%s)",
		e.OccurredAt.Format(time.RFC3339), e.Kind, e.Detail, status)
}

// Log holds an ordered list of audit events.
type Log struct {
	mu     sync.Mutex
	events []Event
	max    int
}

// New creates a Log that retains at most max events (oldest evicted first).
// If max is zero it defaults to 256.
func New(max int) *Log {
	if max <= 0 {
		max = 256
	}
	return &Log{max: max}
}

// Record appends a new event to the log.
func (l *Log) Record(kind Kind, detail string, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.events) >= l.max {
		l.events = l.events[1:]
	}
	l.events = append(l.events, Event{
		Kind:       kind,
		Detail:     detail,
		OccurredAt: time.Now(),
		Err:        err,
	})
}

// List returns a shallow copy of all recorded events.
func (l *Log) List() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Clear removes all events from the log.
func (l *Log) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = l.events[:0]
}

// Fprint writes all events to w, one per line.
func Fprint(w io.Writer, l *Log) {
	for _, e := range l.List() {
		fmt.Fprintln(w, e.String())
	}
}
