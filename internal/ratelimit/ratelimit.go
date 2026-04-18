// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the speed of gRPC replay operations.
package ratelimit

import (
	"context"
	"time"
)

// Limiter controls the rate at which operations are allowed to proceed.
type Limiter struct {
	ticker *time.Ticker
	tokens chan struct{}
	stop   chan struct{}
}

// New creates a Limiter that allows up to rps operations per second.
// rps must be greater than zero.
func New(rps int) *Limiter {
	if rps <= 0 {
		rps = 1
	}
	interval := time.Second / time.Duration(rps)
	l := &Limiter{
		ticker: time.NewTicker(interval),
		tokens: make(chan struct{}, rps),
		stop:   make(chan struct{}),
	}
	go l.fill()
	return l
}

func (l *Limiter) fill() {
	for {
		select {
		case <-l.ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
			}
		case <-l.stop:
			return
		}
	}
}

// Wait blocks until a token is available or the context is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases resources held by the Limiter.
func (l *Limiter) Stop() {
	l.ticker.Stop()
	close(l.stop)
}
