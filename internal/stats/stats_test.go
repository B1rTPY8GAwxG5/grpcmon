package stats_test

import (
	"testing"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/stats"
	"google.golang.org/grpc/codes"
)

func makeEntries() []capture.Entry {
	return []capture.Entry{
		{Method: "/svc/MethodA", StatusCode: codes.OK, Latency: 10 * time.Millisecond},
		{Method: "/svc/MethodA", StatusCode: codes.OK, Latency: 20 * time.Millisecond},
		{Method: "/svc/MethodB", StatusCode: codes.NotFound, Latency: 5 * time.Millisecond},
		{Method: "/svc/MethodC", StatusCode: codes.Internal, Latency: 50 * time.Millisecond},
	}
}

func TestCompute_EmptyEntries(t *testing.T) {
	s := stats.Compute(nil)
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
	if s.StatusCodes == nil {
		t.Error("expected non-nil StatusCodes map")
	}
}

func TestCompute_Totals(t *testing.T) {
	s := stats.Compute(makeEntries())
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.SuccessCount != 2 {
		t.Errorf("expected SuccessCount=2, got %d", s.SuccessCount)
	}
	if s.ErrorCount != 2 {
		t.Errorf("expected ErrorCount=2, got %d", s.ErrorCount)
	}
}

func TestCompute_Latency(t *testing.T) {
	s := stats.Compute(makeEntries())
	if s.MinLatency != 5*time.Millisecond {
		t.Errorf("expected MinLatency=5ms, got %v", s.MinLatency)
	}
	if s.MaxLatency != 50*time.Millisecond {
		t.Errorf("expected MaxLatency=50ms, got %v", s.MaxLatency)
	}
	want := (10 + 20 + 5 + 50) * time.Millisecond / 4
	if s.AvgLatency != want {
		t.Errorf("expected AvgLatency=%v, got %v", want, s.AvgLatency)
	}
}

func TestCompute_StatusCodes(t *testing.T) {
	s := stats.Compute(makeEntries())
	if s.StatusCodes[codes.OK] != 2 {
		t.Errorf("expected 2 OK, got %d", s.StatusCodes[codes.OK])
	}
	if s.StatusCodes[codes.NotFound] != 1 {
		t.Errorf("expected 1 NotFound, got %d", s.StatusCodes[codes.NotFound])
	}
}

func TestCompute_TopMethods(t *testing.T) {
	s := stats.Compute(makeEntries())
	if len(s.TopMethods) == 0 {
		t.Fatal("expected non-empty TopMethods")
	}
	if s.TopMethods[0].Method != "/svc/MethodA" {
		t.Errorf("expected top method to be /svc/MethodA, got %s", s.TopMethods[0].Method)
	}
	if s.TopMethods[0].Count != 2 {
		t.Errorf("expected count 2, got %d", s.TopMethods[0].Count)
	}
}
