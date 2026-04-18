package summary_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"grpcmon/internal/stats"
	"grpcmon/internal/summary"
)

func sampleStats() stats.Stats {
	return stats.Stats{
		Total:      10,
		Successful: 8,
		Failed:     2,
		ErrorRate:  0.2,
		P50:        50 * time.Millisecond,
		P99:        200 * time.Millisecond,
		TopMethods: []stats.MethodCount{{Method: "/svc/Hello", Count: 5}},
	}
}

func TestFromStats_Fields(t *testing.T) {
	r := summary.FromStats(sampleStats())
	if r.Total != 10 {
		t.Errorf("Total: got %d want 10", r.Total)
	}
	if r.Failed != 2 {
		t.Errorf("Failed: got %d want 2", r.Failed)
	}
	if r.TopMethod != "/svc/Hello" {
		t.Errorf("TopMethod: got %q", r.TopMethod)
	}
}

func TestFromStats_NoTopMethod(t *testing.T) {
	s := sampleStats()
	s.TopMethods = nil
	r := summary.FromStats(s)
	if r.TopMethod != "" {
		t.Errorf("expected empty TopMethod, got %q", r.TopMethod)
	}
}

func TestFprint_ContainsKeyFields(t *testing.T) {
	var buf bytes.Buffer
	summary.Fprint(&buf, summary.FromStats(sampleStats()))
	out := buf.String()
	for _, want := range []string{"Total", "Error rate", "p50", "p99", "/svc/Hello"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestFprint_NoTopMethodLine(t *testing.T) {
	s := sampleStats()
	s.TopMethods = nil
	var buf bytes.Buffer
	summary.Fprint(&buf, summary.FromStats(s))
	if strings.Contains(buf.String(), "Top method") {
		t.Error("should not print Top method line when empty")
	}
}
