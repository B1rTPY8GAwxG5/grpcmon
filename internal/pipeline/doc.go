// Package pipeline provides a composable entry-processing pipeline.
//
// Processors are plain functions that accept and return slices of
// capture.Entry, making it straightforward to chain filtering, deduplication,
// truncation, or any other transformation in a defined order.
//
// Example:
//
//	p := pipeline.New(
//		filter.Apply(opts),
//		dedupe.Filter,
//	)
//	result := p.Run(store.List())
package pipeline
