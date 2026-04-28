// Package prestige scores captured gRPC entries by how "interesting" they are,
// combining recency, latency, and error status into a single numeric value.
//
// Higher-scored entries are surfaced first, making it easy to focus attention
// on slow or failing calls during development.
//
// # Usage
//
//	opts := prestige.DefaultOptions()
//	scores := prestige.Rank(store.List(), opts)
//	for _, s := range scores {
//		fmt.Printf("%.3f  %s\n", s.Value, s.Entry.Method)
//	}
//
// # Weights
//
// The three components (recency, latency, error) each have a configurable
// weight. Adjust Options to emphasise whichever dimension matters most for
// your workflow.
package prestige
