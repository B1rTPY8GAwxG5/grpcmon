package tui

import (
	"fmt"

	"github.com/user/grpcmon/internal/capture"
)

// entryItem wraps a capture.Entry so it satisfies the list.Item interface.
type entryItem struct {
	entry capture.Entry
}

// Title returns the list item title shown in the TUI.
func (i entryItem) Title() string {
	status := i.entry.Status
	if status == "" {
		status = "UNKNOWN"
	}
	return fmt.Sprintf("%-40s  %s", i.entry.Method, status)
}

// Description returns secondary text shown beneath the title.
func (i entryItem) Description() string {
	ts := ""
	if !i.entry.Timestamp.IsZero() {
		ts = i.entry.Timestamp.Format("15:04:05.000")
	}
	dur := ""
	if i.entry.Duration > 0 {
		dur = i.entry.Duration.String()
	}
	return fmt.Sprintf("%s  latency: %s", ts, dur)
}

// FilterValue satisfies list.Item; enables future search support.
func (i entryItem) FilterValue() string {
	return i.entry.Method
}
