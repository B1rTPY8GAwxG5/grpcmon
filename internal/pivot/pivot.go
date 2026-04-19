// Package pivot provides utilities for pivoting captured entries
// into tabular summaries keyed by a chosen dimension.
package pivot

import (
	"sort"

	"grpcmon/internal/capture"
)

// Dimension selects the grouping key for a pivot table.
type Dimension int

const (
	ByMethod Dimension = iota
	ByStatus
)

// Row holds aggregated metrics for a single dimension value.
type Row struct {
	Key        string
	Count      int
	ErrorCount int
	TotalMS    float64
	AvgMS      float64
}

// Table is an ordered slice of pivot rows.
type Table []Row

// Build constructs a pivot Table from entries grouped by dim.
func Build(entries []capture.Entry, dim Dimension) Table {
	index := map[string]*Row{}

	for _, e := range entries {
		var key string
		switch dim {
		case ByStatus:
			key = e.Status.String()
		default:
			key = e.Method
		}

		r, ok := index[key]
		if !ok {
			r = &Row{Key: key}
			index[key] = r
		}
		r.Count++
		ms := float64(e.Duration.Milliseconds())
		r.TotalMS += ms
		if e.Status != 0 {
			r.ErrorCount++
		}
	}

	table := make(Table, 0, len(index))
	for _, r := range index {
		if r.Count > 0 {
			r.AvgMS = r.TotalMS / float64(r.Count)
		}
		table = append(table, *r)
	}
	sort.Slice(table, func(i, j int) bool {
		return table[i].Key < table[j].Key
	})
	return table
}
