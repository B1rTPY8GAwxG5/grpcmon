package format

import (
	"fmt"
	"io"
	"time"

	"github.com/grpcmon/internal/capture"
)

const (
	colWidth = 28
	header   = "%-6s %-28s %-8s %-10s %s\n"
	row      = "%-6s %-28s %-8s %-10s %s\n"
)

// WriteTable renders entries as a fixed-width table to w.
func WriteTable(w io.Writer, entries []capture.Entry) error {
	if _, err := fmt.Fprintf(w, header,
		"STATUS", "METHOD", "LATENCY", "TIMESTAMP", "ID",
	); err != nil {
		return err
	}

	sep := fmt.Sprintf("%s\n", repeatDash(80))
	if _, err := fmt.Fprint(w, sep); err != nil {
		return err
	}

	for _, e := range entries {
		method := truncate(e.Method, colWidth)
		latency := e.Duration.Round(time.Millisecond).String()
		ts := e.Timestamp.Format("15:04:05")
		status := statusLabel(e.StatusCode)

		if _, err := fmt.Fprintf(w, row, status, method, latency, ts, e.ID); err != nil {
			return err
		}
	}
	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func repeatDash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '-'
	}
	return string(b)
}
