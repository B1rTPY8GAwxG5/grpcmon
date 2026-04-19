// Package aggregate groups captured gRPC entries into fixed-duration time
// windows and exposes per-window metrics such as request count, error count,
// and average latency.
//
// Usage:
//
//	a := aggregate.New(time.Minute)
//	for _, e := range store.List() {
//		a.Add(e)
//	}
//	for _, w := range a.Windows() {
//		fmt.Printf("%s count=%d errors=%d avgMs=%.1f\n",
//			w.Start.Format(time.RFC3339), w.Count, w.ErrorCount, w.AvgLatencyMS())
//	}
package aggregate
