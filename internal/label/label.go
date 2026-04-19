// Package label provides colour-coded string labels for gRPC status codes
// and latency bands used across formatting and TUI components.
package label

import "google.golang.org/grpc/codes"

// Colour escape codes.
const (
	reset  = "\033[0m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	cyan   = "\033[36m"
)

// ForStatus returns a short, optionally coloured label for a gRPC status code.
func ForStatus(code codes.Code, colour bool) string {
	text := code.String()
	if !colour {
		return text
	}
	switch code {
	case codes.OK:
		return green + text + reset
	case codes.DeadlineExceeded, codes.Unavailable, codes.Internal:
		return red + text + reset
	default:
		return yellow + text + reset
	}
}

// LatencyBand returns a human-readable band label for a latency in milliseconds.
func LatencyBand(ms float64) string {
	switch {
	case ms < 50:
		return "fast"
	case ms < 200:
		return "ok"
	case ms < 1000:
		return "slow"
	default:
		return "very slow"
	}
}

// LatencyBandColour returns a coloured latency band label.
func LatencyBandColour(ms float64, colour bool) string {
	band := LatencyBand(ms)
	if !colour {
		return band
	}
	switch band {
	case "fast":
		return cyan + band + reset
	case "ok":
		return green + band + reset
	case "slow":
		return yellow + band + reset
	default:
		return red + band + reset
	}
}
