// Package alert provides threshold-based alerting for captured gRPC traffic.
package alert

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/grpcmon/internal/capture"
)

// Rule defines conditions that trigger an alert.
type Rule struct {
	MaxErrorRate float64
	MaxLatency   time.Duration
	Method       string
}

// Alert describes a triggered rule.
type Alert struct {
	Rule    Rule
	Message string
	At      time.Time
}

// Evaluator checks entries against a set of rules.
type Evaluator struct {
	rules []Rule
}

// New returns an Evaluator configured with the given rules.
func New(rules []Rule) *Evaluator {
	return &Evaluator{rules: rules}
}

// Evaluate inspects entries and returns any triggered alerts.
func (e *Evaluator) Evaluate(_ context.Context, entries []capture.Entry) []Alert {
	var alerts []Alert
	for _, rule := range e.rules {
		matched := filterByMethod(entries, rule.Method)
		if len(matched) == 0 {
			continue
		}
		if rule.MaxErrorRate > 0 {
			if rate := errorRate(matched); rate > rule.MaxErrorRate {
				alerts = append(alerts, Alert{
					Rule:    rule,
					Message: fmt.Sprintf("error rate %.2f exceeds threshold %.2f", rate, rule.MaxErrorRate),
					At:      time.Now(),
				})
			}
		}
		if rule.MaxLatency > 0 {
			if p99 := p99Latency(matched); p99 > rule.MaxLatency {
				alerts = append(alerts, Alert{
					Rule:    rule,
					Message: fmt.Sprintf("p99 latency %s exceeds threshold %s", p99, rule.MaxLatency),
					At:      time.Now(),
				})
			}
		}
	}
	return alerts
}

func filterByMethod(entries []capture.Entry, method string) []capture.Entry {
	if method == "" {
		return entries
	}
	var out []capture.Entry
	for _, e := range entries {
		if e.Method == method {
			out = append(out, e)
		}
	}
	return out
}

func errorRate(entries []capture.Entry) float64 {
	var errs int
	for _, e := range entries {
		if e.StatusCode != 0 {
			errs++
		}
	}
	return float64(errs) / float64(len(entries))
}

func p99Latency(entries []capture.Entry) time.Duration {
	latencies := make([]time.Duration, len(entries))
	for i, e := range entries {
		latencies[i] = e.Latency
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	idx := int(float64(len(latencies))*0.99)
	if idx >= len(latencies) {
		idx = len(latencies) - 1
	}
	return latencies[idx]
}
