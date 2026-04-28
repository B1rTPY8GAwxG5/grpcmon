package prestige_test

import (
	"testing"
	"time"

	"github.com/example/grpcmon/internal/capture"
	"github.com/example/grpcmon/internal/prestige"
)

func makeEntry(method string, statusCode int, latencyMS int64, age time.Duration) capture.Entry {
	return capture.Entry{
		ID:         "id-" + method,
		Method:     method,
		StatusCode: statusCode,
		LatencyMS:  latencyMS,
		Timestamp:  time.Now().Add(-age),
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := prestige.DefaultOptions()
	if opts.RecencyWeight != 0.3 {
		t.Errorf("RecencyWeight = %v, want 0.3", opts.RecencyWeight)
	}
	if opts.LatencyWeight != 0.4 {
		t.Errorf("LatencyWeight = %v, want 0.4", opts.LatencyWeight)
	}
	if opts.ErrorWeight != 0.3 {
		t.Errorf("ErrorWeight = %v, want 0.3", opts.ErrorWeight)
	}
	if opts.LatencyThresholdMS != 500 {
		t.Errorf("LatencyThresholdMS = %v, want 500", opts.LatencyThresholdMS)
	}
}

func TestRank_EmptyEntries(t *testing.T) {
	result := prestige.Rank(nil, prestige.DefaultOptions())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestRank_HighLatencyScoresHigher(t *testing.T) {
	low := makeEntry("/svc/Low", 0, 50, time.Second)
	high := makeEntry("/svc/High", 0, 1000, time.Second)

	scores := prestige.Rank([]capture.Entry{low, high}, prestige.DefaultOptions())
	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}
	if scores[0].Entry.Method != "/svc/High" {
		t.Errorf("expected High to rank first, got %s", scores[0].Entry.Method)
	}
}

func TestRank_ErrorScoresHigher(t *testing.T) {
	ok := makeEntry("/svc/OK", 0, 10, time.Second)
	err := makeEntry("/svc/Err", 2, 10, time.Second)

	scores := prestige.Rank([]capture.Entry{ok, err}, prestige.DefaultOptions())
	if scores[0].Entry.Method != "/svc/Err" {
		t.Errorf("expected Err to rank first, got %s", scores[0].Entry.Method)
	}
}

func TestRank_RecentScoresHigher(t *testing.T) {
	old := makeEntry("/svc/Old", 0, 10, 10*time.Minute)
	recent := makeEntry("/svc/Recent", 0, 10, time.Millisecond)

	scores := prestige.Rank([]capture.Entry{old, recent}, prestige.DefaultOptions())
	if scores[0].Entry.Method != "/svc/Recent" {
		t.Errorf("expected Recent to rank first, got %s", scores[0].Entry.Method)
	}
}

func TestTop_LimitsResults(t *testing.T) {
	entries := []capture.Entry{
		makeEntry("/a", 0, 100, time.Second),
		makeEntry("/b", 0, 200, time.Second),
		makeEntry("/c", 0, 300, time.Second),
	}
	top := prestige.Top(entries, 2, prestige.DefaultOptions())
	if len(top) != 2 {
		t.Errorf("expected 2, got %d", len(top))
	}
}

func TestTop_NGreaterThanEntries(t *testing.T) {
	entries := []capture.Entry{
		makeEntry("/a", 0, 100, time.Second),
	}
	top := prestige.Top(entries, 10, prestige.DefaultOptions())
	if len(top) != 1 {
		t.Errorf("expected 1, got %d", len(top))
	}
}
