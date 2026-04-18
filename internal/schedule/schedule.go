// Package schedule provides periodic replay scheduling for captured gRPC entries.
package schedule

import (
	"context"
	"time"

	"github.com/grpcmon/internal/capture"
	"github.com/grpcmon/internal/replay"
)

// Job holds configuration for a scheduled replay job.
type Job struct {
	Interval time.Duration
	Options  replay.Options
	OnResult func(replay.Result)
}

// Scheduler runs replay jobs on a fixed interval.
type Scheduler struct {
	store *capture.Store
	job   Job
}

// New creates a new Scheduler.
func New(store *capture.Store, job Job) *Scheduler {
	if job.Interval <= 0 {
		job.Interval = 30 * time.Second
	}
	return &Scheduler{store: store, job: job}
}

// Run starts the scheduler, blocking until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.job.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			s.runOnce(ctx)
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	entries := s.store.List()
	if len(entries) == 0 {
		return
	}
	r, err := replay.New(s.job.Options)
	if err != nil {
		return
	}
	defer r.Close()
	ch := r.ReplayAll(ctx, entries)
	for res := range ch {
		if s.job.OnResult != nil {
			s.job.OnResult(res)
		}
	}
}
