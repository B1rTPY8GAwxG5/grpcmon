package metric

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestSummarise_Totals(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.Record(makeEntry("/svc/A", 0, 10, now))
	tr.Record(makeEntry("/svc/A", 2, 20, now))

	sums := tr.Summarise()
	if len(sums) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(sums))
	}
	s := sums[0]
	if s.Total != 2 {
		t.Errorf("total: want 2, got %d", s.Total)
	}
	if s.Errors != 1 {
		t.Errorf("errors: want 1, got %d", s.Errors)
	}
	if s.ErrorRate != 0.5 {
		t.Errorf("error rate: want 0.5, got %f", s.ErrorRate)
	}
	if s.AvgLatency != 15.0 {
		t.Errorf("avg latency: want 15, got %f", s.AvgLatency)
	}
}

func TestSummarise_Sorted(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.Record(makeEntry("/z", 0, 1, now))
	tr.Record(makeEntry("/a", 0, 1, now))

	sums := tr.Summarise()
	if sums[0].Method != "/a" {
		t.Errorf("expected /a first, got %s", sums[0].Method)
	}
}

func TestFprint_ContainsMethod(t *testing.T) {
	tr := New(time.Minute)
	tr.Record(makeEntry("/svc/Hello", 0, 5, time.Now()))

	var buf bytes.Buffer
	Fprint(&buf, tr.Summarise())

	if !strings.Contains(buf.String(), "/svc/Hello") {
		t.Errorf("output missing method name")
	}
}

func TestSummarise_Empty(t *testing.T) {
	tr := New(time.Minute)
	if sums := tr.Summarise(); len(sums) != 0 {
		t.Errorf("expected empty, got %d", len(sums))
	}
}
