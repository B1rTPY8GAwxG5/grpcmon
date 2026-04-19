package label_test

import (
	"strings"
	"testing"

	"google.golang.org/grpc/codes"

	"github.com/example/grpcmon/internal/label"
)

func TestForStatus_NoColour(t *testing.T) {
	got := label.ForStatus(codes.OK, false)
	if got != "OK" {
		t.Fatalf("expected OK, got %q", got)
	}
}

func TestForStatus_ColourOK(t *testing.T) {
	got := label.ForStatus(codes.OK, true)
	if !strings.Contains(got, "OK") {
		t.Fatalf("expected label to contain OK, got %q", got)
	}
	if !strings.Contains(got, "\033[") {
		t.Fatalf("expected ANSI escape in coloured output")
	}
}

func TestForStatus_ColourError(t *testing.T) {
	got := label.ForStatus(codes.Internal, true)
	if !strings.Contains(got, "Internal") {
		t.Fatalf("expected Internal in label, got %q", got)
	}
}

func TestLatencyBand(t *testing.T) {
	cases := []struct {
		ms   float64
		want string
	}{
		{10, "fast"},
		{100, "ok"},
		{500, "slow"},
		{2000, "very slow"},
	}
	for _, c := range cases {
		got := label.LatencyBand(c.ms)
		if got != c.want {
			t.Errorf("LatencyBand(%.0f) = %q, want %q", c.ms, got, c.want)
		}
	}
}

func TestLatencyBandColour_NoColour(t *testing.T) {
	got := label.LatencyBandColour(10, false)
	if got != "fast" {
		t.Fatalf("expected fast, got %q", got)
	}
}

func TestLatencyBandColour_WithColour(t *testing.T) {
	got := label.LatencyBandColour(2000, true)
	if !strings.Contains(got, "very slow") {
		t.Fatalf("expected 'very slow' in output, got %q", got)
	}
	if !strings.Contains(got, "\033[") {
		t.Fatalf("expected ANSI escape in coloured output")
	}
}
