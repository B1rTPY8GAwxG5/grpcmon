// Package watch implements a lightweight polling watcher for the
// grpcmon capture store.
//
// A Watcher periodically checks whether new entries have been added to
// a [capture.Store] and invokes a caller-supplied [Handler] with the
// slice of entries that appeared since the previous poll.
//
// Typical usage:
//
//	w := watch.New(store, 500*time.Millisecond, func(entries []capture.Entry) {
//		for _, e := range entries {
//			fmt.Println("new entry:", e.Method)
//		}
//	})
//	w.Run(ctx)
package watch
