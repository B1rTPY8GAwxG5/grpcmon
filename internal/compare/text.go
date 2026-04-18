package compare

import (
	"fmt"
	"io"
	"strings"
)

// WriteTextReport writes a detailed text diff report to w.
func WriteTextReport(w io.Writer, r Report) {
	total := r.MatchCount + r.MismatchCount
	fmt.Fprintf(w, "=== Compare Report ===\n")
	fmt.Fprintf(w, "Total: %d | Match: %d | Mismatch: %d\n", total, r.MatchCount, r.MismatchCount)
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, res := range r.Results {
		if res.Match {
			fmt.Fprintf(w, "  [MATCH]    %s\n", res.Method)
		} else {
			fmt.Fprintf(w, "  [MISMATCH] %s\n", res.Method)
			for _, d := range res.Diffs {
				fmt.Fprintf(w, "    - %s\n", d)
			}
		}
	}
	fmt.Fprintln(w, strings.Repeat("=", 40))
}
