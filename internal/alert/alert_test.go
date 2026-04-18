package alert_test

import (
	"context"
	"testing"
	"time"

	"github.com/grpcmon/internal/alert"
	"github.com/grpcmon/internal/capture"
)

func makeEntries(methods []string, codes []uint32, latencies []time.Duration) []capture.Entry {
	var out []capture.Entry
	for i := range methods {
		out = append(out, capture.Entry{
			Method:     methods[i],
			StatusCode: codes[i],
			Latency:    latencies[i],
		})
	}
	return out
}

func TestEvaluate_NoAlerts_WhenBelowThresholds(t *testing.T) {
	entries := makeEntries(
		[]string{"/svc/Method", "/svc/Method"},
		[]uint32{0, 0},
		[]time.Duration{10 * time.Millisecond, 20 * time.Millisecond},
	)
	e := alert.New([]alert.Rule{{MaxErrorRate: 0.5, MaxLatency: time.Second}})
	alerts := e.Evaluate(context.Background(), entries)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts, got %d", len(alerts))
	}
}

func TestEvaluate_ErrorRateAlert(t *testing.T) {
	entries := makeEntries(
		[]string{"/svc/M", "/svc/M", "/svc/M"},
		[]uint32{2, 2, 0},
		[]time.Duration{1, 1, 1},
	)
	e := alert.New([]alert.Rule{{MaxErrorRate: 0.5}})
	alerts := e.Evaluate(context.Background(), entries)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
}

func TestEvaluate_LatencyAlert(t *testing.T) {
	latencies := make([]time.Duration, 10)
	for i := range latencies {
		latencies[i] = time.Duration(i+1) * 100 * time.Millisecond
	}
	entries := makeEntries(
		make([]string, 10),
		make([]uint32, 10),
		latencies,
	)
	e := alert.New([]alert.Rule{{MaxLatency: 500 * time.Millisecond}})
	alerts := e.Evaluate(context.Background(), entries)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
}

func TestEvaluate_MethodFilter(t *testing.T) {
	entries := makeEntries(
		[]string{"/a/M", "/b/M", "/b/M"},
		[]uint32{2, 0, 0},
		[]time.Duration{1, 1, 1},
	)
	// Rule applies only to /a/M which has 100% errors, but /b/M is fine.
	e := alert.New([]alert.Rule{{Method: "/a/M", MaxErrorRate: 0.5}})
	alerts := e.Evaluate(context.Background(), entries)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert for /a/M, got %d", len(alerts))
	}
}

func TestEvaluate_EmptyEntries_NoAlerts(t *testing.T) {
	e := alert.New([]alert.Rule{{MaxErrorRate: 0.1}})
	alerts := e.Evaluate(context.Background(), nil)
	if len(alerts) != 0 {
		t.Fatalf("expected no alerts for empty entries")
	}
}
