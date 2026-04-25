package rollup_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/rollup"
)

func entry(method string, latency int64, ts time.Time) capture.Entry {
	return capture.Entry{
		Method:    method,
		LatencyMS: latency,
		Timestamp: ts,
		Response:  method + "-resp",
	}
}

var (
	t1 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 = t1.Add(5 * time.Second)
	t3 = t1.Add(10 * time.Second)
)

func TestMerge_EmptyEntries_ReturnsZero(t *testing.T) {
	result := rollup.Merge(nil, rollup.DefaultOptions())
	if result.Method != "" {
		t.Fatalf("expected zero entry, got %+v", result)
	}
}

func TestMerge_AveragesLatency(t *testing.T) {
	entries := []capture.Entry{
		entry("/svc/A", 100, t1),
		entry("/svc/A", 200, t2),
		entry("/svc/A", 300, t3),
	}

	result := rollup.Merge(entries, rollup.DefaultOptions())
	if result.LatencyMS != 200 {
		t.Fatalf("expected avg latency 200, got %d", result.LatencyMS)
	}
}

func TestMerge_KeepsNewestTimestamp(t *testing.T) {
	entries := []capture.Entry{
		entry("/svc/A", 10, t1),
		entry("/svc/A", 20, t3),
		entry("/svc/A", 30, t2),
	}

	result := rollup.Merge(entries, rollup.DefaultOptions())
	if !result.Timestamp.Equal(t3) {
		t.Fatalf("expected newest timestamp %v, got %v", t3, result.Timestamp)
	}
}

func TestMerge_KeepFirstTimestamp(t *testing.T) {
	opts := rollup.Options{KeepFirstTimestamp: true}
	entries := []capture.Entry{
		entry("/svc/A", 10, t2),
		entry("/svc/A", 20, t1),
		entry("/svc/A", 30, t3),
	}

	result := rollup.Merge(entries, opts)
	if !result.Timestamp.Equal(t1) {
		t.Fatalf("expected oldest timestamp %v, got %v", t1, result.Timestamp)
	}
}

func TestMergeAll_GroupsByMethod(t *testing.T) {
	entries := []capture.Entry{
		entry("/svc/A", 100, t1),
		entry("/svc/B", 50, t1),
		entry("/svc/A", 200, t2),
	}

	result := rollup.MergeAll(entries, rollup.DefaultOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result))
	}

	byMethod := make(map[string]capture.Entry)
	for _, r := range result {
		byMethod[r.Method] = r
	}

	if byMethod["/svc/A"].LatencyMS != 150 {
		t.Fatalf("expected /svc/A avg latency 150, got %d", byMethod["/svc/A"].LatencyMS)
	}
	if byMethod["/svc/B"].LatencyMS != 50 {
		t.Fatalf("expected /svc/B latency 50, got %d", byMethod["/svc/B"].LatencyMS)
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := rollup.DefaultOptions()
	if opts.KeepFirstTimestamp {
		t.Fatal("expected KeepFirstTimestamp to be false by default")
	}
}
