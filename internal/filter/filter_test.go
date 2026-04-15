package filter_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/filter"
)

func makeEntry(method, status string, duration time.Duration) capture.Entry {
	return capture.Entry{
		ID:         "test-id",
		Method:     method,
		StatusCode: status,
		Duration:   duration,
		Timestamp:  time.Now(),
	}
}

func TestMatch_ByMethod(t *testing.T) {
	e := makeEntry("/helloworld.Greeter/SayHello", "OK", 10*time.Millisecond)

	if !filter.Match(e, filter.Criteria{Method: "Greeter"}) {
		t.Error("expected match on partial method name")
	}
	if filter.Match(e, filter.Criteria{Method: "Farewell"}) {
		t.Error("expected no match on unrelated method")
	}
}

func TestMatch_ByStatusCode(t *testing.T) {
	e := makeEntry("/svc/Method", "NOT_FOUND", 5*time.Millisecond)

	if !filter.Match(e, filter.Criteria{StatusCode: "not_found"}) {
		t.Error("expected case-insensitive status match")
	}
	if filter.Match(e, filter.Criteria{StatusCode: "OK"}) {
		t.Error("expected no match on different status")
	}
}

func TestMatch_ByLatency(t *testing.T) {
	e := makeEntry("/svc/Method", "OK", 50*time.Millisecond)

	if !filter.Match(e, filter.Criteria{MinLatency: 10, MaxLatency: 100}) {
		t.Error("expected match within latency range")
	}
	if filter.Match(e, filter.Criteria{MinLatency: 100}) {
		t.Error("expected no match when below min latency")
	}
	if filter.Match(e, filter.Criteria{MaxLatency: 20}) {
		t.Error("expected no match when above max latency")
	}
}

func TestApply_FiltersEntries(t *testing.T) {
	entries := []capture.Entry{
		makeEntry("/svc/Foo", "OK", 10*time.Millisecond),
		makeEntry("/svc/Bar", "NOT_FOUND", 200*time.Millisecond),
		makeEntry("/svc/Foo", "OK", 30*time.Millisecond),
	}

	result := filter.Apply(entries, filter.Criteria{Method: "Foo", MaxLatency: 50})
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestApply_EmptyCriteriaReturnsAll(t *testing.T) {
	entries := []capture.Entry{
		makeEntry("/svc/A", "OK", 1*time.Millisecond),
		makeEntry("/svc/B", "OK", 2*time.Millisecond),
	}

	result := filter.Apply(entries, filter.Criteria{})
	if len(result) != len(entries) {
		t.Fatalf("expected %d results, got %d", len(entries), len(result))
	}
}
