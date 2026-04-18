// Package schedule provides periodic replay scheduling.
//
// A Scheduler replays captured gRPC entries against a target at a
// configurable interval, making it easy to run regression checks
// during development without manual intervention.
//
// Example:
//
//	store, _ := capture.NewStore(100)
//	job := schedule.Job{
//		Interval: 30 * time.Second,
//		Options:  replay.DefaultOptions(),
//		OnResult: func(r replay.Result) { fmt.Println(r) },
//	}
//	sched := schedule.New(store, job)
//	sched.Run(ctx)
package schedule
