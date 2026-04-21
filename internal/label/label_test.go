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

// TestForStatus_NoColour_AllCodes verifies that ForStatus with colour=false
// returns the plain string representation of the status code for a selection
// of common codes.
func TestForStatus_NoColour_AllCodes(t *testing.T) {
	cases := []struct {
		code codes.Code
		want string
	}{
		{codes.OK, "OK"},
		{codes.Canceled, "Canceled"},
		{codes.NotFound, "NotFound"},
		{codes.Internal, "Internal"},
		{codes.Unavailable, "Unavailable"},
	}
	for _, c := range cases {
		got := label.ForStatus(c.code, false)
		if got != c.want {
			t.Errorf("ForStatus(%v, false) = %q, want %q", c.code, got, c.want)
		}
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
