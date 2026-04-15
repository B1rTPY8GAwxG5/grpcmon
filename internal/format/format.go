// Package format provides utilities for rendering gRPC capture entries
// as human-readable text for CLI output.
package format

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/grpcmon/internal/capture"
)

const (
	// StatusOK is the gRPC status code for a successful call.
	StatusOK = 0
)

// Formatter controls how entries are rendered.
type Formatter struct {
	// Verbose includes request/response bodies when true.
	Verbose bool
	// TimeFormat is the layout string passed to time.Format.
	TimeFormat string
}

// DefaultFormatter returns a Formatter with sensible defaults.
func DefaultFormatter() Formatter {
	return Formatter{
		Verbose:    false,
		TimeFormat: time.RFC3339,
	}
}

// Fprint writes a single entry to w using the formatter's settings.
func (f Formatter) Fprint(w io.Writer, e capture.Entry) error {
	status := statusLabel(e.StatusCode)
	latency := e.Duration.Round(time.Millisecond)

	_, err := fmt.Fprintf(w, "[%s] %s %s (%s) id=%s\n",
		e.Timestamp.Format(f.TimeFormat),
		status,
		e.Method,
		latency,
		e.ID,
	)
	if err != nil {
		return err
	}

	if f.Verbose {
		if len(e.Request) > 0 {
			_, err = fmt.Fprintf(w, "  req: %s\n", strings.TrimSpace(string(e.Request)))
			if err != nil {
				return err
			}
		}
		if len(e.Response) > 0 {
			_, err = fmt.Fprintf(w, "  res: %s\n", strings.TrimSpace(string(e.Response)))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// FprintAll writes every entry in the slice to w.
func (f Formatter) FprintAll(w io.Writer, entries []capture.Entry) error {
	for _, e := range entries {
		if err := f.Fprint(w, e); err != nil {
			return err
		}
	}
	return nil
}

func statusLabel(code int32) string {
	if code == StatusOK {
		return "OK  "
	}
	return fmt.Sprintf("ERR(%d)", code)
}
